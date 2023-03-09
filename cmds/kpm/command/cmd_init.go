package command

import (
	"github.com/urfave/cli/v2"
)

func NewInitCmd() *cli.Command {
	return &cli.Command{
		Hidden: false,
		Name:   "init",
		Usage:  "initialize new module in current directory",
		Action: func(c *cli.Context) error {
			println("Create kcl.mod success!")
			return nil
		},
	}
}
