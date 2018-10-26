package commands

import (
	"esdt/cli/io"
	"github.com/fatih/color"
	"github.com/homee-engineering/go-commons/slice"
	"github.com/urfave/cli"
	"os"
	"path"
	"strings"
	"time"
)

const timeFormatString = "20060102150405"

var GenerateOperationCommand = cli.Command{
	Name:      "operation",
	Usage:     "Generate an Elasticsearch data operation",
	ArgsUsage: "[Flags] [Name]",
	Aliases:   []string{"op", "o"},
	Subcommands: []cli.Command{
		HelpCommand,
	},
	Action: generateOperationAction,
	Flags:  generateOperationFlags,
}

var generateOperationFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "method, m",
		Usage: "The HTTP Method to be used for Elasticsearch. Can be GET, PUT, POST, HEAD, DELETE\tDefault: GET",
		Value: "get",
	},
	cli.StringFlag{
		Name:  "uri, u",
		Usage: "The URI to be used against Elasticsearch e.g. _bulk. Will be concatenated with the conn arg\tRequired",
	},
}

var validElasticSearchHttpMethods = []string{
	"GET",
	"PUT",
	"POST",
	"HEAD",
	"DELETE",
}

type templateModel struct {
	Method         string
	Uri            string
	OppositeMethod string
}

func generateOperationAction(c *cli.Context) error {
	name := c.Args().First()
	method := c.String("method")
	uri := c.String("uri")
	e := newEsdt(c)

	if method != "" && !slice.ContainsStringCaseInsensitive(validElasticSearchHttpMethods, method) {
		return cli.NewExitError(color.RedString("HTTP method must be one of %s", strings.Join(validElasticSearchHttpMethods, ", ")), 1)
	}

	if name == "" {
		return cli.NewExitError(color.RedString("A data template name is required"), 1)
	}

	timestamp := time.Now().Format(timeFormatString)
	fileName := timestamp + "_" + name + ".json"
	oppositeMethod := "delete"
	switch strings.ToLower(method) {
	case "delete":
		oppositeMethod = "post"
	case "post", "put":
		oppositeMethod = "delete"
	}
	fp, err := io.ApplyTemplate("template.json", templateModel{
		Method:         strings.ToUpper(method),
		Uri:            uri,
		OppositeMethod: strings.ToUpper(oppositeMethod),
	})

	if err != nil {
		return cli.NewExitError(color.RedString("Failed to generate json: %s", err), 1)
	}

	err = os.Rename(fp, path.Join(e.GetConfig().TargetDir, fileName))
	if err != nil {
		return cli.NewExitError(color.RedString("Failed to generate json: %s", err), 1)
	}

	return nil
}
