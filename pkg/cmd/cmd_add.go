// Copyright 2023 The KCL Authors. All rights reserved.

package cmd

import (
	"fmt"
	"os"
	"strings"

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

			addOpts, err := parseAddOptions(c, globalPkgPath, settings)
			if err != nil {
				return err
			}

			err = addOpts.Validate()
			if err != nil {
				return err
			}

			err = kclPkg.AddDeps(addOpts, settings)
			if err != nil {
				return err
			}
			reporter.Report("kpm: add dependency successfully.")
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

func parseAddOptions(c *cli.Context, localPath string, settings *settings.Settings) (*opt.AddOptions, error) {
	if c.NArg() == 0 {
		gitOpts, err := parseGitRegistryOptions(c)
		if err != nil {
			return nil, err
		}
		return &opt.AddOptions{
			LocalPath:    localPath,
			RegistryOpts: *gitOpts,
		}, nil
	} else {
		ociReg, err := parseOciRegistryOptions(c, settings)
		if err != nil {
			return nil, err
		}
		return &opt.AddOptions{
			LocalPath:    localPath,
			RegistryOpts: *ociReg,
		}, nil
	}

}

func parseGitRegistryOptions(c *cli.Context) (*opt.RegistryOptions, error) {
	gitUrl, err := onlyOnceOption(c, "git")

	if err != nil {
		return nil, nil
	}

	gitTag, err := onlyOnceOption(c, "tag")

	if err != nil {
		return nil, err
	}

	return &opt.RegistryOptions{
		Git: &opt.GitOptions{
			Url: *gitUrl,
			Tag: *gitTag,
		},
	}, nil
}

func parseOciRegistryOptions(c *cli.Context, settings *settings.Settings) (*opt.RegistryOptions, error) {
	ociPkgRef := c.Args().First()
	name, version := parseNameAndVersion(ociPkgRef)
	if len(version) == 0 {
		reporter.Report("kpm: default version 'latest' of the package will be downloaded.")
		version = opt.DEFAULT_OCI_TAG
	}

	return &opt.RegistryOptions{
		Oci: &opt.OciOptions{
			Reg:     settings.DefauleOciRegistry(),
			Repo:    settings.DefauleOciRepo(),
			PkgName: name,
			Tag:     version,
		},
	}, nil
}

func parseNameAndVersion(s string) (string, string) {
	parts := strings.Split(s, "@")
	if len(parts) == 1 {
		return parts[0], ""
	}

	if len(parts) > 2 {
		return "", ""
	}

	return parts[0], parts[1]
}
