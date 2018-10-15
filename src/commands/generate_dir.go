package commands

import (
	"esdt/src"
	"esdt/src/io"
	"github.com/urfave/cli"
	"os"
	"path"
)

var GenerateDirCommand = cli.Command{
	Name:      "dir",
	Usage:     "Generate the directory used by esdt",
	ArgsUsage: "[Flags]",
	Subcommands: []cli.Command{
		src.HelpCommand,
	},
	Action: generateDirAction,
	Flags:  generateDirFlags,
}

var generateDirFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "target, t",
		Usage: "The target directory to generate the directory within.\tDefault: ./",
		Value: ".",
	},
}

func generateDirAction(c *cli.Context) error {
	target := c.String("target")

	os.MkdirAll(path.Join(target, "es/operations"), os.ModePerm)

	fp, _ := io.GetFile("template-config.yml")
	os.Rename(fp, path.Join(target, "es", "config.yml"))

	return nil
}
