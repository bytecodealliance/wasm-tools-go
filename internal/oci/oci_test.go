//go:build !wasip1 && !wasip2

package oci

import "testing"

func TestIsOCIPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid OCI references
		{
			name:     "ValidOCIReferenceWithTag",
			input:    "nginx:latest",
			expected: true,
		},
		{
			name:     "ValidOCIReferenceWithRegistry",
			input:    "ghcr.io/webassembly/wasi/http:0.2.0",
			expected: true,
		},
		{
			name:     "ValidOCIReferenceWithDigest",
			input:    "alpine@sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			expected: true,
		},
		{
			name:     "ValidOCIReferenceWithLocalhost",
			input:    "localhost:5000/wasi/http",
			expected: true,
		},
		// Invalid OCI references (should return false)
		{
			name:     "InvalidOCIReferenceSpecialChars",
			input:    "invalid$$$",
			expected: false,
		},
		{
			name:     "InvalidOCIReferenceTooManyColons",
			input:    "wasi/http:latest:extra",
			expected: false,
		},
		{
			name:     "InvalidOCIReferenceEmpty",
			input:    "",
			expected: false,
		},
		// Local file paths (should return false)
		{
			name:     "LocalFilePathAbsolute",
			input:    "/usr/local/bin/myfile",
			expected: false,
		},
		{
			name:     "LocalDirectoryRelative",
			input:    "./mydirectory",
			expected: false,
		},
		{
			name:     "LocalFilePath",
			input:    "./http.wit",
			expected: false,
		},
		// URLs (should return false)
		{
			name:     "URLInput",
			input:    "http://example.com/image",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := IsOCIPath(tt.input)
			if result != tt.expected {
				t.Errorf("IsOCIPath(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
