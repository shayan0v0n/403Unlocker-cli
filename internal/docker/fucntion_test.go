package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDockerImageName(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		expected bool
	}{
		// Valid image names
		{"Valid image name without registry or tag", "ubuntu", true},
		{"Valid image name with namespace", "library/ubuntu", true},
		{"Valid image name with registry", "docker.io/library/ubuntu", true},
		{"Valid image name with custom registry and port", "localhost:5000/myproject/ubuntu", true},
		{"Valid image name with tag", "myregistry/myproject/ubuntu:latest", true},
		{"Valid image name with digest", "myregistry/myproject/ubuntu@sha256:abc123", true},

		// Invalid image names
		{"Invalid image name with uppercase letters in repository", "MyRegistry/MyProject/Ubuntu", false},
		{"Invalid image name with invalid characters", "invalid!@#/image/name", false},
		{"Invalid image name with empty repository", "", false},
		{"Invalid image name with only slashes", "///", false},
		{"Invalid image name with multiple @ symbols", "invalid@image@name", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DockerImageValidator(tt.image)
			assert.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
