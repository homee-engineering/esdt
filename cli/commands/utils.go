package commands

import (
	"esdt/esdt"
	"github.com/urfave/cli"
)

func newEsdt(ctx *cli.Context) esdt.Esdt {
	configFile := ctx.GlobalString("config")
	targetDir := ctx.GlobalString("dir")
	conn := ctx.GlobalString("conn")
	env := ctx.GlobalString("env")

	in := &esdt.Config{
		ConfigFile: configFile,
		TargetDir:  targetDir,
		Conn:       conn,
		Env:        env,
	}

	return esdt.New(in)
}
