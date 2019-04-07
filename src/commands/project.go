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

func cloneProjects(c *cli.Context) error {
	fmt.Println("Cloning projects...")

	stacks := viper.GetStringMap("stacks")
	stacksPath := utils.ResolvePath(viper.GetString("paths.stacks"))

	utils.EnsureDirectoryExists(stacksPath)

	for stackKey := range stacks {
		var projects = viper.GetStringMap("stacks." + stackKey)

		fmt.Printf("%s\n", Auroro.Bold(stackKey))
		for projectKey := range projects {
			var project = viper.GetStringMap("stacks." + stackKey + "." + projectKey)
			stackPath := path.Join(stacksPath, stackKey)
			projectPath := path.Join(stacksPath, stackKey, projectKey)

			utils.EnsureDirectoryExists(stackPath)

			if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
				stdout, stderr, err := utils.ExecuteCwd("git status", projectPath)
				utils.HandleIoError(stdout, stderr, err)

				if strings.Contains(stdout, "Your branch is up to date") {
					stdout = fmt.Sprintf("%s %s", "and is", Auroro.Green("up to date"))
				} else {
					stdout = fmt.Sprintf("%s %s", "but is", Auroro.Red("not up to date"))
				}

				fmt.Printf("  %s %s, %s. Skipping clone.\n", projectKey, Auroro.Brown("already exists"), stdout)
			} else {
				stdout, stderr, err := utils.ExecuteCwd(fmt.Sprintf("git clone %s %s", project["url"], projectKey), stackPath)
				utils.HandleIoError(stdout, stderr, err)
				fmt.Printf("  %s %s.\n", projectKey, Auroro.Green("successfully cloned"))
			}
		}
		fmt.Printf("\n")
	}

	return nil
}

func buildOnMachine(wg *sync.WaitGroup, bright bool, ip string, id string, machineName string, projectPath string, verbose bool) {
	fmt.Printf("  Executing %s on machine %s (%s)\n", Auroro.Cyan("docker-compose build"), Auroro.Cyan(machineName), ip)
	cmd := fmt.Sprintf("eval `docker-machine env %s` && docker-compose build", machineName)

	_, errText, err := utils.ExecuteCwdStream(cmd, projectPath, func(stdout string) {
		if !verbose {
			if strings.HasPrefix(stdout, "Successfully tagged ") {
				fmt.Printf("  Successfully built %s on machine %s\n", Auroro.Green(strings.TrimPrefix(stdout, "Successfully tagged ")), Auroro.Green(ip))
			}
		} else {
			if bright {
				fmt.Printf("    %s (%s): %s\n", Auroro.Bold(id), ip, stdout)
			} else {
				fmt.Printf("\u001b[30;1m    %s (%s): %s\u001b[0m\n", id, ip, stdout)
			}
		}
	})

	if err != nil {
		log.Fatalf("Build %s on machine %s with error:\n%s\n\n%s", Auroro.Red("failed"), Auroro.Red(ip), errText, Auroro.Magenta("Try running the command again with the --verbose flag for more information"))
	}
	defer wg.Done()
}

func buildProjects(c *cli.Context) error {
	fmt.Println("Building images...")

	stacks := viper.GetStringMap("stacks")
	stacksPath := utils.ResolvePath(viper.GetString("paths.stacks"))

	for stackKey := range stacks {
		projects := viper.GetStringMap("stacks." + stackKey)

		for projectKey := range projects {
			compose := viper.GetBool("stacks." + stackKey + "." + projectKey + ".compose")
			if compose {
				fmt.Printf("%s -> %s\n", stackKey, projectKey)

				environment := "development"
				projectPath := path.Join(stacksPath, stackKey, projectKey, environment)

				resp, err := utils.UnixGet("/v1.24/nodes")
				if err != nil {
					log.Fatal(err)
				}

				machines := gjson.Parse(resp).Array()

				var wg sync.WaitGroup
				wg.Add(len(machines))
				for i, machine := range machines {
					machineIP := machine.Get("Status.Addr").String()
					machineID := machine.Get("ID").String()
					machineName := machine.Get("Description.Hostname").String()

					if i == 0 {
						go buildOnMachine(&wg, true, machineIP, machineID, machineName, projectPath, c.GlobalBool("verbose"))
					} else {
						go buildOnMachine(&wg, false, machineIP, machineID, machineName, projectPath, c.GlobalBool("verbose"))
					}
				}
				wg.Wait()

				fmt.Printf("\n")
			}
		}
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
