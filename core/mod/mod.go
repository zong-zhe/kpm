// Copyright 2022 The KCL Authors. All rights reserved.

// Package mod is the core abstraction in kpm for working with a package or module of kcl.
package mod

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	conf "kusionstack.io/kpm/core/conf"
	"kusionstack.io/kpm/core/reporter"
	"kusionstack.io/kpm/gen/pkg"
)

const kclMod = "kcl.mod"

// KclPkg is the package of kcl.
type KclPkg struct {
	HomePath string       `toml:"-"`
	Pkg      *pkg.Package `toml:"package"`
}

// NewKclPkg new a kcl package `KclPkg`` based on `Config`.
func NewKclPkg(conf conf.Config) KclPkg {
	return KclPkg{
		HomePath: conf.ExecPath,
		Pkg: &pkg.Package{
			Name:    conf.Name,
			Version: conf.Version,
			Edition: conf.Edition,
		},
	}
}

// InitEmptyModule inits an empty kcl module and create a default kcl.mod.
func (kclPkg KclPkg) InitEmptyModule() error {
	kclPkg.HomePath = filepath.Join(kclPkg.HomePath, kclMod)
	_, err := os.Stat(kclPkg.HomePath)
	if os.IsNotExist(err) {
		reporter.Report("kpm: creating new kcl.mod:", kclPkg.HomePath)
		genKclMod(kclPkg)
	}
	return fmt.Errorf("kpm: '%s' already exists", kclPkg.HomePath)
}

func (kclPkg KclPkg) ContainsDepNamed(name string) bool {
	_, ok := kclPkg.Pkg.Dependencies[name]
	return ok
}

func genKclMod(kclPkg KclPkg) error {
	var buf bytes.Buffer
	err := toml.NewEncoder(&buf).Encode(kclPkg)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(kclPkg.HomePath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (kclPkg KclPkg) AddDeps(dep *pkg.Dependency) error {

	if kclPkg.ContainsDepNamed(dep.GetName()) {
		reporter.Report("kpm: ", dep.GetName(), "has already exists.")
	}

	genKclMod(kclPkg)

	findOrDownloadDeps(dep)

	genKclModLock(dep)

	return nil
}


func findOrDownloadDeps(deps []*pkg.Dependency) error {
	return nil
}

func genKclModLock(deps []*pkg.Dependency) error {
	return nil
}
