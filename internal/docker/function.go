package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/salehborhani/403Unlocker-cli/internal/check"
	"github.com/salehborhani/403Unlocker-cli/internal/common"
	"github.com/urfave/cli/v2"
)

// DockerImageValidator validates a Docker image name using a regular expression.
func DockerImageValidator(imageName string) bool {
	pattern := `^(?:[a-zA-Z0-9\-._]+(?::[0-9]+)?/)?` +
		`(?:[a-z0-9\-._]+/)?` +
		`[a-z0-9\-._]+` +
		`(?::[a-zA-Z0-9\-._]+)?` +
		`(?:@[a-zA-Z0-9\-._:]+)?$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(imageName) && !strings.Contains(imageName, "@@")
}

// customTransport tracks the number of bytes transferred during HTTP requests.
type customTransport struct {
	Transport http.RoundTripper
	Bytes     int64
}

// RoundTrip implements the http.RoundTripper interface and wraps the response body to count bytes read.
func (c *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	resp.Body = &countingReader{inner: resp.Body, Bytes: &c.Bytes}
	return resp, nil
}

// countingReader wraps an io.ReadCloser and counts the bytes read.
type countingReader struct {
	inner io.ReadCloser
	Bytes *int64
}

func (cr *countingReader) Read(p []byte) (int, error) {
	n, err := cr.inner.Read(p)
	atomic.AddInt64(cr.Bytes, int64(n))
	return n, err
}

func (cr *countingReader) Close() error {
	return cr.inner.Close()
}

// DownloadDockerImage downloads a Docker image from a registry and tracks the bytes downloaded.
func DownloadDockerImage(ctx context.Context, imageName, registry, outputPath string) (int64, error) {

	fullImageName := registry + "/" + imageName

	// Parse the image reference.
	ref, err := name.ParseReference(fullImageName)
	if err != nil {
		return 0, fmt.Errorf("failed to parse image reference: %v", err)
	}

	auth := authn.DefaultKeychain
	transport := &customTransport{Transport: http.DefaultTransport}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(auth), remote.WithContext(ctx), remote.WithTransport(transport))
	if err != nil {
		return transport.Bytes, fmt.Errorf("failed to download image: %v", err)
	}

	// Ensure output directory exists.
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return transport.Bytes, fmt.Errorf("failed to create output directory: %v", err)
	}

	// Save the image as a tarball.
	tarballPath := filepath.Join(outputPath, filepath.Base(imageName)+".tar")
	if err := tarball.WriteToFile(tarballPath, ref, img); err != nil {
		return transport.Bytes, nil
	}

	return transport.Bytes, nil
}

// CheckWithDockerImage downloads the image from multiple registries and reports the downloaded data size.
func CheckWithDockerImage(c *cli.Context) error {
	registrySizeMap := make(map[string]int64)
	timeout := c.Int("timeout")
	imageName := c.Args().First()
	tempDir := time.Now().UnixMilli()
	fmt.Println("Timeout:", timeout)
	fmt.Println("Docker Image: ", imageName)

	if imageName == "" {
		return fmt.Errorf("image name cannot be empty")
	}

	registryList, err := check.ReadDNSFromFile("config/dockerRegistry.conf")
	if err != nil {
		log.Printf("Error reading registry list: %v", err)
		return err
	}

	for _, registry := range registryList {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		size, err := DownloadDockerImage(ctx, imageName, registry, fmt.Sprintf("/tmp/%v", tempDir))
		if err != nil {
			fmt.Printf("%s: %v\n", registry, "failed")
			continue
		}
		registrySizeMap[registry] += size
		fmt.Printf("%s downloaded : %v\n", registry, common.FormatDataSize(size))
	}
	// Determine which DNS downloaded the most data
	var maxRegistry string
	var maxSize int64
	for dns, size := range registrySizeMap {
		if size > maxSize {
			maxRegistry = dns
			maxSize = size
		}
	}
	if maxRegistry != "" {
		fmt.Printf("best Registry is %s and downloaded the most data: %v\n", maxRegistry, common.FormatDataSize(maxSize))
	} else {
		fmt.Println("No DNS server was able to download any data.")
	}
	os.RemoveAll(fmt.Sprintf("/tmp/%v", tempDir))
	return nil
}
