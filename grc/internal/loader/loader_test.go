package loader

import (
	"testing"

	"github.com/carlosrabelo/grc/grc/internal/testutils"
)

func TestLoadConfig_Valid(t *testing.T) {
	content := "filters:\n  - name: newsletter\n    expression: subject:newsletter\n    action:\n      label: Newsletters\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	cfg, err := LoadConfig(path)
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

func TestLoadConfig_MultipleFilters(t *testing.T) {
	content := "filters:\n  - name: a\n    expression: x\n    action:\n      label: Inbox\n  - name: b\n    expression: y\n    action:\n      label: Inbox\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(cfg.Filters))
	}
}

func TestLoadConfig_WithDefaults(t *testing.T) {
	content := `filters:
  - name: test
    expression: x
    action:
      label: Inbox
defaults:
  author:
    name: Default Author
    email: default@example.com
`
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Defaults == nil {
		t.Fatal("expected defaults")
	}
	if cfg.Defaults.Author.Name != "Default Author" {
		t.Errorf("got %q", cfg.Defaults.Author.Name)
	}
}

func TestLoadConfig_WithAuthor(t *testing.T) {
	content := "filters:\n  - name: test\n    expression: x\n    action:\n      label: Inbox\n    author:\n      name: Author\n      email: author@example.com\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Filters[0].Author.Name != "Author" {
		t.Errorf("got %q", cfg.Filters[0].Author.Name)
	}
}

func TestLoadConfig_BadExtension(t *testing.T) {
	path, err := testutils.CreateTempFile("", "*.json", "{}")
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for bad extension")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	path, err := testutils.CreateTempYAMLFile("", "{invalid")
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoadConfig_UnknownField(t *testing.T) {
	content := "filters:\n  - name: test\n    expression: x\n    action:\n      label: Inbox\n    bogus: true\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestLoadConfig_MissingFilterName(t *testing.T) {
	content := "filters:\n  - expression: x\n    action:\n      label: Inbox\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for missing filter name")
	}
}

func TestLoadConfig_MissingExpression(t *testing.T) {
	content := "filters:\n  - name: test\n    action:\n      label: Inbox\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for missing expression")
	}
}

func TestLoadConfig_MissingLabel(t *testing.T) {
	content := "filters:\n  - name: test\n    expression: x\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for missing action.label")
	}
}

func TestLoadConfig_InvalidAuthorEmail(t *testing.T) {
	content := "filters:\n  - name: test\n    expression: x\n    action:\n      label: Inbox\n    author:\n      name: Author\n      email: not-an-email\n"
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid author email")
	}
}

func TestLoadConfig_EmptyFile(t *testing.T) {
	path, err := testutils.CreateTempYAMLFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer testutils.CleanupFile(path)

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for empty config")
	}
}
