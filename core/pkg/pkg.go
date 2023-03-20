package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	modfile "kusionstack.io/kpm/core/mod"
	"kusionstack.io/kpm/core/opt"
	"kusionstack.io/kpm/core/reporter"
	"kusionstack.io/kpm/core/utils"
)

type KclPkg struct {
	modFile  modfile.ModFile
	lockFile modfile.ModLockFile
	HomePath string
}

func NewKclPkg(opt *opt.InitOptions) KclPkg {
	return KclPkg{
		modFile:  *modfile.NewModFile(opt, filepath.Join(opt.InitPath, modfile.File)),
		lockFile: *modfile.NewModLockFile(opt, filepath.Join(opt.InitPath, modfile.LockFile)),
		HomePath: opt.InitPath,
	}
}

func LoadKclPkg(pkgPath string) (*KclPkg, error) {
	modFile, err := modfile.LoadModFile(pkgPath)
	modLockFile, err := modfile.LoadModLockFile(pkgPath)
	if err != nil {
		return nil, err
	}
	return &KclPkg{
		modFile:  *modFile,
		lockFile: *modLockFile,
		HomePath: pkgPath,
	}, nil
}

// InitEmptyModule inits an empty kcl module and create a default kcl.modfile.
func (kclPkg KclPkg) InitEmptyPkg() error {
	err := createFileIfNotExist(kclPkg.modFile.HomePath, "kcl.mod", kclPkg.modFile.Store)
	if err != nil {
		return err
	}

	err = createFileIfNotExist(kclPkg.lockFile.HomePath, "kcl.modfile.lock", kclPkg.lockFile.Store)
	if err != nil {
		return err
	}

	return nil
}

func createFileIfNotExist(filePath string, fileName string, storeFunc func() error) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		reporter.Report("kpm: creating new "+fileName+":", filePath)
		err := storeFunc()
		if err != nil {
			reporter.Report("kpm: failed to create "+fileName+",", err)
			return err
		}
	} else {
		reporter.Report("kpm: '%s' already exists", filePath)
		return err
	}
	return nil
}

// InitEmptyModule inits an empty kcl module and create a default kcl.modfile.
func (kclPkg KclPkg) AddDeps(opt *opt.AddOptions) error {

	d := modfile.ParseOpt(&opt.RegistryOpts)
	if !modfile.DepEqual(kclPkg.modFile.Dependencies.Deps[d.Name], *d) {
		// the dep passed on the cli is different from the jsonnetFile
		kclPkg.modFile.Dependencies.Deps[d.Name] = *d

		// we want to install the passed version (ignore the lock)
		delete(kclPkg.lockFile.Dependencies.Deps, d.Name)
	}

	changedDeps, err := getDeps(kclPkg.modFile.Dependencies, kclPkg.lockFile.Dependencies, opt.LocalPath)

	if err != nil {
		reporter.ExitWithReport("kpm: failed to download dependancies.")
	}

	for k, v := range changedDeps.Deps {
		kclPkg.modFile.Dependencies.Deps[k] = v
		kclPkg.lockFile.Dependencies.Deps[k] = v
	}

	err = kclPkg.modFile.Store()
	if err != nil {
		return err
	}
	err = kclPkg.lockFile.Store()
	if err != nil {
		return err
	}

	return nil
}

func getDeps(deps modfile.Dependencies, lockDeps modfile.Dependencies, localPath string) (*modfile.Dependencies, error) {
	newDeps := modfile.Dependencies{
		Deps: make(map[string]modfile.Dependency),
	}

	for _, d := range deps.Deps {
		if len(d.Name) == 0 {
			reporter.ExitWithReport("kpm: invalid dependencies.")
			return nil, fmt.Errorf("kpm: invalid dependencies.")
		}
		l, present := lockDeps.Deps[d.Name]

		// already locked and the integrity is intact
		if present {
			d.Version = lockDeps.Deps[d.Name].Version

			if check(l, localPath) {
				newDeps.Deps[d.Name] = l
				continue
			}
		}
		expectedSum := lockDeps.Deps[d.Name].Sum

		dir := filepath.Join(localPath, d.Name)
		os.RemoveAll(dir)

		lockedDep, err := d.Download(dir)
		if err != nil {
			return nil, fmt.Errorf("checksum mismatch")
		}
		if expectedSum != "" && lockedDep.Sum != expectedSum {
			return nil, fmt.Errorf("checksum mismatch")
		}
		newDeps.Deps[d.Name] = *lockedDep
		lockDeps.Deps[d.Name] = *lockedDep
	}

	for _, d := range newDeps.Deps {
		modfile, err := modfile.LoadModFile(filepath.Join(localPath, d.Name))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		nested, err := getDeps(modfile.Dependencies, lockDeps, localPath)
		if err != nil {
			return nil, err
		}

		for _, d := range nested.Deps {
			if _, ok := newDeps.Deps[d.Name]; !ok {
				newDeps.Deps[d.Name] = d
			}
		}
	}

	return &newDeps, nil
}

func check(d modfile.Dependency, vendorDir string) bool {

	if d.Sum == "" {
		return false
	}

	dir := filepath.Join(vendorDir, d.Name)
	sum := utils.HashDir(dir)
	return d.Sum == sum
}
