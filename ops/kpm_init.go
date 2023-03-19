// Copyright 2022 The KCL Authors. All rights reserved.

package ops

import (
	conf "kusionstack.io/kpm/core/conf"
)

// KpmInit initializes an empty kcl module.
func KpmInit(conf conf.Config) error {
	return nil
	// return mod.NewKclPkg(conf).InitEmptyModule()
}
