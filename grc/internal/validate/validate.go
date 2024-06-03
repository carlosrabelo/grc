package validate

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func Extension(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".yaml" && ext != ".yml" {
		return fmt.Errorf("invalid extension %q, expected .yaml or .yml", ext)
	}
	return nil
}

func Author(a *types.Author) error {
	if a == nil {
		return nil
	}
	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("author name is required")
	}
	if a.Email != "" && !types.EmailRegex.MatchString(a.Email) {
		return fmt.Errorf("invalid author email %q", a.Email)
	}
	return nil
}

func Filter(f *types.Filter) error {
	if strings.TrimSpace(f.Name) == "" {
		return fmt.Errorf("filter name is required")
	}
	if strings.TrimSpace(f.Expression) == "" {
		return fmt.Errorf("filter %q: must have either expression or criteria", f.Name)
	}
	if f.Author != nil {
		if err := Author(f.Author); err != nil {
			return fmt.Errorf("filter %q: %w", f.Name, err)
		}
	}
	return nil
}

func Email(email string) error {
	if email != "" && !types.EmailRegex.MatchString(email) {
		return fmt.Errorf("invalid email %q", email)
	}
	return nil
}

func Domain(domain string) error {
	if domain != "" && !types.DomainRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain %q", domain)
	}
	return nil
}

func Config(cfg *types.FiltersConfig) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	for i := range cfg.Filters {
		if err := Filter(&cfg.Filters[i]); err != nil {
			return err
		}
	}
	if cfg.Defaults != nil && cfg.Defaults.Author != nil {
		if err := Author(cfg.Defaults.Author); err != nil {
			return fmt.Errorf("defaults author: %w", err)
		}
	}
	return nil
}
