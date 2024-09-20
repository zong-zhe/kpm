package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kcl-lang.io/kpm/pkg/downloader"
	"kcl-lang.io/kpm/pkg/reporter"
	"kcl-lang.io/kpm/pkg/settings"
	"kcl-lang.io/kpm/pkg/utils"
)

type LoadOptions struct {
	// The package path.
	PkgPath string
	// The settings with default oci registry.
	Settings *settings.Settings
}

type LoadOption func(*LoadOptions)

// WithPkgPath sets the package path.
func WithPkgPath(pkgPath string) LoadOption {
	return func(opts *LoadOptions) {
		opts.PkgPath = pkgPath
	}
}

// WithSettings sets the settings with default oci registry.
func WithSettings(settings *settings.Settings) LoadOption {
	return func(opts *LoadOptions) {
		opts.Settings = settings
	}
}

// LoadKclPkgWithOpts loads a package from the file system with options.
// The options include the package path and the settings with default oci registry.
func LoadKclPkgWithOpts(options ...LoadOption) (*KclPkg, error) {
	opts := &LoadOptions{}
	for _, opt := range options {
		opt(opts)
	}

	pkgPath := opts.PkgPath

	modFile := new(ModFile)
	err := modFile.LoadModFile(filepath.Join(pkgPath, MOD_FILE))
	if err != nil {
		return nil, fmt.Errorf("could not load 'kcl.mod' in '%s': %w", pkgPath, err)
	}

	// load the kcl.mod.lock file.
	// Get dependencies from kcl.mod.lock.
	deps, err := LoadLockDeps(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("could not load 'kcl.mod' in '%s': %w", pkgPath, err)
	}

	// pre-process the package.
	// 1. Transform the local path to the absolute path.
	err = convertDepsLocalPathToAbsPath(&modFile.Dependencies, pkgPath)
	if err != nil {
		return nil, fmt.Errorf("could not load 'kcl.mod' in '%s': %w", pkgPath, err)
	}
	// 2. Fill the default oci registry, the default oci registry is in the settings.
	err = fillDepsInfoWithSettings(&modFile.Dependencies, opts.Settings)
	if err != nil {
		return nil, fmt.Errorf("could not load 'kcl.mod' in '%s': %w", pkgPath, err)
	}
	// 3. Sync the dependencies information in kcl.mod.lock with the dependencies in kcl.mod.
	for _, name := range modFile.Dependencies.Deps.Keys() {
		modDep, ok := modFile.Dependencies.Deps.Get(name)
		if !ok {
			return nil, fmt.Errorf("could not load 'kcl.mod' in '%s': %w", pkgPath, err)
		}
		if lockDep, ok := deps.Deps.Get(name); !ok {
			lockDep.Source = modDep.Source
			lockDep.FullName = modDep.GenDepFullName()
			deps.Deps.Set(name, lockDep)
		}
	}

	return &KclPkg{
		ModFile:      *modFile,
		HomePath:     pkgPath,
		Dependencies: *deps,
	}, nil
}

// LoadAndFillModFileWithOpts loads a mod file from the file system with options.
// It will load the mod file, convert the local path to the absolute path, and fill the default oci registry.
func LoadAndFillModFileWithOpts(options ...LoadOption) (*ModFile, error) {
	opts := &LoadOptions{}
	for _, opt := range options {
		opt(opts)
	}

	pkgPath := opts.PkgPath

	// Load the mod file.
	// The content of the `ModFile` is the same as the content in kcl.mod
	// The `ModFile` lacks some information of the dependencies.
	modFile := new(ModFile)
	err := modFile.LoadModFile(filepath.Join(pkgPath, MOD_FILE))
	if err != nil {
		return nil, fmt.Errorf("could not load 'kcl.mod' in '%s': %w", pkgPath, err)
	}

	// pre-process the package.
	// 1. Transform the local path to the absolute path.
	err = convertDepsLocalPathToAbsPath(&modFile.Dependencies, pkgPath)
	if err != nil {
		return nil, fmt.Errorf("could not load 'kcl.mod' in '%s': %w", pkgPath, err)
	}
	// 2. Fill the default oci registry, the default oci registry is in the settings.
	err = fillDepsInfoWithSettings(&modFile.Dependencies, opts.Settings)
	if err != nil {
		return nil, fmt.Errorf("could not load 'kcl.mod' in '%s': %w", pkgPath, err)
	}

	return modFile, nil
}

// `convertDepsLocalPathToAbsPath` will transform the local path to the absolute path from `rootPath` in dependencies.
func convertDepsLocalPathToAbsPath(deps *Dependencies, rootPath string) error {
	for _, name := range deps.Deps.Keys() {
		dep, ok := deps.Deps.Get(name)
		if !ok {
			break
		}
		// Transform the local path to the absolute path.
		if dep.Local != nil {
			var localFullPath string
			var err error
			if filepath.IsAbs(dep.Local.Path) {
				localFullPath = dep.Local.Path
			} else {
				localFullPath, err = filepath.Abs(filepath.Join(rootPath, dep.Local.Path))
				if err != nil {
					return fmt.Errorf("failed to get the absolute path of the local dependency %s: %w", name, err)
				}
			}
			dep.LocalFullPath = localFullPath
		}
		deps.Deps.Set(name, dep)
	}

	return nil
}

