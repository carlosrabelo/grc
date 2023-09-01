package testutils

import (
	"os"
)

func CreateTempFile(dir, pattern, content string) (string, error) {
	if dir == "" {
		dir = os.TempDir()
	}
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	if err := f.Close(); err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

func CreateTempYAMLFile(dir, content string) (string, error) {
	return CreateTempFile(dir, "*.yaml", content)
}

func CleanupFile(path string) {
	if path != "" {
		os.Remove(path)
	}
}

func BoolPtr(b bool) *bool {
	return &b
}
