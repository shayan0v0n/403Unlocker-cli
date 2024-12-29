package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "403unlocker-cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "check",
				Value: "https://proxy.golang.org/",
				Usage: "Check some urls with provided DNS",
			},
		},
		Action: CheckWithDNS,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func CheckWithDNS(cCtx *cli.Context) error {
	url := cCtx.String("check")
	// Check with these DNS
	dns := GetDNS("dns.txt")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, d := range dns {
		client := ChangeDNS(d)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

		if err != nil {
			fmt.Println("err: ", err)
		}

		resp, err := client.Do(req)

		if err != nil {
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				fmt.Println("The problematic DNS: ", d)
				fmt.Println("err: ", err)
				continue
			} else {
				fmt.Println("err: ", err)
			}
		}

		resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println("Its not ok with dns: ", d)
		} else {
			fmt.Println("Its ok")
		}
	}

	return nil
}

func ChangeDNS(dns string) *http.Client {
	dialer := &net.Dialer{}

	// Use a custom resolver with specific DNS server(s)
	customResolver := &net.Resolver{
		PreferGo: true, // Force the Go resolver
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			// Replace with your DNS server IP and port
			dnsServer := fmt.Sprintf("%s:53", dns) // Example: Google DNS
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

func GetDNS(path string) []string {
	var dns []string

	file, err := os.Open(path) // Replace with your file name if different
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	// Create a scanner to read the file line-by-line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Read the line and split by spaces
		line := scanner.Text()
		ips := strings.Fields(line)
		// Iterate through each IP address
		for _, d := range ips {
			dns = append(dns, d)
		}
	}

	// Check for errors while reading the file
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return dns
}
