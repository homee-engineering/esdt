package commands

import (
	"errors"
	"esdt/src"
	"esdt/src/utils"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var RollbackCommand = cli.Command{
	Name:      "rollback",
	Usage:     "Rollback a single data template or a set of data templates. The rollback ID is defined as the filename of the data template minus the extension.",
	ArgsUsage: "[Flags] [Rollback ID]",
	Aliases:   []string{"roll"},
	Subcommands: []cli.Command{
		src.HelpCommand,
	},
	Action: rollbackAction,
	Flags:  rollbackFlags,
}

var rollbackFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "from, f",
		Usage: "The rollback ID to rollback from.\tOptional",
	},
}

func rollbackAction(c *cli.Context) error {
	rollbackId := c.Args().First()
	from := c.String("from")
	dir := c.GlobalString("dir")
	conn := c.GlobalString("conn")

	if rollbackId == "" {
		return cli.NewExitError(color.RedString("A rollback ID is required"), 1)
	}

	if from == "" {
		rollback(conn, dir, rollbackId+".json")
	} else {
		fi, err := ioutil.ReadDir(dir)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not find directory %s", dir))
		}
		for _, f := range fi {
			if f.Name() >= from && f.Name() <= rollbackId+".json" {
				rollback(conn, dir, f.Name())
			}
		}
	}

	return nil
}

func rollback(conn string, dir string, dataTemplateFilename string) {
	dt, err := utils.LoadDataTemplate(dir, dataTemplateFilename)
	rollbackId := strings.TrimSuffix(dataTemplateFilename, filepath.Ext(dataTemplateFilename))
	if err != nil {
		color.Red("Failed to rollback %s due to %s\n", rollbackId, err.Error())
	}
	err = utils.RollbackDataTemplate(conn, dt)
	if err != nil {
		if strings.Contains(err.Error(), utils.NoRollbackFieldErrorMsg) {
			color.Yellow("Rollback %s has no valid rollback field listed", rollbackId)
		} else {
			color.Red("Failed to rollback %s due to %s\n", rollbackId, err.Error())
		}
	} else {

		color.Green("Successfully rolled back %s", rollbackId)

		utils.DeleteOperationIndex(conn, rollbackId)
	}
}
