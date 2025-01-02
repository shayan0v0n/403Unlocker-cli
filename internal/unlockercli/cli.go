package

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func Run() {
	app := &cli.App{
		Name: "403unlocker-403cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "check",
				Value: "https://proxy.golang.org/",
				Usage: "Check some urls with provided DNS",
			},
			&cli.StringFlag{
				Name:  "dns",
				Value: "8.8.8.8", // Default to Google DNS
				Usage: "DNS server to use for resolution",
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

	// Use the resolver in a custom dialer
	customDialer := &net.Dialer{
		Resolver: customResolver,
	}

	// Create HTTP transport with custom dialer
	transport := &http.Transport{
		DialContext: customDialer.DialContext,
	}

	// Create HTTP client with custom transport
	client := &http.Client{
		Transport: transport,
	}
	return client
}

func CheckWithDNS(c *cli.Context) error {
	url := c.String("check")
	dns := c.String("dns")

	client := ChangeDNS(dns)

	// Extract the hostname from the URL
	hostname := strings.TrimPrefix(url, "https://")
	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.Split(hostname, "/")[0]

	// Step 1: Verify DNS resolution
	startTime := time.Now()
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return fmt.Errorf("failed to resolve hostname: %v", err)
	}
	resolutionTime := time.Since(startTime)

	log.Printf("Resolved IPs for %s: %v\n", hostname, ips)
	log.Printf("DNS resolution took: %v\n", resolutionTime)

	// Step 2: Verify HTTP request
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Response status for %s: %s\n", url, resp.Status)

	return nil
}