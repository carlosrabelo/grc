package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/carlosrabelo/grc/grc/internal/testutils"
)

func writeTestYAML(t *testing.T, content string) string {
	t.Helper()
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testutils.CleanupFile(path) })
	return path
}

func TestRun_Basic(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n")
	output := filepath.Join(t.TempDir(), "out.xml")

	code := Run([]string{input, output})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}

func TestRun_CustomOutput(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n")
	output := filepath.Join(t.TempDir(), "custom.xml")

	code := Run([]string{"--verbose", input, output})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}

func TestRun_DefaultOutput(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n")
	expected := input[:len(input)-5] + ".xml"
	t.Cleanup(func() { os.Remove(expected) })

	code := Run([]string{input})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}

func TestRun_Force(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n")
	output := filepath.Join(t.TempDir(), "force.xml")

	code := Run([]string{input, output})
	if code != 0 {
		t.Fatalf("first run failed: %d", code)
	}

	code = Run([]string{"--force", input, output})
	if code != 0 {
		t.Fatalf("expected force to succeed, got %d", code)
	}
}

func TestRun_OverwriteProtected(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n")
	output := filepath.Join(t.TempDir(), "noforce.xml")

	code := Run([]string{input, output})
	if code != 0 {
		t.Fatalf("first run failed: %d", code)
	}

	code = Run([]string{input, output})
	if code == 0 {
		t.Fatal("expected error for overwrite without force")
	}
}

func TestRun_Version(t *testing.T) {
	Version = "1.0.0"
	BuildTime = "2024-01-01T00:00:00Z"

	code := Run([]string{"--version"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

func TestRun_Help(t *testing.T) {
	code := Run([]string{"--help"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

func TestRun_NoArgs(t *testing.T) {
	code := Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestRun_InvalidInput(t *testing.T) {
	code := Run([]string{"/nonexistent/file.yaml"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestRun_InvalidYAML(t *testing.T) {
	input := writeTestYAML(t, "{invalid")
	output := filepath.Join(t.TempDir(), "out.xml")

	code := Run([]string{input, output})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestRun_Verbose(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n")
	output := filepath.Join(t.TempDir(), "verbose.xml")

	code := Run([]string{"--verbose", input, output})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

func TestRun_MultipleFilters(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: a\n    expression: x\n  - name: b\n    expression: y\n")
	output := filepath.Join(t.TempDir(), "multi.xml")

	code := Run([]string{"--force", input, output})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRun_UnknownFlag(t *testing.T) {
	code := Run([]string{"--bogus"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}
