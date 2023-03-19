// Copyright 2022 The KCL Authors. All rights reserved.

package ops

import (
	conf "kusionstack.io/kpm/core/conf"
	"kusionstack.io/kpm/gen/pkg"
)

// KpmInit initializes an empty kcl module.
func KpmAdd(conf *conf.Config, dep *pkg.Dependency) error {
	return nil
	// kclPkg, err := mod.LoadKclPkg(conf.KclModPath)
	// if err != nil {
	// 	reporter.ExitWithReport("kpm: failed to load kcl.mod from", conf.KclModPath)
	// }

	// return kclPkg.AddDeps(conf, dep)
}
