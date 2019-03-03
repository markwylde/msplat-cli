package commands

import (
	"fmt"
	"log"

	Auroro "github.com/logrusorgru/aurora"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli"
	utils "msplat-cli/src/utils"
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
							containerId := container.Get("Id")

							fmt.Printf("  Restarting %s\n", Auroro.Cyan(containerId))
							cmd := fmt.Sprintf("docker rm -f %s\n", containerId)
							stdout, stderr, err := utils.ExecuteCwd(cmd, "")
							utils.HandleIoError(stdout, stderr, err)
						}

						return nil
					},
				},

				{
					Name:  "logs",
					Usage: "Output the logs for a service",
					Action: func(c *cli.Context) error {
						cmd := fmt.Sprintf("docker service logs -f %s", c.Args().First())
						_, errText, err := utils.ExecuteCwdStream(cmd, "", func(stdout string) {
							fmt.Println(stdout)
						})

						if err != nil {
							log.Fatalf("Error %i: %s", err, errText)
						}

						return nil
					},
				},
			},
		},
	}
}
