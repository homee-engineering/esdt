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
	pw := ctx.GlobalString("password")
	user := ctx.GlobalString("username")

	in := &esdt.Config{
		ConfigFile: configFile,
		TargetDir:  targetDir,
		Conn:       conn,
		Env:        env,
		Password:   pw,
		Username:   user,
	}

	return esdt.New(in)
}
