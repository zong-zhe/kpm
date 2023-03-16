// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"github.com/urfave/cli/v2"
	reporter "kusionstack.io/kpm/core/reporter"
	"kusionstack.io/kpm/gen/pkg"
	"kusionstack.io/kpm/ops"
)

// NewInitCmd new a Command for `kpm init`.
func NewAddCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "add",
		Usage:  "add new dependancy",
		Flags: []cli.Flag{&cli.StringSliceFlag{
			Name:  "git",
			Usage: "Git repository location",
		},
		},
		Action: func(c *cli.Context) error {
			return addGitDep(c)
		},
	}
}

func addGitDep(c *cli.Context) error {
	gitPath := c.StringSlice("git")
	if len(gitPath) > 1 {
		reporter.ExitWithReport("kpm: the argument '--git <URI>' cannot be used multiple times")
	}
	dep := pkg.Dependency{
		Dependency: &pkg.Dependency_Git{
			Git: &pkg.GitDependency{
				Git: gitPath[0],
			},
		},
	}

	return ops.KpmAdd(&dep)
}
