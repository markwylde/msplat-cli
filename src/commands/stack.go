package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

// CreateStackCommands : Creates a command for "add"
func CreateStackCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "stack",
			Usage: "tasks for managing stacks",
			Subcommands: []cli.Command{
				{
					Name:  "up",
					Usage: "bring up a selection of stacks",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "rm",
					Usage: "remove a selection of stacks",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}
}
