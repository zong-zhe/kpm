package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"kusionstack.io/kpm/core/conf"
	"kusionstack.io/kpm/core/git"
	"kusionstack.io/kpm/core/mod"
	"kusionstack.io/kpm/core/reporter"
	"kusionstack.io/kpm/core/utils"
)

type KclPkg struct {
	modFile  mod.ModFile
	lockFile mod.ModLockFile
	HomePath string
}

func NewKclPkg(conf conf.Config) KclPkg {
	return KclPkg{
		modFile:  *mod.NewModFile(&conf, filepath.Join(conf.KclModPath, mod.File)),
		lockFile: *mod.NewModLockFile(&conf, filepath.Join(conf.KclModPath, mod.LockFile)),
		HomePath: conf.KclModPath,
	}
}

func LoadKclPkg(conf *conf.Config) (*KclPkg, error) {
	modFile, err := mod.LoadModFile(conf.KclModPath)
	modLockFile, err := mod.LoadModLockFile(conf.KclModPath)
	if err != nil {
		return nil, err
	}
	return &KclPkg{
		modFile:  *modFile,
		lockFile: *modLockFile,
		HomePath: conf.KclModPath,
	}, nil
}

// InitEmptyModule inits an empty kcl module and create a default kcl.mod.
func (kclPkg KclPkg) InitEmptyPkg() error {
	_, err := os.Stat(kclPkg.modFile.HomePath)
	if os.IsNotExist(err) {
		reporter.Report("kpm: creating new kcl.mod:", kclPkg.HomePath)
		err := kclPkg.modFile.Store()
		if err != nil {
			reporter.Report(err)
		}
	}
	_, err = os.Stat(kclPkg.lockFile.HomePath)
	if os.IsNotExist(err) {
		reporter.Report("kpm: creating new kcl.mod.lock:", kclPkg.HomePath)
		err = kclPkg.lockFile.Store()
		if err != nil {
			reporter.Report(err)
		}
	}
	return fmt.Errorf("kpm: '%s' already exists", kclPkg.modFile.HomePath)
}

// InitEmptyModule inits an empty kcl module and create a default kcl.mod.
func (kclPkg KclPkg) AddDeps(gitOpts []git.GitOption, localpath string) error {

	for _, opt := range gitOpts {
		d := mod.ParseOpt(&opt)
		if !mod.DepEqual(kclPkg.modFile.Dependencies.Deps[d.Name], d) {
			// the dep passed on the cli is different from the jsonnetFile
			kclPkg.modFile.Dependencies.Deps[d.Name] = d

			// we want to install the passed version (ignore the lock)
			delete(kclPkg.lockFile.Dependencies.Deps, d.Name)
		}
	}

	changedDeps, _ := getDeps(kclPkg.modFile.Dependencies, kclPkg.lockFile.Dependencies, localpath)

	fmt.Println(changedDeps)

	for k, v := range changedDeps.Deps {
		kclPkg.modFile.Dependencies.Deps[k] = v
		kclPkg.lockFile.Dependencies.Deps[k] = v
	}

	// store kcl.mod.lock
	// 这里只有新加入的，需要增量写入
	kclPkg.modFile.Store()
	kclPkg.lockFile.Store()

	return fmt.Errorf("kpm: '%s' already exists", kclPkg.modFile.HomePath)
}

