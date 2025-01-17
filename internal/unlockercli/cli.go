package unlockercli

import (
	"fmt"
	"log"
	"os"

	"github.com/salehborhani/403Unlocker-cli/internal/check"
	"github.com/salehborhani/403Unlocker-cli/internal/dns"
	"github.com/salehborhani/403Unlocker-cli/internal/docker"
	"github.com/urfave/cli/v2"
)

func Run() {
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "403unlocker",
		Usage:                "403Unlocker-CLI is a versatile command-line tool designed to bypass 403 restrictions effectively",
		Commands: []*cli.Command{
			{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "Checks if the DNS SNI-Proxy can bypass 403 error for a specific domain",
				Description: `Examples:
    403unlocker check https://pkg.go.dev`,
				Action: func(cCtx *cli.Context) error {
					if check.DomainValidator(cCtx.Args().First()) {
						return check.CheckWithDNS(cCtx)
					} else {
						err := cli.ShowSubcommandHelp(cCtx)
						if err != nil {
							fmt.Println(err)
						}
					}
					return nil
				},
			},
			{
				Name:    "fastdocker",
				Aliases: []string{"docker"},
				Usage:   "Finds the fastest docker registries for a specific docker image",
				Description: `Examples:
    403unlocker fastdocker --timeout 15 gitlab/gitlab-ce:17.0.0-ce.0`,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "timeout",
						Usage:   "Sets timeout",
						Value:   10,
						Aliases: []string{"t"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					if docker.DockerImageValidator(cCtx.Args().First()) {
						return docker.CheckWithDockerImage(cCtx)
					} else {
						err := cli.ShowSubcommandHelp(cCtx)
						if err != nil {
							fmt.Println(err)
						}
					}
					return nil
				},
			},
			{
				Name:    "bestdns",
				Aliases: []string{"dns"},
				Usage:   "Finds the fastest DNS SNI-Proxy for downloading a specific URL",
				Description: `Examples:
			403unlocker bestdns --timeout 15 https://packages.gitlab.com/gitlab/gitlab-ce/packages/el/7/gitlab-ce-16.8.0-ce.0.el7.x86_64.rpm/download.rpm`,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "timeout",
						Usage:   "Sets timeout in seconds",
						Value:   10,
						Aliases: []string{"t"},
					},
					&cli.BoolFlag{
						Name:    "check",
						Usage:   "Update the DNS cache before running the check",
						Aliases: []string{"c"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					// Validate the URL argument
					if cCtx.Args().Len() < 1 {
						fmt.Println("Error: URL is required")
						return cli.ShowSubcommandHelp(cCtx)
					}

					// Validate the provided URL
					url := cCtx.Args().First()
					if !dns.URLValidator(url) {
						fmt.Println("Error: Invalid URL")
						return cli.ShowSubcommandHelp(cCtx)
					}

					// Call CheckWithURL with the current context
					return dns.CheckWithURL(cCtx)
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
