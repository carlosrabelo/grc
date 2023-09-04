package validate

import (
	"testing"

	"github.com/carlosrabelo/grc/grc/internal/testutils"
	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func TestExtension(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"yaml", "config.yaml", false},
		{"yml", "config.yml", false},
		{"upper YAML", "config.YAML", false},
		{"no extension", "config", true},
		{"wrong extension", "config.json", true},
		{"txt", "config.yaml.txt", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Extension(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extension(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestYAML_Valid(t *testing.T) {
	yaml := `
filters:
  - name: newsletter
    expression: subject:newsletter
    author:
      name: News Bot
      email: news@example.com
`
	cfg, err := YAML([]byte(yaml))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(cfg.Filters))
	}
	if cfg.Filters[0].Name != "newsletter" {
		t.Errorf("got name %q", cfg.Filters[0].Name)
	}
}

func TestYAML_UnknownField(t *testing.T) {
	yaml := `
filters:
  - name: test
    expression: x
    unknown_field: true
`
	_, err := YAML([]byte(yaml))
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestYAML_InvalidSyntax(t *testing.T) {
	_, err := YAML([]byte("{invalid"))
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestConfigFile_Valid(t *testing.T) {
	content := "filters:\n  - name: test\n    expression: x\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	cfg, err := ConfigFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(cfg.Filters))
	}
}

func TestConfigFile_BadExtension(t *testing.T) {
	path, err := testutils.CreateTempFile("", "*.json", "{}")
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = ConfigFile(path)
	if err == nil {
		t.Fatal("expected error for bad extension")
	}
}

func TestConfigFile_NotFound(t *testing.T) {
	_, err := ConfigFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAuthor_Valid(t *testing.T) {
	err := Author(&types.Author{Name: "John", Email: "john@example.com"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAuthor_EmptyName(t *testing.T) {
	err := Author(&types.Author{Name: "", Email: "john@example.com"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAuthor_InvalidEmail(t *testing.T) {
	err := Author(&types.Author{Name: "John", Email: "not-an-email"})
	if err == nil {
		t.Fatal("expected error for invalid email")
	}
}

func TestAuthor_Nil(t *testing.T) {
	err := Author(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFilter_Valid(t *testing.T) {
	err := Filter(&types.Filter{Name: "test", Expression: "subject:hello"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFilter_EmptyName(t *testing.T) {
	err := Filter(&types.Filter{Name: "", Expression: "x"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestFilter_EmptyExpression(t *testing.T) {
	err := Filter(&types.Filter{Name: "test", Expression: ""})
	if err == nil {
		t.Fatal("expected error for empty expression")
	}
}

func TestFilter_InvalidAuthor(t *testing.T) {
	err := Filter(&types.Filter{
		Name:       "test",
		Expression: "x",
		Author:     &types.Author{Name: "", Email: "bad"},
	})
	if err == nil {
		t.Fatal("expected error for invalid author")
	}
}

func TestEmail_Valid(t *testing.T) {
	if err := Email("user@example.com"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := Email(""); err != nil {
		t.Errorf("unexpected error for empty string: %v", err)
	}
}

func TestEmail_Invalid(t *testing.T) {
	if err := Email("not-email"); err == nil {
		t.Fatal("expected error")
	}
}

func TestDomain_Valid(t *testing.T) {
	if err := Domain("example.com"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := Domain(""); err != nil {
		t.Errorf("unexpected error for empty string: %v", err)
	}
}

func TestDomain_Invalid(t *testing.T) {
	if err := Domain("not-a-domain"); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_Valid(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "a", Expression: "x"},
			{Name: "b", Expression: "y"},
		},
	}
	if err := Config(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfig_Nil(t *testing.T) {
	if err := Config(nil); err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestConfig_InvalidFilter(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "", Expression: "x"},
		},
	}
	if err := Config(cfg); err == nil {
		t.Fatal("expected error for invalid filter")
	}
}

func TestConfig_InvalidDefaultsAuthor(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "a", Expression: "x"},
		},
		Defaults: &types.Defaults{
			Author: &types.Author{Name: "", Email: ""},
		},
	}
	if err := Config(cfg); err == nil {
		t.Fatal("expected error for invalid defaults author")
	}
}

func TestConfigFile_UnknownField(t *testing.T) {
	content := "filters:\n  - name: test\n    expression: x\n    bogus: true\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = ConfigFile(path)
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}
