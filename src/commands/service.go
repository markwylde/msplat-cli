package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

// CreateServiceCommands : Creates a command for "add"
func CreateServiceCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "services",
			Usage:   "Tasks for managing services",
			Aliases: []string{"sv"},
			Subcommands: []cli.Command{
				{
					Name:  "restart",
					Usage: "Restart a selection of services",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}
}
