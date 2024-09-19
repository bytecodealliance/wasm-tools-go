//go:build !wasip1 && !wasip2 && !tinygo

package oci

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsOCIPath(t *testing.T) {
	// Create a temporary directory for local file path tests
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		input    string
		local    bool
		expected bool
	}{
		// Valid OCI references
		{
			name:     "ValidOCIReferenceWithTag",
			input:    "nginx:latest",
			local:    false,
			expected: true,
		},
		{
			name:     "ValidOCIReferenceWithRegistry",
			input:    "ghcr.io/webassembly/wasi/http:0.2.0",
			local:    false,
			expected: true,
		},
		{
			name:     "ValidOCIReferenceWithDigest",
			input:    "alpine@sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			local:    false,
			expected: true,
		},
		{
			name:     "ValidOCIReferenceWithLocalhost",
			input:    "localhost:5000/wasi/http",
			local:    false,
			expected: true,
		},
		// Invalid OCI references (should return false)
		{
			name:     "InvalidOCIReferenceSpecialChars",
			input:    "invalid$$$",
			local:    false,
			expected: false,
		},
		{
			name:     "InvalidOCIReferenceTooManyColons",
			input:    "wasi/http:latest:extra",
			local:    false,
			expected: false,
		},
		{
			name:     "InvalidOCIReferenceEmpty",
			input:    "",
			local:    false,
			expected: false,
		},
		{
			name:     "InvalidOCIReferenceDash",
			input:    "-",
			local:    false,
			expected: false,
		},
		// Local file paths (should return false)
		{
			name:     "LocalFilePathAbsolute",
			input:    "usr/local/bin/myfile",
			local:    true,
			expected: false,
		},
		{
			name:     "LocalDirectoryRelative",
			input:    "./mydirectory",
			local:    true,
			expected: false,
		},
		{
			name:     "LocalFilePath",
			input:    "http.wit",
			local:    true,
			expected: false,
		},
		{
			name:     "LocalFilePathLong",
			input:    "testdata/issues/issue163.wit",
			local:    true,
			expected: false,
		},
		{
			name:     "LocalFilePathLong",
			input:    "./testdata/issues/issue163.wit",
			local:    true,
			expected: false,
		},
		// URLs (should return false)
		{
			name:     "URLInput",
			input:    "http://example.com/image",
			local:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input
			if tt.local {
				localPath := filepath.Join(tempDir, tt.input)
				if err := os.MkdirAll(localPath, 0o755); err != nil {
					t.Fatalf("Failed to create directory %q: %v", localPath, err)
				}
				input = localPath
			}
			result := IsOCIPath(input)
			if result != tt.expected {
				t.Errorf("IsOCIPath(%q) = %v; want %v", input, result, tt.expected)
			}
		})
	}
}