// `fillDepsInfoWithSettings` will fill the default oci registry info in dependencies.
func fillDepsInfoWithSettings(deps *Dependencies, settings *settings.Settings) error {
	for _, name := range deps.Deps.Keys() {
		dep, ok := deps.Deps.Get(name)
		if !ok {
			break
		}
		// Fill the default oci registry.
		if dep.Source.Oci != nil {
			if len(dep.Source.Oci.Reg) == 0 {
				dep.Source.Oci.Reg = settings.DefaultOciRegistry()
			}

			if len(dep.Source.Oci.Repo) == 0 {
				urlpath := utils.JoinPath(settings.DefaultOciRepo(), dep.Name)
				dep.Source.Oci.Repo = urlpath
			}
		}
		if dep.Source.Registry != nil {
			if len(dep.Source.Registry.Reg) == 0 {
				dep.Source.Registry.Reg = settings.DefaultOciRegistry()
			}

			if len(dep.Source.Registry.Repo) == 0 {
				urlpath := utils.JoinPath(settings.DefaultOciRepo(), dep.Name)
				dep.Source.Registry.Repo = urlpath
			}

			dep.Version = dep.Source.Registry.Version
		}
		if dep.Source.IsNilSource() {
			dep.Source.Registry = &downloader.Registry{
				Name:    dep.Name,
				Version: dep.Version,
				Oci: &downloader.Oci{
					Reg:  settings.DefaultOciRegistry(),
					Repo: utils.JoinPath(settings.DefaultOciRepo(), dep.Name),
					Tag:  dep.Version,
				},
			}
		}
		deps.Deps.Set(name, dep)
	}

	return nil
}

// LoadKclPkg will load a package from the 'pkgPath'
// The default oci registry in '$KCL_PKG_PATH/.kpm/config/kpm.json' will be used.
func LoadKclPkg(pkgPath string) (*KclPkg, error) {
	return LoadKclPkgWithOpts(WithPkgPath(pkgPath), WithSettings(settings.GetSettings()))
}

// FindFirstKclPkgFrom will find the first kcl package from the 'path'
// The default oci registry in '$KCL_PKG_PATH/.kpm/config/kpm.json' will be used.
func FindFirstKclPkgFrom(pkgpath string) (*KclPkg, error) {
	matches, _ := filepath.Glob(filepath.Join(pkgpath, "*.tar"))
	if matches == nil || len(matches) != 1 {
		// then try to glob tgz file
		matches, _ = filepath.Glob(filepath.Join(pkgpath, "*.tgz"))
		if matches == nil || len(matches) != 1 {
			pkg, err := LoadKclPkg(pkgpath)
			if err != nil {
				return nil, reporter.NewErrorEvent(
					reporter.InvalidKclPkg,
					err,
					fmt.Sprintf("failed to find the kcl package tar from '%s'.", pkgpath),
				)
			}

			return pkg, nil
		}
	}

	tarPath := matches[0]
	unTarPath := filepath.Dir(tarPath)
	var err error
	if utils.IsTar(tarPath) {
		err = utils.UnTarDir(tarPath, unTarPath)
	} else {
		err = utils.ExtractTarball(tarPath, unTarPath)
	}
	if err != nil {
		return nil, reporter.NewErrorEvent(
			reporter.FailedUntarKclPkg,
			err,
			fmt.Sprintf("failed to untar the kcl package tar from '%s' into '%s'.", tarPath, unTarPath),
		)
	}

	// After untar the downloaded kcl package tar file, remove the tar file.
	if utils.DirExists(tarPath) {
		rmErr := os.Remove(tarPath)
		if rmErr != nil {
			return nil, reporter.NewErrorEvent(
				reporter.FailedUntarKclPkg,
				err,
				fmt.Sprintf("failed to untar the kcl package tar from '%s' into '%s'.", tarPath, unTarPath),
			)
		}
	}

	pkg, err := LoadKclPkg(unTarPath)
	if err != nil {
		return nil, reporter.NewErrorEvent(
			reporter.InvalidKclPkg,
			err,
			fmt.Sprintf("failed to find the kcl package tar from '%s'.", pkgpath),
		)
	}

	return pkg, nil
}

// LoadKclPkgFromTar loads a package *.tar file from the 'pkgTarPath'
// The default oci registry in '$KCL_PKG_PATH/.kpm/config/kpm.json' will be used.
func LoadKclPkgFromTar(pkgTarPath string) (*KclPkg, error) {
	destDir := strings.TrimSuffix(pkgTarPath, filepath.Ext(pkgTarPath))
	err := utils.UnTarDir(pkgTarPath, destDir)
	if err != nil {
		return nil, err
	}
	return LoadKclPkg(destDir)
}
