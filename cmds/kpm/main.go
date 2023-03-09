package main

import (
	"github.com/urfave/cli/v2"
	"kusionstack.io/kpm/cmds/kpm/command"
)

func main() {
	app := cli.NewApp()
	app.Name = "kpm"
	app.Usage = "kpm is a kcl package manager"
	app.Version = "v0.0.1-alpha.1"
	app.UsageText = "kpm  <command> [arguments]..."
	app.Commands = []*cli.Command{
		command.NewInitCmd(),
		command.NewAddCmd(),
		command.NewDownloadCmd(),
	}
}
