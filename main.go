package main

import (
	"esdt/cli/commands"
	"esdt/esdt"
	"github.com/urfave/cli"
	"log"
	"os"
)

var version string

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "c, conn",
		Usage:  "Specify the Elasticsearch cluster the tool points to. Accepts env variable ELASTICSEARCH_URL.\tDefault: " + esdt.DefaultConnUrl,
		EnvVar: "ELASTICSEARCH_URL",
	},
	cli.StringFlag{
		Name:   "d, dir",
		Usage:  "The target directory for all esdt data. Accepts env variable ESDT_TARGET_DIR.\tDefault: " + esdt.DefaultTargetDir,
		EnvVar: "ESDT_TARGET_DIR",
	},
	cli.StringFlag{
		Name:  "conf, config",
		Usage: "The location of your config YAML.\tDefault: ./es/config.yml",
		Value: "es/config.yml",
	},
	cli.StringFlag{
		Name:   "e, env",
		Usage:  "The environment to run the tool in. Accepts env variable ESDT_ENV\tDefault: dev",
		Value:  "dev",
		EnvVar: "ESDT_ENV",
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "esdt"
	app.Usage = "Elasticsearch Data Tool. For initializing data on Elasticsearch"
	app.ArgsUsage = "[Command]"
	app.Flags = GlobalFlags
	app.Version = version
	app.Commands = []cli.Command{
		commands.RunCommand,
		commands.GenerateCommand,
		commands.RollbackCommand,
	}

	app.CustomAppHelpTemplate = commands.AppHelpTemplate
	cli.AppHelpTemplate = commands.AppHelpTemplate
	cli.CommandHelpTemplate = commands.CommandHelpTemplate
	cli.SubcommandHelpTemplate = commands.SubCommandHelpTemplate

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
