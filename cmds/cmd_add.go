// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"os"

	"github.com/urfave/cli/v2"
	"kusionstack.io/kpm/core/conf"
	"kusionstack.io/kpm/core/git"
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
			pwd, err := os.Getwd()

			if err != nil {
				reporter.Fatal("kpm: internal bugs, please contact us to fix it")
			}

			gitPath := c.StringSlice("git")
			if len(gitPath) > 1 {
				reporter.ExitWithReport("kpm: the argument '--git <URI>' cannot be used multiple times")
			}
			if len(gitPath) != 0 {
				conf := conf.NewEmptyConf().SetKclModPath(pwd)
				return addGitDep(&conf, git.NewGitOption().SetUrl(gitPath[0]))
			}
			return nil
		},
	}
}

func addGitDep(conf *conf.Config, gitOpt *git.GitOption) error {

	dep := pkg.Dependency{
		Name: git.ParseRepoNameFromGitUrl(gitOpt.Url()),
		Dependency: &pkg.Dependency_Git{
			Git: &pkg.GitDependency{
				Git: gitOpt.Url(),
			},
		},
	}

	return ops.KpmAdd(conf, &dep)
}
