package client

import (
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
	"kcl-lang.io/kpm/pkg/errors"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/utils"
)

func (c *KpmClient) vendorDeps(kclPkg *pkg.KclPkg, vendorPath string) error {
	lockDeps := make([]pkg.Dependency, 0, kclPkg.Dependencies.Deps.Len())
	for _, k := range kclPkg.Dependencies.Deps.Keys() {
		d, _ := kclPkg.Dependencies.Deps.Get(k)
		lockDeps = append(lockDeps, d)
	}

	// Traverse all dependencies in kcl.mod.lock.
	for i := 0; i < len(lockDeps); i++ {
		d := lockDeps[i]
		if len(d.Name) == 0 {
			return errors.InvalidDependency
		}
		// If the dependency is from the local path, do not vendor it, vendor its dependencies.
		if d.IsFromLocal() {
			dpkg, err := c.LoadPkgFromPath(d.GetLocalFullPath(kclPkg.HomePath))
			if err != nil {
				return err
			}
			err = c.vendorDeps(dpkg, vendorPath)
			if err != nil {
				return err
			}
			continue
		} else {
			vendorFullPath := filepath.Join(vendorPath, d.GenPathSuffix())

			// If the package already exists in the 'vendor', do nothing.
			if utils.DirExists(vendorFullPath) {
				d.LocalFullPath = vendorFullPath
				lockDeps[i] = d
				continue
			} else {
				// If not in the 'vendor', check the global cache.
				cacheFullPath := c.getDepStorePath(c.homePath, &d, false)
				if utils.DirExists(cacheFullPath) {
					// If there is, copy it into the 'vendor' directory.
					err := copy.Copy(cacheFullPath, vendorFullPath)
					if err != nil {
						return err
					}
				} else {
					// re-download if not.
					err := c.AddDepToPkg(kclPkg, &d)
					if err != nil {
						return err
					}
					// re-vendor again with new kcl.mod and kcl.mod.lock
					err = c.vendorDeps(kclPkg, vendorPath)
					if err != nil {
						return err
					}
					return nil
				}
			}

			if d.GetPackage() != "" {
				tempVendorFullPath, err := utils.FindPackage(vendorFullPath, d.GetPackage())
				if err != nil {
					return err
				}
				vendorFullPath = tempVendorFullPath
			}

			dpkg, err := c.LoadPkgFromPath(vendorFullPath)
			if err != nil {
				return err
			}

			// Vendor the dependencies of the current dependency.
			err = c.vendorDeps(dpkg, vendorPath)
			if err != nil {
				return err
			}
			d.LocalFullPath = vendorFullPath
			lockDeps[i] = d
		}
	}

	// Update the dependencies in kcl.mod.lock.
	for _, d := range lockDeps {
		kclPkg.Dependencies.Deps.Set(d.Name, d)
	}

	return nil
}

// VendorDeps will vendor all the dependencies of the current kcl package.
func (c *KpmClient) VendorDeps(kclPkg *pkg.KclPkg) error {
	// Mkdir the dir "vendor".
	vendorPath := kclPkg.LocalVendorPath()
	err := os.MkdirAll(vendorPath, 0755)
	if err != nil {
		return err
	}

	return c.vendorDeps(kclPkg, vendorPath)
}
