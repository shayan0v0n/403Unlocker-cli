package dns

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/salehborhani/403Unlocker-cli/internal/check"
	"github.com/salehborhani/403Unlocker-cli/internal/common"
	"github.com/urfave/cli/v2"
)

func URLValidator(URL string) bool {
	// Parse the URL
	u, err := url.Parse(URL)
	if err != nil {
		return false
	}
	// Check if the scheme is either "http" or "https"
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	// Check if the host is present
	if u.Host == "" {
		return false
	}
	return true
}
func CheckWithURL(c *cli.Context) error {
	fileToDownload := c.Args().First()
	timeout := c.Int("timeout")
	dnsList, err := check.ReadDNSFromFile("config/dns.conf")
	if err != nil {
		fmt.Println("Error reading DNS list:", err)
		return err
	}
	// Map to store the total size downloaded by each DNS
	dnsSizeMap := make(map[string]int64)
	fmt.Println("Timeout:", timeout)
	fmt.Println("URL: ", fileToDownload)
	tempDir := time.Now().UnixMilli()
	for _, dns := range dnsList {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		// Create a custom HTTP client with the specified DNS
		clientWithCustomDNS := check.ChangeDNS(dns)
		client := grab.NewClient()
		client.HTTPClient = clientWithCustomDNS
		// Create a new download request
		req, err := grab.NewRequest(fmt.Sprintf("/tmp/%v", tempDir), fileToDownload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request for DNS %s: %v\n", dns, err)
		}
		req = req.WithContext(ctx)
		// Start the download
		resp := client.Do(req)
		// Update the total size downloaded by this DNS
		dnsSizeMap[dns] += resp.BytesComplete() // Use BytesComplete() for partial downloads
		fmt.Printf("%v\tDNS: %s\n", common.FormatDataSize(resp.BytesComplete()), dns)
	}
	// Determine which DNS downloaded the most data
	var maxDNS string
	var maxSize int64
	for dns, size := range dnsSizeMap {
		if size > maxSize {
			maxDNS = dns
			maxSize = size
		}
	}
	if maxDNS != "" {
		fmt.Printf("best DNS is %s and downloaded the most data: %v\n", maxDNS, common.FormatDataSize(maxSize))
	} else {
		fmt.Println("No DNS server was able to download any data.")
	}
	os.RemoveAll(fmt.Sprintf("/tmp/%v", tempDir))
	return nil
}
