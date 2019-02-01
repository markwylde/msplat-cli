package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

// CreateServiceCommands : Creates a command for "add"
func CreateServiceCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "service",
			Usage: "tasks for managing services",
			Subcommands: []cli.Command{
				{
					Name:  "restart",
					Usage: "restart a selection of services",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}
}
