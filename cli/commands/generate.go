package commands

import (
	"github.com/urfave/cli"
)

var GenerateCommand = cli.Command{
	Name:      "gen",
	Usage:     "Generate an esdt resource",
	ArgsUsage: "[Flags]",
	Subcommands: []cli.Command{
		GenerateOperationCommand,
		GenerateDirCommand,
		HelpCommand,
	},
}
