package commands

import (
	"github.com/urfave/cli"
	"strings"
)

const AppHelpTemplate = `
{{.Name}} - {{.Usage}}
{{if .Description}}
Description:
  {{.Description}}
{{end}}
Usage:
  {{.HelpName}} {{.ArgsUsage}}
{{if .Commands}}
Available Commands:	
  {{range .Commands}}{{join .Names ", "}}{{"\t"}}{{.Usage}}{{ "\n  " }}{{end}}
{{end -}}
{{if .Flags -}}
Available Flags:
  {{range .Flags}} -{{.Name}}{{"\t"}}{{"\t"}}{{"\t"}}{{.Usage}}{{"\n  "}}{{end}}
{{end -}}
`

const CommandHelpTemplate = `
{{.HelpName}} - {{.Usage}}
{{if .Description}}
Description:
  {{.Description}}
{{end}}
Usage:
  {{.HelpName}} {{.ArgsUsage}}
{{if .Subcommands}}
Available Commands:	
  {{range .Subcommands}}{{join .Names ", "}}{{"\t"}}{{.Usage}}{{"\n  "}}{{end}}
{{end -}}
{{if .Flags -}}
Available Flags:
  {{range .Flags}} -{{.Name}}{{"\t"}}{{"\t"}}{{"\t"}}{{.Usage}}{{"\n  "}}{{end}}
{{end -}}
`

const SubCommandHelpTemplate = `
{{.HelpName}} - {{.Usage}}
{{if .Description}}
Description:
  {{.Description}}
{{end}}
Usage:
  {{.HelpName}} {{.ArgsUsage}}
{{if .Subcommands}}
Available Commands:	
  {{range .Subcommands}}-{{.Name}}{{"\t"}}{{.Usage}}{{ "\n  " }}{{end}}
{{end -}}
{{if .Flags -}}
Available Flags:
    {{range .Flags}} -{{.Name}}{{"\t"}}{{"\t"}}{{"\t"}}{{.Usage}}{{"\t"}}{{if .Value}}Default: {{.Value}}{{else}}Required{{end}}{{"\n  "}}{{end}}
{{end -}}
`

var HelpCommand = cli.Command{
	Name:      "help",
	Aliases:   []string{"h"},
	Usage:     "Shows a list of commands or help for one command",
	ArgsUsage: "[command]",
	Action: func(c *cli.Context) error {
		args := strings.Fields(c.Command.HelpName)
		if len(args) >= 2 {
			commandString := args[len(args)-2]
			ctx := c.Parent().Parent()
			command := ctx.App.Command(commandString)
			if command == nil {
				cli.ShowAppHelp(c)
				return nil
			} else {
				return cli.ShowCommandHelp(ctx, commandString)
			}
		} else {
			cli.ShowAppHelp(c)
			return nil
		}
	},
}

var HelpSubcommand = cli.Command{
	Name:      "help",
	Aliases:   []string{"h"},
	Usage:     "Shows a list of commands or help for one command",
	ArgsUsage: "[command]",
	Action: func(c *cli.Context) error {
		cmd := strings.Fields(c.Command.HelpName)
		commandString := cmd[len(cmd)-2]

		if !c.Command.HasName(commandString) {
			commandString = ""
		}

		if commandString == "" {
			cli.ShowCommandHelp(c, commandString)
			return nil
		} else {
			return cli.ShowSubcommandHelp(c)
		}
	},
}
