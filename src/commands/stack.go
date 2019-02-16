package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

// CreateStackCommands : Creates a command for "add"
func CreateStackCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "stacks",
			Usage:   "Tasks for managing stacks",
			Aliases: []string{"st"},
			Subcommands: []cli.Command{
				{
					Name:  "up",
					Usage: "Bring up a selection of stacks",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "rm",
					Usage: "Remove a selection of stacks",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}
}
