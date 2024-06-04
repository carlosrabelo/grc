package save

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SaveFile(data []byte, path, ext string, force bool) (string, error) {
	path = ensureExtension(path, ext)
	if err := checkOverwrite(path, force); err != nil {
		return path, err
	}
	return path, writeFile(data, path)
}

func SaveXML(data []byte, path string, force bool) error {
	_, err := SaveFile(data, path, ".xml", force)
	return err
}

func ensureExtension(path, ext string) string {
	if ext == "" {
		return path
	}
	if !strings.EqualFold(filepath.Ext(path), ext) {
		return path + ext
	}
	return path
}

func checkOverwrite(path string, force bool) error {
	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("file %q already exists (use --force to overwrite)", path)
	}
	return nil
}

func writeFile(data []byte, path string) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}
