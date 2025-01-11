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
    403unlocker fastdocker gitlab/gitlab-ce:17.0.0-ce.0 --timeout 15`,
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
    403unlocker bestdns https://packages.gitlab.com/gitlab/gitlab-ce/packages/el/7/gitlab-ce-16.8.0-ce.0.el7.x86_64.rpm/download.rpm`,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "timeout",
						Usage:   "Sets timeout",
						Value:   10,
						Aliases: []string{"t"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					if dns.URLValidator(cCtx.Args().First()) {
						return dns.CheckWithURL(cCtx)
					} else {
						err := cli.ShowSubcommandHelp(cCtx)
						if err != nil {
							fmt.Println(err)
						}
					}
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
