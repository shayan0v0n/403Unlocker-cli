package dns

import (
	"net/url"

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
	return nil
}
