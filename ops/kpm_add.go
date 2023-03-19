// Copyright 2022 The KCL Authors. All rights reserved.

package ops

import (
	"fmt"

	conf "kusionstack.io/kpm/core/conf"
	"kusionstack.io/kpm/core/git"
	"kusionstack.io/kpm/core/pkg"
)

// KpmInit initializes an empty kcl module.
func KpmAdd(conf *conf.Config, opts *git.GitOption) error {
	kclPkg, err := pkg.LoadKclPkg(conf)

	if err != nil {
		fmt.Println(err)
	}

	return kclPkg.AddDeps([]git.GitOption{*opts}, conf.KclModPath)
}
