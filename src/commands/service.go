package commands

import (
	"fmt"
	"log"

	utils "msplat-cli/src/utils"

	Auroro "github.com/logrusorgru/aurora"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
)

func outputLogsByService(serviceName string) {
	cmd := fmt.Sprintf("docker service logs -f %s", serviceName)
	utils.ExecuteStream(cmd, "", func(stdout string, stderr string) {
		if stdout != "" {
			fmt.Println(stdout)
		}
		if stderr != "" {
			fmt.Println(stderr)
		}
	})
}

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
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "logs"},
					},
					Action: func(c *cli.Context) error {
						fmt.Println("Searching for containers:", c.Args().First())

						url := fmt.Sprintf(`/v1.24/containers/json?filters={"label":["com.docker.swarm.service.name=%s"]}`, c.Args().First())

						resp, err := utils.UnixGet(url)
						if err != nil {
							log.Fatal(err)
						}

						containers := gjson.Parse(resp).Array()

						if len(containers) == 0 {
							fmt.Println("  No containers found")
						}

						for _, container := range containers {
							containerID := container.Get("Id")

							fmt.Printf("  Restarting %s\n", Auroro.Cyan(containerID))
							cmd := fmt.Sprintf("docker rm -f %s\n", containerID)
							stdout, stderr, err := utils.ExecuteCwd(cmd, "")
							utils.HandleIoError(stdout, stderr, err)
						}

						if c.Bool("logs") {
							outputLogsByService(c.Args().First())
						}

						return nil
					},
				},

				{
					Name:  "logs",
					Usage: "Output the logs for a service",
					Action: func(c *cli.Context) error {
						outputLogsByService(c.Args().First())

						return nil
					},
				},
			},
		},
	}
}