func getDeps(deps mod.Dependencies, lockDeps mod.Dependencies, localPath string) (*mod.Dependencies, error) {
	newDeps := mod.Dependencies{
		Deps: make(map[string]mod.Dependency),
	}

	for _, d := range deps.Deps {
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
			return nil, nil
		}
		if expectedSum != "" && lockedDep.Sum != expectedSum {
			return nil, fmt.Errorf("checksum mismatch")
		}
		newDeps.Deps[d.Name] = *lockedDep
		// we settled on a new version, add it to the locks for recursion
		lockDeps.Deps[d.Name] = *lockedDep
	}

	for _, d := range newDeps.Deps {

		modfile, err := mod.LoadModFile(filepath.Join(localPath, d.Name))
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

// check returns whether the files present at the vendor/ folder match the
// sha256 sum of the package. local-directory dependencies are not checked as
// their purpose is to change during development where integrity checking would
// be a hindrance.
func check(d mod.Dependency, vendorDir string) bool {
	// assume a local dependency is intact as long as it exists
	// if d.Source.LocalSource != nil {
	// 	x, err := mod.ModFileExists(filepath.Join(vendorDir, d.Name()))
	// 	if err != nil {
	// 		return false
	// 	}
	// 	return x
	// }

	if d.Sum == "" {
		// no sum available, need to download
		return false
	}

	dir := filepath.Join(vendorDir, d.Name)
	sum := utils.HashDir(dir)
	return d.Sum == sum
}

// func LoadKclPkg(homePath string) (*KclPkg, error) {

// 	if err != nil {
// 		return nil, err
// 	}
// 	return kclPkg, nil
// }

// func genKclModLock(kclPkg KclPkg) error {
// 	return nil
// }

// func (kclPkg KclPkg) AddDeps(conf *conf.Config, dep *pkg.Dependency) error {
// 	if kclPkg.ContainsDepNamed(dep.GetName()) {
// 		reporter.Report("kpm: '", dep.GetName(), "' has already exists.")
// 	}

// 	_, err := newpkg.GetDeps(directDeps, localPath, oldLocks)

// 	kclPkg.Pkg.Dependencies[dep.GetName()] = dep
// 	err = genKclMod(kclPkg)

// 	if err != nil {
// 		reporter.ExitWithReport("kpm: failed to update ", kclPkg.HomePath)
// 	}

// 	reporter.Report("kpm: '", dep.GetName(), "' added successfully.")

// 	// genKclModLock(kclPkg)

// 	return nil
// }

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"

// 	kclmodfile "kusionstack.io/kpm/core/modfile"
// )

// func GetDeps(directDeps map[string]kclmodfile.Dependency, localPath string, locks map[string]kclmodfile.Dependency) (map[string]kclmodfile.Dependency, error) {
// 	deps := make(map[string]kclmodfile.Dependency)

// 	for _, d := range directDeps {
// 		l, present := locks[d.Name]

// 		// already locked and the integrity is intact
// 		if present {
// 			d.Version = locks[d.Name].Version

// 			if check(l, localPath) {
// 				deps[d.Name] = l
// 				continue
// 			}
// 		}
// 		expectedSum := locks[d.Name].Sum

// 		// either not present or not intact: download again
// 		dir := filepath.Join(localPath, d.Name)
// 		os.RemoveAll(dir)

// 		locked, err := d.Download(localPath)
// 		if err != nil {
// 			return nil, fmt.Errorf("err")
// 		}
// 		if expectedSum != "" && locked.Sum != expectedSum {
// 			return nil, fmt.Errorf("checksum mismatch for %s. Expected %s but got %s", d.Name, expectedSum, locked.Sum)
// 		}
// 		deps[d.Name] = *locked
// 		// we settled on a new version, add it to the locks for recursion
// 		locks[d.Name] = *locked
// 	}

// 	for _, d := range deps {

// 		f, err := kclmodfile.LoadModFile(filepath.Join(localPath, d.Name, kclmodfile.File))
// 		if err != nil {
// 			if os.IsNotExist(err) {
// 				continue
// 			}
// 			return nil, err
// 		}

// 		nested, err := GetDeps(f.Dependencies, localPath, locks)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, d := range nested {
// 			if _, ok := deps[d.Name]; !ok {
// 				deps[d.Name] = d
// 			}
// 		}
// 	}

// 	return deps, nil
// }

// func check(d kclmodfile.Dependency, vendorDir string) bool {

// 	if d.Sum == "" {
// 		// no sum available, need to download
// 		return false
// 	}

// 	dir := filepath.Join(vendorDir, d.Name)
// 	sum := kclmodfile.HashDir(dir)
// 	return d.Sum == sum
// }
