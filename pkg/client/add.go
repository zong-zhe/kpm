package client

import (
	"fmt"
	"path/filepath"

	"kcl-lang.io/kpm/pkg/constants"
	"kcl-lang.io/kpm/pkg/downloader"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/utils"
	"kcl-lang.io/kpm/pkg/visitor"
)

// AddOptions contains the options for adding a dependency.
type AddOptions struct {
	// Sources is the sources of the package.
	// It can be a local *.k path, a local *.tar/*.tgz path, a local directory, a remote git/oci path,.
	Sources []*downloader.Source
	// The kcl package to add the dependencies.
	KclPkg pkg.KclPkg
}

type AddOption func(*AddOptions) error

// WithAddKclPkg sets the kcl package to add the dependencies.
func WithAddKclPkg(kclPkg pkg.KclPkg) AddOption {
	return func(ao *AddOptions) error {
		ao.KclPkg = kclPkg
		return nil
	}
}

func WithAddSourceUrl(url string) AddOption {
	return func(ao *AddOptions) error {
		if ao.Sources == nil {
			ao.Sources = make([]*downloader.Source, 0)
		}
		source, err := downloader.NewSourceFromStr(url)
		if err != nil {
			return err
		}
		ao.Sources = append(ao.Sources, source)
		return nil
	}
}

// WithAddSource sets the source for the dependency.
func WithAddSource(source *downloader.Source) AddOption {
	return func(ao *AddOptions) error {
		if ao.Sources == nil {
			ao.Sources = make([]*downloader.Source, 0)
		}
		ao.Sources = append(ao.Sources, source)
		return nil
	}
}

// WithAddSources sets the sources for the dependencies.
func WithAddSources(sources []*downloader.Source) AddOption {
	return func(ao *AddOptions) error {
		ao.Sources = sources
		return nil
	}
}

func (c *KpmClient) Add(options ...AddOption) error {
	opts := &AddOptions{}
	for _, option := range options {
		if err := option(opts); err != nil {
			return err
		}
	}

	visitorSelector := func(source *downloader.Source) (visitor.Visitor, error) {
		pkgVisitor := &visitor.PkgVisitor{
			Settings:  &c.settings,
			LogWriter: c.logWriter,
		}

		if source.IsRemote() {
			return &visitor.RemoteVisitor{
				PkgVisitor:            pkgVisitor,
				Downloader:            c.DepDownloader,
				InsecureSkipTLSverify: c.insecureSkipTLSverify,
				EnableCache:           true,
				CachePath:             c.homePath,
				VisitedPath:           c.homePath,
			}, nil
		} else if source.IsLocalTarPath() || source.IsLocalTgzPath() {
			return visitor.NewArchiveVisitor(pkgVisitor), nil
		} else if source.IsLocalPath() {
			rootPath, err := source.FindRootPath()
			if err != nil {
				return nil, err
			}
			kclmodpath := filepath.Join(rootPath, constants.KCL_MOD)
			if utils.DirExists(kclmodpath) {
				return pkgVisitor, nil
			} else {
				return visitor.NewVirtualPkgVisitor(pkgVisitor), nil
			}
		} else {
			return nil, fmt.Errorf("unsupported source")
		}
	}

	addedPkg := &opts.KclPkg

	for _, depSource := range opts.Sources {
		// Set the default OCI registry and repo if the source is nil and the package spec is not nil.
		if depSource.IsNilSource() && !depSource.PkgSpec.IsNil() {
			depSource.Oci = &downloader.Oci{
				Reg:  c.GetSettings().Conf.DefaultOciRegistry,
				Repo: utils.JoinPath(c.GetSettings().Conf.DefaultOciRepo, depSource.PkgSpec.Name),
				Tag:  depSource.PkgSpec.Version,
			}
		}
		visitor, err := visitorSelector(depSource)
		if err != nil {
			return err
		}

		err = visitor.Visit(depSource, func(depPkg *pkg.KclPkg) error {
			dep := pkg.Dependency{
				Name:          depPkg.ModFile.Pkg.Name,
				FullName:      depPkg.GetPkgFullName(),
				Version:       depPkg.ModFile.Pkg.Version,
				LocalFullPath: depPkg.HomePath,
				Source:        *depSource,
			}

			sum, err := c.AcquireDepSum(dep)
			if err != nil {
				return err
			}
			dep.Sum = sum

			if modExistDep, ok := addedPkg.ModFile.Dependencies.Deps.Get(dep.Name); ok {
				if less, err := modExistDep.VersionLessThan(&dep); less && err == nil {
					addedPkg.ModFile.Dependencies.Deps.Set(dep.Name, dep)
				}
			} else {
				addedPkg.ModFile.Dependencies.Deps.Set(dep.Name, dep)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	_, err := c.Update(
		WithUpdatedKclPkg(addedPkg),
	)

	if err != nil {
		return err
	}

	return nil
}
