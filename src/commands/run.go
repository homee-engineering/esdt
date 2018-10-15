package commands

import (
	"esdt/src"
	"esdt/src/utils"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var RunCommand = cli.Command{
	Name:      "run",
	Usage:     "Run a directory of data templates",
	ArgsUsage: "[Flags]",
	Aliases:   []string{"r"},
	Subcommands: []cli.Command{
		src.HelpCommand,
	},
	Action: runAction,
}

func runAction(c *cli.Context) error {
	config, err := utils.GetConfig(c)

	if err != nil {
		return cli.NewExitError(color.RedString(err.Error()), 1)
	}

	ex, err := utils.OperationsIndexExists(config.Conn)
	if err != nil {
		cli.NewExitError(fmt.Sprintf("Failed to check if index exists: %s", err.Error()), 1)
	}

	if !ex {
		err = utils.CreateOperationsIndex(config.Conn)
		if err != nil {
			cli.NewExitError(fmt.Sprintf("Failed to create operations index: %s", err.Error()), 1)
		}
	}

	utils.RunFiles(config, config.TargetDir)

	return nil
}
