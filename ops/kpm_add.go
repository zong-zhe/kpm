// Copyright 2022 The KCL Authors. All rights reserved.

package ops

import (
	"os"

	mod "kusionstack.io/kpm/core/mod"
	"kusionstack.io/kpm/core/reporter"
	"kusionstack.io/kpm/gen/pkg"
)

// KpmInit initializes an empty kcl module.
func KpmAdd(dep *pkg.Dependency) error {
	pwd, err := os.Getwd()
	kclPkg, err := mod.LoadKclPkg(pwd)
	if err != nil {
		reporter.ExitWithReport("kpm: failed to load kcl.mod from", kclPkg.HomePath)
	}

	return kclPkg.AddDeps(dep)
}
