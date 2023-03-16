// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"github.com/urfave/cli/v2"
	reporter "kusionstack.io/kpm/core/reporter"
	"kusionstack.io/kpm/core/store"
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

	depName:=store.ParseGitNameFromUrl(gitPath[0])
	
}
