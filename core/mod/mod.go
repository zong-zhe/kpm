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
	kclModPath := filepath.Join(kclPkg.HomePath, kclMod)
	_, err := os.Stat(kclModPath)
	if os.IsNotExist(err) {
		reporter.Report("kpm: creating new kcl.mod:", kclModPath)
		var buf bytes.Buffer
		err := toml.NewEncoder(&buf).Encode(kclPkg)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(kclModPath, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("kpm: '%s' already exists", kclModPath)
}
