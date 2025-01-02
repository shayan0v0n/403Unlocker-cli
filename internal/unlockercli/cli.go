package unlockercli

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func Run() {
	app := &cli.App{
		Name: "403unlocker-cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "check",
				Value: "https://pkg.go.dev/",
				Usage: "Check some urls with provided DNS",
			},
		},
		Action: CheckWithDNS,
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
	url := c.String("check")

	dnsList, err := ReadDNSFromFile("config/dns.conf")

	if err != nil {
		fmt.Println(err)
	}

	for _, dns := range dnsList {

		client := ChangeDNS(dns)

		hostname := strings.TrimPrefix(url, "https://")
		hostname = strings.TrimPrefix(hostname, "http://")
		hostname = strings.Split(hostname, "/")[0]

		startTime := time.Now()
		ips, err := net.LookupIP(hostname)
		if err != nil {
			continue
			//return fmt.Errorf("failed to resolve hostname: %v", err)
		}
		resolutionTime := time.Since(startTime)

		log.Printf("Resolved IPs for %s: %v\n", hostname, ips)
		log.Printf("DNS resolution took: %v\n", resolutionTime)

		resp, err := client.Get(url)
		if err != nil {
			continue
			//return fmt.Errorf("failed to fetch URL: %v", err)
		}
		defer resp.Body.Close()

		log.Printf("Response status for %s: %s\n", url, resp.Status)
	}
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
