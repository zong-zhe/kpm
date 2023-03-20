// Copyright 2022 The KCL Authors. All rights reserved.

package ops

import (
	"kusionstack.io/kpm/core/opt"
	"kusionstack.io/kpm/core/pkg"
)

// KpmInit initializes an empty kcl module.
func KpmInit(opt *opt.InitOptions) error {
	return pkg.NewKclPkg(opt).InitEmptyPkg()
}
