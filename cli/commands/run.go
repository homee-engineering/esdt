package commands

import (
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var RunCommand = cli.Command{
	Name:      "run",
	Usage:     "Run a directory of data templates",
	ArgsUsage: "[Flags]",
	Aliases:   []string{"r"},
	Subcommands: []cli.Command{
		HelpCommand,
	},
	Action: runAction,
}

func runAction(c *cli.Context) error {
	e := newEsdt(c)

	err := e.RunAll()
	if err != nil {
		return cli.NewExitError(color.RedString("Failed to create operations index: %s", err.Error()), 1)
	}

	return nil
}
