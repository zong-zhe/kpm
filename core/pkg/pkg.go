package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	kclmodfile "kusionstack.io/kpm/core/modfile"
)

func GetDeps(directDeps map[string]kclmodfile.Dependency, localPath string, locks map[string]kclmodfile.Dependency) (map[string]kclmodfile.Dependency, error) {
	deps := make(map[string]kclmodfile.Dependency)

	for _, d := range directDeps {
		l, present := locks[d.Name]

		// already locked and the integrity is intact
		if present {
			d.Version = locks[d.Name].Version

			if check(l, localPath) {
				deps[d.Name] = l
				continue
			}
		}
		expectedSum := locks[d.Name].Sum

		// either not present or not intact: download again
		dir := filepath.Join(localPath, d.Name)
		os.RemoveAll(dir)

		locked, err := d.Download(localPath)
		if err != nil {
			return nil, fmt.Errorf("err")
		}
		if expectedSum != "" && locked.Sum != expectedSum {
			return nil, fmt.Errorf("checksum mismatch for %s. Expected %s but got %s", d.Name, expectedSum, locked.Sum)
		}
		deps[d.Name] = *locked
		// we settled on a new version, add it to the locks for recursion
		locks[d.Name] = *locked
	}

	for _, d := range deps {

		f, err := kclmodfile.LoadModFile(filepath.Join(localPath, d.Name, kclmodfile.File))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		nested, err := GetDeps(f.Dependencies, localPath, locks)
		if err != nil {
			return nil, err
		}

		for _, d := range nested {
			if _, ok := deps[d.Name]; !ok {
				deps[d.Name] = d
			}
		}
	}

	return deps, nil
}

func check(d kclmodfile.Dependency, vendorDir string) bool {

	if d.Sum == "" {
		// no sum available, need to download
		return false
	}

	dir := filepath.Join(vendorDir, d.Name)
	sum := kclmodfile.HashDir(dir)
	return d.Sum == sum
}
