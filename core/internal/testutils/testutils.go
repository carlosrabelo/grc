package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

// CreateTempYAMLFile creates a temporary YAML file with the provided content
func CreateTempYAMLFile(t *testing.T, content string) string {
	return CreateTempFile(t, "test.yaml", content)
}

// CreateTempFile creates a temporary file with the provided name and content
func CreateTempFile(t *testing.T, name, content string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, name)
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tmpFile
}

// CleanupFile removes a file safely, ignoring errors
func CleanupFile(path string) {
	_ = os.Remove(path)
}

// BoolPtr returns a pointer to the provided boolean value
func BoolPtr(v bool) *bool {
	return &v
}