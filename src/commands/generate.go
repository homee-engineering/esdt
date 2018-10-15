package commands

import (
	"esdt/src"
	"github.com/urfave/cli"
)

var GenerateCommand = cli.Command{
	Name:      "gen",
	Usage:     "Generate an esdt resource",
	ArgsUsage: "[Flags]",
	Subcommands: []cli.Command{
		GenerateTemplateCommand,
		GenerateDirCommand,
		src.HelpCommand,
	},
}
