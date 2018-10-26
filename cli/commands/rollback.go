package commands

import (
	"errors"
	"esdt/esdt"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var RollbackCommand = cli.Command{
	Name:      "rollback",
	Usage:     "RollbackFile a single data template or a set of data templates. The data template ID is defined as the filename of the data template minus the extension.",
	ArgsUsage: "[Flags] [Data Template ID]",
	Aliases:   []string{"roll"},
	Subcommands: []cli.Command{
		HelpCommand,
	},
	Action: rollbackAction,
	Flags:  rollbackFlags,
}

var rollbackFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "from, f",
		Usage: "The data template ID to rollback from.\tOptional",
	},
}

func rollbackAction(c *cli.Context) error {
	rollbackId := c.Args().First()
	from := c.String("from")

	if rollbackId == "" {
		return cli.NewExitError(color.RedString("A data template ID is required"), 1)
	}

	if esdt.JsonRegEx.MatchString(rollbackId) {
		rollbackId = strings.TrimSuffix(rollbackId, filepath.Ext(rollbackId))
	}
	if esdt.JsonRegEx.MatchString(from) {
		from = strings.TrimSuffix(from, filepath.Ext(from))
	}

	e := newEsdt(c)

	if from == "" {
		err := e.RollbackFile(rollbackId + ".json")

		handleRollbackError(err, rollbackId)
	} else {
		fi, err := ioutil.ReadDir(e.GetConfig().TargetDir)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not find directory %s", e.GetConfig().TargetDir))
		}
		for _, f := range fi {
			if esdt.JsonRegEx.MatchString(f.Name()) && f.Name() >= from && f.Name() <= rollbackId+".json" {
				err := e.RollbackFile(f.Name())
				handleRollbackError(err, strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())))
			}
		}
	}

	return nil
}

func handleRollbackError(err error, rollbackId string) {
	if err != nil {
		if strings.Contains(err.Error(), esdt.NoRollbackFieldErrorMsg) {
			color.Yellow("RollbackFile %s has no valid rollback field listed", rollbackId)
		} else {
			color.Red("Failed to rollback %s due to %s\n", rollbackId, err.Error())
		}
	} else {
		color.Green("Successfully rolled back %s", rollbackId)
	}
}
