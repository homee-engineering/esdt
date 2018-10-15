package main

import (
	"esdt/src"
	"esdt/src/commands"
	"esdt/src/utils"
	"github.com/urfave/cli"
	"log"
	"os"
)

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "c, conn",
		Usage:  "Specify the Elasticsearch cluster the tool points to. Accepts env variable ELASTICSEARCH_URL.\tDefault: " + utils.ElastisearchDefaultUrl,
		Value:  utils.ElastisearchDefaultUrl,
		EnvVar: "ELASTICSEARCH_URL",
	},
	cli.StringFlag{
		Name:  "d, dir",
		Usage: "The target directory for all esdt data.\tDefault: " + utils.DefaultTargetDir,
		Value: utils.DefaultTargetDir,
	},
	cli.StringFlag{
		Name:  "conf, config",
		Usage: "The location of your config YAML.\tDefault: ./es/config.yml",
		Value: "es/config.yml",
	},
	cli.StringFlag{
		Name:  "e, env",
		Usage: "The environment to run the tool in.\tDefault: dev",
		Value: "dev",
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "esdt"
	app.Usage = "Elasticsearch Data Tool. For initializing data on Elasticsearch"
	app.ArgsUsage = "[Command]"
	app.Flags = GlobalFlags
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		commands.RunCommand,
		commands.GenerateCommand,
		commands.RollbackCommand,
	}

	app.CustomAppHelpTemplate = src.AppHelpTemplate
	cli.AppHelpTemplate = src.AppHelpTemplate
	cli.CommandHelpTemplate = src.CommandHelpTemplate
	cli.SubcommandHelpTemplate = src.SubCommandHelpTemplate

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
