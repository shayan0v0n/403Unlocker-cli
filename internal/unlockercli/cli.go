package unlockercli

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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
					if URLValidator(cCtx.Args().First()) {
						return CheckWithDNS(cCtx)
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
					if DockerImageValidator(cCtx.Args().First()) {
						return CheckWithDockerImage(cCtx)
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
					if URLValidator(cCtx.Args().First()) {
						return CheckWithURL(cCtx)
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

func ChangeDNS(dns string) *http.Client {
	dialer := &net.Dialer{}
	customResolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dnsServer := fmt.Sprintf("%s:53", dns)
			log.Printf("Using DNS server: %s\n", dnsServer)
			return dialer.DialContext(ctx, "udp", dnsServer)
		},
	}
	customDialer := &net.Dialer{
		Resolver: customResolver,
	}
	transport := &http.Transport{
		DialContext: customDialer.DialContext,
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}

func CheckWithDNS(c *cli.Context) error {
	url := c.Args().First()
	dnsList, err := ReadDNSFromFile("config/dns.conf")
	if err != nil {
		fmt.Println(err)
		return err
	}

	var wg sync.WaitGroup
	for _, dns := range dnsList {
		wg.Add(1)
		go func(dns string) {
			defer wg.Done()

			client := ChangeDNS(dns)

			hostname := strings.TrimPrefix(url, "https://")
			hostname = strings.TrimPrefix(hostname, "http://")
			hostname = strings.Split(hostname, "/")[0]

			startTime := time.Now()
			ips, err := net.LookupIP(hostname)
			if err != nil {
				log.Printf("Failed to resolve hostname %s with DNS %s: %v\n", hostname, dns, err)
				return
			}
			resolutionTime := time.Since(startTime)

			log.Printf("Resolved IPs for %s: %v (DNS: %s)\n", hostname, ips, dns)
			log.Printf("DNS resolution took: %v\n", resolutionTime)

			resp, err := client.Get(url)
			if err != nil {
				log.Printf("Failed to fetch URL %s with DNS %s: %v\n", url, dns, err)
				return
			}
			defer resp.Body.Close()

			log.Printf("Response status for %s (DNS: %s): %s\n", url, dns, resp.Status)
		}(dns)
	}

	wg.Wait()
	return nil
}

func ReadDNSFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	dnsServers := strings.Fields(string(data))
	return dnsServers, nil
}

// ################### need to be completed ########################
func URLValidator(URL string) bool {
	return false
}

func DockerImageValidator(URL string) bool {
	return false
}

func CheckWithURL(c *cli.Context) error {
	return nil
}

func CheckWithDockerImage(c *cli.Context) error {
	return nil
}
