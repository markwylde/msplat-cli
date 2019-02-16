package commands

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	Auroro "github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"

	utils "msplat-cli/src/utils"
)

func ensureDirectoryExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func cloneProjects(c *cli.Context) error {
	fmt.Println("Cloning projects...")

	stacks := viper.GetStringMap("stacks")
	stacksPath := utils.ResolvePath(viper.GetString("paths.stacks"))

	ensureDirectoryExists(stacksPath)

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
}

func buildProjects(c *cli.Context) error {
	fmt.Println("Building images...")

	stacks := viper.GetStringMap("stacks")
	stacksPath := utils.ResolvePath(viper.GetString("paths.stacks"))

	for stackKey := range stacks {
		fmt.Printf("%s\n", Auroro.Bold(stackKey))
		projectKey := "configuration"
		environment := "development"
		projectPath := path.Join(stacksPath, stackKey, projectKey, environment)

		resp, err := utils.UnixGet("/v1.24/nodes")
		if err != nil {
			log.Fatal(err)
		}

		machines := gjson.Get(resp, "#.Status.Addr").Array()

		var wg sync.WaitGroup
		wg.Add(len(machines))
		for _, addr := range machines {
			go func(ip string) {
				tag := 0
				fmt.Printf("  Executing %s on machine %s\n", Auroro.Cyan("docker-compose build"), ip)
				_, errText, err := utils.ExecuteCwdStream(fmt.Sprintf("DOCKER_HOST=tcp://%s:2376 docker-compose build", ip), projectPath, func(stdout string) {
					if !c.GlobalBool("verbose") {
						if tag == 1 {
							fmt.Printf("  Successfully built %s on machine %s\n", Auroro.Green(stdout), Auroro.Green(ip))
							tag = 0
						}
						if tag == 0 && stdout == "tagged" {
							tag = 1
						}
					}
				})

				if err != nil {
					log.Fatalf("Build %s on machine %s with error:\n%s\n\n%s", Auroro.Red("failed"), Auroro.Red(ip), errText, Auroro.Magenta("Try running the command again with the --verbose flag for more information"))
				}
				defer wg.Done()
			}(addr.String())
		}
		wg.Wait()

		fmt.Printf("\n")
	}
	fmt.Printf("%s\nEverything is complete.\n", Auroro.Green("Images built successfully"))

	return nil
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
					Name:   "clone",
					Usage:  "Clone all or some projects",
					Action: cloneProjects,
				},
				{
					Name:   "build",
					Usage:  "Build a selection of projects",
					Action: buildProjects,
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
