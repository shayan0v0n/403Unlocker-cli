package check

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"
)

func ChangeDNS(dns string) *http.Client {
	dialer := &net.Dialer{}
	customResolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dnsServer := fmt.Sprintf("%s:53", dns)
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

	url = ensureHTTPS(url)

	var wg sync.WaitGroup
	for _, dns := range dnsList {
		wg.Add(1)
		go func(dns string) {
			defer wg.Done()

			client := ChangeDNS(dns)

			resp, err := client.Get(url)
			if err != nil {
				return
			}

			defer resp.Body.Close()
			code := strings.Split(resp.Status, " ")
			fmt.Printf("DNS: %s %s\n", dns, code[1])
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
func DomainValidator(domain string) bool {
	// Regular expression to validate domain names
	// This regex ensures:
	// - The domain contains only alphanumeric characters, hyphens, and dots.
	// - It does not start or end with a hyphen or dot.
	// - It has at least one dot.
	domainRegex := `^(http[s]?:\/\/)?([a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}).*?$`
	// Match the domain against the regex
	match, _ := regexp.MatchString(domainRegex, domain)
	if !match {
		return false
	}
	// Additional checks:
	// 1. The total length of the domain should not exceed 253 characters.
	if len(domain) > 253 {
		return false
	}

	// 2. Each segment between dots should be between 1 and 63 characters long.
	segments := strings.Split(domain, ".")
	for _, segment := range segments {
		if len(segment) < 1 || len(segment) > 63 {
			return false
		}
	}

	return true
}

func ensureHTTPS(url string) string {
	// Regex to check if the URL starts with https://
	regex := `^(https)://`
	re, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return url
	}
	regexHTTP := `^(http)://`

	reHTTP, err := regexp.Compile(regexHTTP)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return url
	}

	if reHTTP.MatchString(url) {
		url = strings.TrimPrefix(url, "http://")
	}

	// If the URL doesn't start with http:// or https://, prepend https://
	if !re.MatchString(url) {
		url = "https://" + url
	}

	return url
}
