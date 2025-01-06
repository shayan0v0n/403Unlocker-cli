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
				Usage:   "Checks if the DNS SNI-Proxy can bypass 403 error for an specific domain",
				Action: func(cCtx *cli.Context) error {
					if check.DomainValidator(cCtx.Args().First()) {
						return check.CheckWithDNS(cCtx)
					} else {
						fmt.Println("need a valid domain		example: https://pkg.go.dev")
					}
					return nil
				},
			},
			{
				Name:    "docker",
				Aliases: []string{"d"},
				Usage:   "Finds the fastest docker registries for an specific docker image",
				Action: func(cCtx *cli.Context) error {
					if docker.DockerImageValidator(cCtx.Args().First()) {
						return docker.CheckWithDockerImage(cCtx)
					} else {
						fmt.Println("need a valid docker image		example: gitlab/gitlab-ce:17.0.0-ce.0")
					}
					return nil
				},
			},
			{
				Name:  "dns",
				Usage: "Finds the fastest DNS SNI-Proxy for downloading an specific URL",
				Action: func(cCtx *cli.Context) error {
					if dns.URLValidator(cCtx.Args().First()) {
						return dns.CheckWithURL(cCtx)
					} else {
						fmt.Println("need a valid URL		example: \"https://packages.gitlab.com/gitlab/gitlab-ce/packages/el/7/gitlab-ce-16.8.0-ce.0.el7.x86_64.rpm/download.rpm\"")
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
