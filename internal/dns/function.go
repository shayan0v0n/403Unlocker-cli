package dns

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/cavaliergopher/grab/v3"
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

func CheckAndCacheDNS(fileToDownload string) error {
	cacheFile := common.DNS_CONFIG_FILE_CACHED

	dnsList, err := common.ReadDNSFromFile(common.DNS_CONFIG_FILE)
	if err != nil {
		err = common.DownloadConfigFile(common.DNS_CONFIG_URL, common.DNS_CONFIG_FILE)
		if err != nil {
			fmt.Println("Error downloading DNS config file:", err)
			return err
		}

		dnsList, err = common.ReadDNSFromFile(common.DNS_CONFIG_FILE)
		if err != nil {
			fmt.Println("Error reading DNS list from file:", err)
			return err
		}
	}

	fmt.Println("\n+--------------------+------------+")
	fmt.Printf("| %-18s | %-10s |\n", "DNS Server", "Status")
	fmt.Println("+--------------------+------------+")

	var validDNSList []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, dns := range dnsList {
		wg.Add(1)
		go func(dns string) {
			defer wg.Done()

			// Set a 10-second timeout for the DNS lookup
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client := common.ChangeDNS(dns)
			req, err := http.NewRequestWithContext(ctx, "GET", fileToDownload, nil)
			if err != nil {
				mu.Lock()
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Red, "Error", common.Reset)
				mu.Unlock()
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				mu.Lock()
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Red, "Failed", common.Reset)
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			statusCode := resp.StatusCode
			mu.Lock()
			if statusCode == http.StatusOK {
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Green, "OK", common.Reset)
				validDNSList = append(validDNSList, dns)
			} else {
				fmt.Printf("| %-18s | %s%-10s%s |\n", dns, common.Red, "Error", common.Reset)
			}
			mu.Unlock()
		}(dns)
	}

	wg.Wait()

	fmt.Println("+--------------------+------------+")

	fmt.Println("Valid DNS List: ", validDNSList)

	if len(validDNSList) > 0 {
		err = common.WriteDNSToFile(cacheFile, validDNSList)
		if err != nil {
			fmt.Println("Error writing to cached DNS file:", err)
			return err
		}
		fmt.Printf("Cached %d valid DNS servers to %s\n", len(validDNSList), cacheFile)
	} else {
		fmt.Println("No valid DNS servers found to cache.")
	}

	return nil
}

func CheckWithURL(c *cli.Context) error {
	fileToDownload := c.Args().First()

	// Determine the DNS list file based on the presence of the `--check` flag
	var dnsFile string
	if c.Bool("check") {
		CheckAndCacheDNS(fileToDownload) // Update cached DNS if `--check` is used
		dnsFile = common.DNS_CONFIG_FILE_CACHED
	} else {
		dnsFile = common.DNS_CONFIG_FILE
	}

	// Read the DNS list from the determined file
	dnsList, err := common.ReadDNSFromFile(dnsFile)
	if err != nil {
		// Fallback to download and read from the original DNS file
		err = common.DownloadConfigFile(common.DNS_CONFIG_URL, common.DNS_CONFIG_FILE)
		if err != nil {
			return fmt.Errorf("error downloading DNS config file: %w", err)
		}
		dnsList, err = common.ReadDNSFromFile(common.DNS_CONFIG_FILE)
		if err != nil {
			return fmt.Errorf("error reading DNS list from file: %w", err)
		}
	}

	dnsSizeMap := make(map[string]int64)

	timeout := c.Int("timeout")
	fmt.Printf("\nTimeout: %d seconds\n", timeout)
	fmt.Printf("URL: %s\n\n", fileToDownload)

	// Print table header
	fmt.Println("+--------------------+----------------+")
	fmt.Printf("| %-18s | %-14s |\n", "DNS Server", "Download Speed")
	fmt.Println("+--------------------+----------------+")

	tempDir := time.Now().UnixMilli()
	var wg sync.WaitGroup
	for _, dns := range dnsList {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		clientWithCustomDNS := common.ChangeDNS(dns)
		client := grab.NewClient()
		client.HTTPClient = clientWithCustomDNS

		req, err := grab.NewRequest(fmt.Sprintf("/tmp/%v", tempDir), fileToDownload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request for DNS %s: %v\n", dns, err)
		}
		req = req.WithContext(ctx)

		resp := client.Do(req)
		dnsSizeMap[dns] = resp.BytesComplete()

		speed := common.FormatDataSize(resp.BytesComplete() / int64(timeout))
		if resp.BytesComplete() == 0 {
			fmt.Printf("| %-18s | %s%-14s%s |\n", dns, common.Red, speed+"/s", common.Reset)
		} else {
			fmt.Printf("| %-18s | %-14s |\n", dns, speed+"/s")
		}

	}

	wg.Wait()
	// Print table footer
	fmt.Println("+--------------------+----------------+")

	// Find and display the best DNS
	var maxDNS string
	var maxSize int64
	for dns, size := range dnsSizeMap {
		if size > maxSize {
			maxDNS = dns
			maxSize = size
		}
	}

	fmt.Println() // Add a blank line for separation
	if maxDNS != "" {
		bestSpeed := common.FormatDataSize(maxSize / int64(timeout))
		fmt.Printf("Best DNS: %s%s%s (%s%s/s%s)\n",
			common.Green, maxDNS, common.Reset,
			common.Green, bestSpeed, common.Reset)
	} else {
		fmt.Println("No DNS server was able to download any data.")
	}

	os.RemoveAll(fmt.Sprintf("/tmp/%v", tempDir))
	return nil
}
