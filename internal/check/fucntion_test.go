package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainValidator(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{
			name:     "Valid URL",
			domain:   "https://pkg.go.dev",
			expected: true,
		},
		{
			name:     "Valid domain",
			domain:   "http://example.com",
			expected: true,
		},
		{
			name:     "Valid subdomain",
			domain:   "https://sub.example.com",
			expected: true,
		},
		{
			name:     "Invalid domain - hyphen at start",
			domain:   "-invalid.com",
			expected: false,
		},
		{
			name:     "Invalid domain - hyphen at end",
			domain:   "invalid-.com",
			expected: false,
		},
		{
			name:     "Invalid domain - missing top-level domain",
			domain:   "example",
			expected: false,
		},
		{
			name:     "Invalid domain - double dots",
			domain:   "invalid..com",
			expected: false,
		},
		{
			name:     "Valid domain with hyphens",
			domain:   "https://valid-domain.org",
			expected: true,
		},
		{
			name:     "Invalid domain - too long",
			domain:   "toolongdomainnamethatiswaylongerthanthemaximumallowedlengthof253charactersandshouldfailvalidationbecauseitistoolongandexceedsthelimit.toolongdomainnamethatiswaylongerthanthemaximumallowedlengthof253charactersandshouldfailvalidationbecauseitistoolongandexceedsthelimit.toolongdomainnamethatiswaylongerthanthemaximumallowedlengthof253charactersandshouldfailvalidationbecauseitistoolongandexceedsthelimit.toolongdomainnamethatiswaylongerthanthemaximumallowedlengthof253charactersandshouldfailvalidationbecauseitistoolongandexceedsthelimit",
			expected: false,
		},
		{
			name:     "Invalid domain - starts with dot",
			domain:   ".invalid",
			expected: false,
		},
		{
			name:     "Invalid domain - ends with dot",
			domain:   "invalid.",
			expected: false,
		},
		{
			name:     "Invalid domain without scheme",
			domain:   "pkg.go.dev",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DomainValidator(tt.domain)
			assert.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestEnsureHTTPS(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "URL with http scheme",
			url:      "http://example.com",
			expected: "https://example.com/",
		},
		{
			name:     "URL with https scheme",
			url:      "https://example.com",
			expected: "https://example.com/",
		},
		{
			name:     "URL without scheme",
			url:      "example.com",
			expected: "https://example.com/",
		},
		{
			name:     "URL with path and query",
			url:      "http://example.com/path?query=123",
			expected: "https://example.com/",
		},
		{
			name:     "URL with subdomain",
			url:      "http://sub.example.com",
			expected: "https://sub.example.com/",
		},
		{
			name:     "URL with double slashes",
			url:      "http://example.com//path",
			expected: "https://example.com/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ensureHTTPS(tt.url)
			assert.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
