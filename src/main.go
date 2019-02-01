package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"

	commands "msplat-cli/src/commands"
)

func main() {
	app := cli.NewApp()

	app.Name = "msplat toolkit"
	app.Description = "A cli for managing msplat environments"
	app.Usage = "A cli for managing msplat environments"
	app.Version = "0.0.1"

	app.EnableBashCompletion = true

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t\t\t\t"}} {{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "./config.yml",
			Usage:  "configuration settings for msplat",
			EnvVar: "MSPLAT_CONFIG",
		},
	}

	app.Commands = append(app.Commands, commands.CreateProjectCommands()...)
	app.Commands = append(app.Commands, commands.CreateServiceCommands()...)
	app.Commands = append(app.Commands, commands.CreateStackCommands()...)

	sort.Sort(cli.CommandsByName(app.Commands))

	app.Before = func(c *cli.Context) error {
		ReadConfig(c.String("config"))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
