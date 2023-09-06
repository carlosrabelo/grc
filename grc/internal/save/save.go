package save

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SaveXML(data []byte, path string, force bool) error {
	path = ensureXMLExtension(path)
	if err := validateFileOverwrite(path, force); err != nil {
		return err
	}
	return writeXMLFile(data, path)
}

func ensureXMLExtension(path string) string {
	if strings.ToLower(filepath.Ext(path)) != ".xml" {
		return path + ".xml"
	}
	return path
}

func validateFileOverwrite(path string, force bool) error {
	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("file %q already exists (use --force to overwrite)", path)
	}
	return nil
}

func writeXMLFile(data []byte, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}
