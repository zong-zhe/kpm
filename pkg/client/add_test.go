package client

import (
	"path/filepath"
	"testing"

	pkg "kcl-lang.io/kpm/pkg/package"
)

func TestAdd(t *testing.T) {
	testDir := getTestDir("add_with_pkg_spec")

	testOciDir := filepath.Join(testDir, "oci")

	kpmcli, err := NewKpmClient()
	if err != nil {
		t.Fatal(err)
	}

	kpkg, err := pkg.LoadKclPkgWithOpts(
		pkg.WithPath(testOciDir),
		pkg.WithSettings(&kpmcli.settings),
	)

	if err != nil {
		t.Fatal(err)
	}

	err = kpmcli.Add(
		WithAddKclPkg(*kpkg),
		WithAddSourceUrl("oci://ghcr.io/kcl-lang/helloworld?tag=0.1.2"),
	)

	if err != nil {
		t.Fatal(err)
	}
}
