// Copyright 2022 The KCL Authors. All rights reserved.

package ops

import (
	"net/url"

	"kusionstack.io/kpm/core/reporter"
	"kusionstack.io/kpm/core/store"
	"kusionstack.io/kpm/gen/pkg"
)

// KpmInit initializes an empty kcl module.
func KpmAddGit(gitDep *pkg.GitDependency) error {

	// 在这里拼出来url

	gitUrl, err := url.Parse(depUrl)

	gitUrl.

	gitQuery := u.Query()
	gitQuery.Add(paramName, paramValue)
 u.RawQuery = q.Encode()

	if err != nil {
		reporter.Report("kpm: Invalid git url")
	}

	// 在这里parse出来名字

	store.GetFromGit(name, gitUrl)
	return
}
