package command

import (
	"os"

	"github.com/urfave/cli/v2"
	ops "kusionstack.io/kpm/ops"
	conf "kusionstack.io/kpm/utils"
)

func NewInitCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "init",
		Usage:  "initialize new module in current directory",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				println("not args...")
				cli.ShowAppHelpAndExit(c, 0)
			}
			println("init...")
			pwd, err := os.Getwd()

			if err != nil {
				println("internal bugs, please contact us to fix it")
			}

			config := conf.NewEmptyConf().SetName(c.Args().First()).SetExecPath(pwd)
			err = ops.KpmInit(config)
			if err == nil {
				println("kpm init finished")
			}
			return err
		},
	}
}
