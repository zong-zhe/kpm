// Copyright 2023 The KCL Authors. All rights reserved.

package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"kusionstack.io/kpm/pkg/env"
	"kusionstack.io/kpm/pkg/opt"
	pkg "kusionstack.io/kpm/pkg/package"
	"kusionstack.io/kpm/pkg/reporter"
	"kusionstack.io/kpm/pkg/settings"
)

// NewAddCmd new a Command for `kpm add`.
func NewAddCmd(settings *settings.Settings) *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "add",
		Usage:  "add new dependancy",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "git",
				Usage: "Git repository location",
			},
			&cli.StringSliceFlag{
				Name:  "tag",
				Usage: "Git repository tag",
			},
		},

		Action: func(c *cli.Context) error {

			if c.NArg() == 0 {
				reporter.Report("kpm: module name must be specified.")
				reporter.ExitWithReport("kpm: run 'kpm init help' for more information.")
			}

			pkgName := c.Args().First()

			pwd, err := os.Getwd()

			if err != nil {
				reporter.Fatal("kpm: internal bugs, please contact us to fix it")
			}

			globalPkgPath, err := env.GetAbsPkgPath()
			if err != nil {
				return err
			}

			kclPkg, err := pkg.LoadKclPkg(pwd)
			if err != nil {
				reporter.Fatal("kpm: could not load `kcl.mod` in `", pwd, "`")
			}

			err = kclPkg.ValidateKpmHome(globalPkgPath)
			if err != nil {
				return err
			}

			gitUrl, err := onlyOnceOption(c, "git")

			if err != nil {
				return nil
			}

			gitTag, err := onlyOnceOption(c, "tag")

			if err != nil {
				return err
			}

			var addOpts opt.AddOptions
			if gitUrl != nil {
				addOpts = opt.AddOptions{
					LocalPath: globalPkgPath,
					RegistryOpts: opt.RegistryOptions{
						Git: &opt.GitOptions{
							Url: *gitUrl,
							Tag: *gitTag,
						},
					},
				}
			} else {
				addOpts = opt.AddOptions{
					LocalPath: globalPkgPath,
					RegistryOpts: opt.RegistryOptions{
						Oci: &opt.OciOptions{
							Reg:     "ghcr.io",
							Repo:    "zong-zhe",
							PkgName: pkgName,
							Tag:     "v0.0.1",
						},
					},
				}
			}

			// err = addOpts.Validate()
			// if err != nil {
			// 	return err
			// }

			err = addGitDep(&addOpts, kclPkg, settings)
			if err != nil {
				return err
			}
			reporter.Report("kpm: add dependency '", *gitUrl, "'", "with tag '", *gitTag, "' successfully.")
			return nil
		},
	}
}

// onlyOnceOption is used to check that the value of some parameters can only appear once.
func onlyOnceOption(c *cli.Context, name string) (*string, error) {
	inputOpt := c.StringSlice(name)
	if len(inputOpt) > 1 {
		reporter.ExitWithReport("kpm: the argument '", name, "' cannot be used multiple times")
		reporter.ExitWithReport("kpm: run 'kpm add help' for more information.")
		return nil, fmt.Errorf("kpm: Invalid command")
	} else if len(inputOpt) == 1 {
		return &inputOpt[0], nil
	} else {
		return nil, nil
	}
}

func addGitDep(opt *opt.AddOptions, kclPkg *pkg.KclPkg, settings *settings.Settings) error {
	// if opt.RegistryOpts.Git == nil {
	// 	reporter.Report("kpm: a value is required for '-git <URI>' but none was supplied")
	// 	reporter.ExitWithReport("kpm: run 'kpm add help' for more information.")
	// }

	return kclPkg.AddDeps(opt, settings)
}
