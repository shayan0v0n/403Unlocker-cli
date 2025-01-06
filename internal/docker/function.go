package docker

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/urfave/cli/v2"
)

func DockerImageValidator(imageName string) bool {
	// Regular expression to match a valid Docker image name
	// This pattern allows for optional registry, namespace, and tag/digest
	pattern := `^(?:[a-zA-Z0-9\-._]+(?::[0-9]+)?/)?` + // Optional registry (e.g., docker.io, localhost:5000)
		`(?:[a-z0-9\-._]+/)?` + // Optional namespace (e.g., library, user)
		`[a-z0-9\-._]+` + // Repository name (required)
		`(?::[a-zA-Z0-9\-._]+)?` + // Optional tag (e.g., latest, v1.0)
		`(?:@[a-zA-Z0-9\-._:]+)?$` // Optional digest (e.g., sha256:...)

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Check if the image name matches the pattern
	return regex.MatchString(imageName) && !strings.Contains(imageName, "@@")
}

// downloadDockerImage downloads a Docker image from a registry and saves it as a tarball.
func DownloadDockerImage(imageName, registry, outputPath string, timeoutSeconds int) error {
	imageName = registry + "/" + imageName
	outputPath = outputPath + "/" + imageName
	// Parse the image reference (e.g., "ubuntu:latest")
	ref, err := name.ParseReference(imageName)
	if err != nil {
		return fmt.Errorf("failed to parse image reference: %v", err)
	}
	// Authenticate with the registry (defaults to anonymous auth)
	auth := authn.DefaultKeychain
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(auth), remote.WithContext(ctx))
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Download timed out after %d seconds, saving partially downloaded image...\n", timeoutSeconds)
		} else {
			return fmt.Errorf("failed to fetch image: %v", err)
		}
	}
	// Save the image as a tarball
	err = tarball.WriteToFile(outputPath, ref, img)
	if err != nil {
		return fmt.Errorf("failed to save image as tarball: %v", err)
	}
	fmt.Printf("Image successfully downloaded and saved to %s\n", outputPath)
	return nil
}
func CheckWithDockerImage(c *cli.Context) error {
	return nil
}
