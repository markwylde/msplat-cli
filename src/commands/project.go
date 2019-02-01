package commands

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// CreateProjectCommands : Creates a command for "add"
func CreateProjectCommands() []cli.Command {
	return []cli.Command{
		{
			Name:  "project",
			Usage: "tasks for managing projects",
			Subcommands: []cli.Command{
				{
					Name:  "clone",
					Usage: "clone all or some projects",
					Action: func(c *cli.Context) error {
						fmt.Println("clone is not implemented yet\nthe following is a list of the stacks projects:\n")

						var stacks = viper.GetStringMap("stacks")

						for stackKey := range stacks {
							var projects = viper.GetStringMap("stacks." + stackKey)

							fmt.Printf("%s\n", stackKey)
							for projectKey := range projects {
								var project = viper.GetStringMap("stacks." + stackKey + "." + projectKey)
								fmt.Printf(" %s=%s\n", projectKey, project["url"])
							}
							fmt.Printf("\n")
						}

						return nil
					},
				},
				{
					Name:  "build",
					Usage: "build a selection of projects",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "query or list the projects",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}
}
