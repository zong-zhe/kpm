// Copyright 2021 The KCL Authors. All rights reserved.

package command

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"kusionstack.io/kpm/core/opt"
	"kusionstack.io/kpm/core/pkg"
	reporter "kusionstack.io/kpm/core/reporter"
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

			kpmHome := os.Getenv("KPM_HOME")
			if kpmHome == "" {
				fmt.Println("kpm: KPM_HOME environment variable is not set")
				fmt.Println("kpm: `add` will be downloaded to directory: ", pwd)
			}

			if err != nil {
				reporter.Fatal("kpm: internal bugs, please contact us to fix it")
			}

			kclPkg, err := pkg.LoadKclPkg(pwd)

			if err != nil {
				reporter.Fatal("kpm: could not load `kcl.mod` in `", pwd, "`")
			}

			gitUrls := c.StringSlice("git")
			if len(gitUrls) > 1 {
				reporter.ExitWithReport("kpm: the argument '--git <URI>' cannot be used multiple times")
			}

			if len(gitUrls) != 0 {
				return addGitDep(&opt.AddOptions{
					LocalPath: kpmHome, // todo: should be KPM_HOME
					RegistryOpts: opt.RegistryOption{
						Git: &opt.GitOption{
							Url: gitUrls[0],
						},
					},
				}, kclPkg)
			}
			return nil
		},
	}
}

func addGitDep(opt *opt.AddOptions, kclPkg *pkg.KclPkg) error {
	return ops.KpmAdd(opt, kclPkg)
}
