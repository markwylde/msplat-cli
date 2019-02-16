package commands

import (
	"fmt"
	"os"
	"path"
	"strings"

	Auroro "github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"github.com/urfave/cli"

	utils "msplat-cli/src/utils"
)

func ensureDirectoryExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

// CreateProjectCommands : Creates a command for "add"
func CreateProjectCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "projects",
			Usage:   "Tasks for managing projects",
			Aliases: []string{"pr"},
			Subcommands: []cli.Command{
				{
					Name:  "clone",
					Usage: "Clone all or some projects",
					Action: func(c *cli.Context) error {

						fmt.Println("Cloning projects...")

						stacksPath := utils.ResolvePath(viper.GetString("paths.stacks"))

						ensureDirectoryExists(stacksPath)

						var stacks = viper.GetStringMap("stacks")

						for stackKey := range stacks {
							var projects = viper.GetStringMap("stacks." + stackKey)

							fmt.Printf("%s\n", Auroro.Bold(stackKey))
							for projectKey := range projects {
								var project = viper.GetStringMap("stacks." + stackKey + "." + projectKey)
								stackPath := path.Join(stacksPath, stackKey)
								projectPath := path.Join(stacksPath, stackKey, projectKey)

								ensureDirectoryExists(stackPath)

								if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
									gitStatus := utils.ExecuteCwd("git status", projectPath)
									if strings.Contains(gitStatus, "Your branch is up to date") {
										gitStatus = fmt.Sprintf("%s %s", "and is", Auroro.Green("up to date"))
									} else {
										gitStatus = fmt.Sprintf("%s %s", "but is", Auroro.Red("not up to date"))
									}

									fmt.Printf("  %s %s, %s. Skipping clone.\n", projectKey, Auroro.Brown("already exists"), gitStatus)
								} else {
									utils.ExecuteCwd(fmt.Sprintf("git clone %s %s", project["url"], projectKey), stackPath)
									fmt.Printf("  %s %s.\n", projectKey, Auroro.Green("successfully cloned"))
								}
							}
							fmt.Printf("\n")
						}

						return nil
					},
				},
				{
					Name:  "build",
					Usage: "Build a selection of projects",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "Query or list the projects",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}
}
