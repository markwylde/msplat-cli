package commands

import (
	"fmt"
	"log"
	utils "msplat-cli/src/utils"
	"path"
	"strings"

	Auroro "github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
)

func ensureSecretsExist(projectPath string, secrets gjson.Result, verbose bool) {
	secrets.ForEach(func(key, value gjson.Result) bool {
		if value.Get("external").Bool() {
			_, stderr, _ := utils.ExecuteCwd(fmt.Sprintf("printf \"a1b2c3d4e5f6g7h8i9j0\" | docker secret create %s -", key), projectPath)

			if strings.Contains(stderr, "code = AlreadyExists") {
				if verbose {
					fmt.Printf("    Secret %s %s\n", key, Auroro.Brown("already exists"))
				}
			} else {
				fmt.Printf("    Secret %s %s\n", key, Auroro.Green("was created"))
			}
		}
		return true // keep iterating
	})
}

func ensureNetworksExist(projectPath string, networks gjson.Result, verbose bool) {
	networks.ForEach(func(key, value gjson.Result) bool {
		if value.Get("external").Bool() {
			name := value.Get("name")
			_, stderr, _ := utils.ExecuteCwd(fmt.Sprintf("docker network create %s --scope swarm --driver overlay --attachable", name), projectPath)

			if strings.Contains(stderr, "already exists") {
				if verbose {
					fmt.Printf("    Network %s (%s) %s\n", name, key, Auroro.Brown("already exists"))
				}
			} else {
				fmt.Printf("    Network %s (%s) %s\n", name, key, Auroro.Green("was created"))
			}
		}
		return true // keep iterating
	})
}

func prepareStack(projectPath string, verbose bool) {
	composeFile := path.Join(projectPath, "docker-compose.yml")
	json, _ := utils.ReadYAMLFileAsJSON(composeFile)

	secrets := gjson.Get(json, "secrets")
	ensureSecretsExist(projectPath, secrets, verbose)

	networks := gjson.Get(json, "networks")
	ensureNetworksExist(projectPath, networks, verbose)
}

func stacksUp(c *cli.Context) error {
	fmt.Println("Starting stacks...")

	stacks := viper.GetStringMap("stacks")
	stacksPath := utils.ResolvePath(viper.GetString("paths.stacks"))

	for stackKey := range stacks {
		projectKey := "configuration"
		environment := "development"
		projectPath := path.Join(stacksPath, stackKey, projectKey, environment)

		fmt.Printf("  Starting %s\n", Auroro.Cyan(stackKey))
		prepareStack(projectPath, c.GlobalBool("verbose"))

		envVars := viper.GetStringMapString(fmt.Sprintf("stacks.%s.configuration.variables", stackKey))

		_, stderr, err := utils.ExecuteCwdStreamWithEnv(fmt.Sprintf("docker stack deploy %s -c docker-compose.yml", stackKey), projectPath, envVars, func(stdout string) {
			if c.GlobalBool("verbose") {
				fmt.Printf("    %s: %s\n", Auroro.Bold(stackKey), stdout)
			}
		})


		if err != nil {
			log.Fatalf("Stacks up error:\n%s", stderr)
		}

		fmt.Printf("\n")
	}
	fmt.Println(Auroro.Green("Stacks brought up successfully"))

	return nil
}

func stacksRm(c *cli.Context) error {
	fmt.Println("Removing stacks...")

	stacks := viper.GetStringMap("stacks")
	stacksPath := utils.ResolvePath(viper.GetString("paths.stacks"))

	for stackKey := range stacks {
		projectKey := "configuration"
		environment := "development"
		projectPath := path.Join(stacksPath, stackKey, projectKey, environment)

		fmt.Printf("  Stopping %s\n", Auroro.Cyan(stackKey))

		_, stderr, err := utils.ExecuteCwdStream(fmt.Sprintf("docker stack rm %s", stackKey), projectPath, func(stdout string) {
			if c.GlobalBool("verbose") {
				fmt.Printf("    %s: %s\n", Auroro.Bold(stackKey), stdout)
			}
		})


		if err != nil {
			log.Fatalf("Stacks rm error:\n%s", stderr)
		}
	}
	fmt.Printf("\n")
	fmt.Println(Auroro.Green("Stacks removed successfully"))

	return nil
}

// CreateStackCommands : Creates a command for "add"
func CreateStackCommands() []cli.Command {
	return []cli.Command{
		{
			Name:    "stacks",
			Usage:   "Tasks for managing stacks",
			Aliases: []string{"st"},
			Subcommands: []cli.Command{
				{
					Name:   "up",
					Usage:  "Bring up a selection of stacks",
					Action: stacksUp,
				},
				{
					Name:  "rm",
					Usage: "Remove a selection of stacks",
					Action: stacksRm,
				},
			},
		},
	}
}
