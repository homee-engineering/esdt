package utils

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io/ioutil"
)

type config struct {
	Conn      string `yaml:"conn"`
	TargetDir string `yaml:"dir"`
}

type yamlConfig map[string]*config

func GetConfig(ctx *cli.Context) (c *config, err error) {
	configFile := ctx.GlobalString("config")
	targetDir := ctx.GlobalString("dir")
	conn := ctx.GlobalString("conn")
	env := ctx.GlobalString("env")
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Problems reading file %s: ", configFile))
		return
	}
	var yc yamlConfig
	err = yaml.Unmarshal(content, &yc)
	if err != nil {
		return
	}

	c = yc[env]

	if conn != ElastisearchDefaultUrl {
		c.Conn = conn
	}

	if targetDir != DefaultTargetDir {
		c.TargetDir = targetDir
	}

	if c == nil {
		err = errors.New(fmt.Sprintf("Could not find environment %s in file %s", env, configFile))
	}

	return
}
