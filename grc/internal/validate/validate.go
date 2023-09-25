package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func Extension(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".yaml" && ext != ".yml" {
		return fmt.Errorf("invalid extension %q, expected .yaml or .yml", ext)
	}
	return nil
}

func YAML(content []byte) (*types.FiltersConfig, error) {
	var cfg types.FiltersConfig
	decoder := yaml.NewDecoder(strings.NewReader(string(content)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("YAML parse error: %w", err)
	}
	return &cfg, nil
}

func ConfigFile(path string) (*types.FiltersConfig, error) {
	if err := Extension(path); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	cfg, err := YAML(data)
	if err != nil {
		return nil, err
	}
	return cfg, nil
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
	hasExpression := strings.TrimSpace(f.Expression) != ""
	hasCriteria := f.Criteria != nil && (f.Criteria.From != "" || f.Criteria.To != "" ||
		f.Criteria.Subject != "" || f.Criteria.HasTheWord != "" ||
		f.Criteria.HasAttachment != nil || f.Criteria.Is != "" ||
		f.Criteria.In != "" || f.Criteria.Cc != "" ||
		f.Criteria.Bcc != "" || f.Criteria.DeliveredTo != "" ||
		f.Criteria.List != "" || f.Criteria.Filename != "" ||
		f.Criteria.After != "" || f.Criteria.Before != "" ||
		f.Criteria.OlderThan != "" || f.Criteria.NewerThan != "" ||
		f.Criteria.Size != "" || f.Criteria.Larger != "" ||
		f.Criteria.Smaller != "" || f.Criteria.Label != "" ||
		f.Criteria.Category != "")
	if !hasExpression && !hasCriteria {
		return fmt.Errorf("filter %q: must have either expression or criteria", f.Name)
	}
	if f.Action == nil || strings.TrimSpace(f.Action.Label) == "" {
		return fmt.Errorf("filter %q: action.label is required (destination folder)", f.Name)
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
