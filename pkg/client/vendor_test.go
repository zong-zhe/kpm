package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/stretchr/testify/assert"
	"kcl-lang.io/kpm/pkg/downloader"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/utils"
)

func TestVendorDeps(t *testing.T) {
	testDir := getTestDir("resolve_deps")
	kpm_home := filepath.Join(testDir, "kpm_home")
	os.RemoveAll(filepath.Join(testDir, "my_kcl"))
	kcl1Sum, _ := utils.HashDir(filepath.Join(kpm_home, "kcl1"))
	kcl2Sum, _ := utils.HashDir(filepath.Join(kpm_home, "kcl2"))

	depKcl1 := pkg.Dependency{
		Name:     "kcl1",
		FullName: "kcl1",
		Version:  "0.0.1",
		Sum:      kcl1Sum,
		Source: downloader.Source{
			Oci: &downloader.Oci{
				Reg:  "ghcr.io",
				Repo: "kcl-lang/kcl1",
				Tag:  "0.0.1",
			},
		},
	}

	depKcl2 := pkg.Dependency{
		Name:     "kcl2",
		FullName: "kcl2",
		Version:  "0.0.1",
		Sum:      kcl2Sum,
		Source: downloader.Source{
			Oci: &downloader.Oci{
				Reg:  "ghcr.io",
				Repo: "kcl-lang/kcl2",
				Tag:  "0.0.1",
			},
		},
	}

	mppTest := orderedmap.NewOrderedMap[string, pkg.Dependency]()
	mppTest.Set("kcl1", depKcl1)
	mppTest.Set("kcl2", depKcl2)

	kclPkg := pkg.KclPkg{
		ModFile: pkg.ModFile{
			HomePath: filepath.Join(testDir, "my_kcl"),
			// Whether the current package uses the vendor mode
			// In the vendor mode, kpm will look for the package in the vendor subdirectory
			// in the current package directory.
			VendorMode: false,
			Dependencies: pkg.Dependencies{
				Deps: mppTest,
			},
		},
		HomePath: filepath.Join(testDir, "my_kcl"),
		// The dependencies in the current kcl package are the dependencies of kcl.mod.lock,
		// not the dependencies in kcl.mod.
		Dependencies: pkg.Dependencies{
			Deps: mppTest,
		},
	}

	mykclVendorPath := filepath.Join(filepath.Join(testDir, "my_kcl"), "vendor")
	assert.Equal(t, utils.DirExists(mykclVendorPath), false)
	kpmcli, err := NewKpmClient()
	kpmcli.homePath = kpm_home
	assert.Equal(t, err, nil)
	err = kpmcli.VendorDeps(&kclPkg)
	assert.Equal(t, err, nil)
	assert.Equal(t, utils.DirExists(mykclVendorPath), true)
	assert.Equal(t, utils.DirExists(filepath.Join(mykclVendorPath, "kcl1_0.0.1")), true)
	assert.Equal(t, utils.DirExists(filepath.Join(mykclVendorPath, "kcl2_0.0.1")), true)

	maps, err := kpmcli.ResolveDepsIntoMap(&kclPkg)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(maps), 2)

	os.RemoveAll(filepath.Join(testDir, "my_kcl"))
}
