package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"kcl-lang.io/kpm/pkg/features"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/settings"
	"kcl-lang.io/kpm/pkg/utils"
)

func TestWithMVS(t *testing.T) {
	features.Enable(features.SupportMVS)
	defer features.Disable(features.SupportMVS)

	testUpdate(t)
	testVendorWithMVS(t)
}

func testUpdate(t *testing.T) {
	testDir := getTestDir("test_update_with_mvs")
	kpmcli, err := NewKpmClient()
	if err != nil {
		t.Fatal(err)
	}

	updates := []struct {
		name   string
		before func() error
	}{
		{
			name:   "update_0",
			before: func() error { return nil },
		},
		{
			name: "update_1",
			before: func() error {
				if err := copy.Copy(filepath.Join(testDir, "update_1", "pkg", "kcl.mod.bk"), filepath.Join(testDir, "update_1", "pkg", "kcl.mod")); err != nil {
					return err
				}
				if err := copy.Copy(filepath.Join(testDir, "update_1", "pkg", "kcl.mod.lock.bk"), filepath.Join(testDir, "update_1", "pkg", "kcl.mod.lock")); err != nil {
					return err
				}
				return nil
			},
		},
	}

	for _, update := range updates {
		if err := update.before(); err != nil {
			t.Fatal(err)
		}

		kpkg, err := kpmcli.LoadPkgFromPath(filepath.Join(testDir, update.name, "pkg"))
		if err != nil {
			t.Fatal(err)
		}

		_, err = kpmcli.Update(WithUpdatedKclPkg(kpkg))
		if err != nil {
			t.Fatal(err)
		}

		expectedMod, err := os.ReadFile(filepath.Join(testDir, update.name, "pkg", "kcl.mod.expect"))
		if err != nil {
			t.Fatal(err)
		}

		expectedModLock, err := os.ReadFile(filepath.Join(testDir, update.name, "pkg", "kcl.mod.lock.expect"))
		if err != nil {
			t.Fatal(err)
		}

		gotMod, err := os.ReadFile(filepath.Join(testDir, update.name, "pkg", "kcl.mod"))
		if err != nil {
			t.Fatal(err)
		}

		gotModLock, err := os.ReadFile(filepath.Join(testDir, update.name, "pkg", "kcl.mod.lock"))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, utils.RmNewline(string(expectedMod)), utils.RmNewline(string(gotMod)))
		assert.Equal(t, utils.RmNewline(string(expectedModLock)), utils.RmNewline(string(gotModLock)))
	}
}

func testVendorWithMVS(t *testing.T) {
	testDir := getTestDir("test_vendor_mvs")
	pkgPath := filepath.Join(testDir, "pkg")
	kPkg, err := pkg.LoadKclPkgWithOpts(
		pkg.WithPath(pkgPath),
		pkg.WithSettings(settings.GetSettings()),
	)
	assert.Equal(t, err, nil)

	kpmcli, err := NewKpmClient()
	assert.Equal(t, err, nil)
	err = kpmcli.VendorDeps(kPkg)
	assert.Equal(t, err, nil)

	assert.Equal(t, utils.DirExists(filepath.Join(pkgPath, "vendor")), true)
	assert.Equal(t, utils.DirExists(filepath.Join(pkgPath, "vendor", "helloworld_0.1.2")), true)
	assert.Equal(t, utils.DirExists(filepath.Join(pkgPath, "vendor", "helloworld_0.1.1")), false)
}
