package app

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

func fileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		t.Errorf("expected output file %q to exist", path)
	}
}

func TestRun_NoArgs_ShowsRootHelp(t *testing.T) {
	code := Run(nil)
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	code := Run([]string{"bogus"})
	if code != exitUsage {
		t.Fatalf("expected exit %d, got %d", exitUsage, code)
	}
}

func TestRun_HelpTopLevel(t *testing.T) {
	code := Run([]string{"--help"})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}

func TestRun_HelpCommand(t *testing.T) {
	code := Run([]string{"help"})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}

func TestRun_HelpUnknownSubcommand(t *testing.T) {
	code := Run([]string{"help", "bogus"})
	if code != exitUsage {
		t.Fatalf("expected exit %d, got %d", exitUsage, code)
	}
}

func TestRun_Version(t *testing.T) {
	Version = "1.0.0"
	BuildTime = "2024-01-01T00:00:00Z"

	code := Run([]string{"version"})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}

func TestRun_VersionHelp(t *testing.T) {
	code := Run([]string{"version", "--help"})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}

func TestRun_AtomBasic(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n    action:\n      label: Inbox\n")
	output := filepath.Join(t.TempDir(), "out.xml")

	code := Run([]string{"atom", input, output})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
	fileExists(t, output)
}

func TestRun_AtomDefaultOutput(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n    action:\n      label: Inbox\n")
	expected := strings.TrimSuffix(input, filepath.Ext(input)) + ".xml"
	t.Cleanup(func() { _ = os.Remove(expected) })

	code := Run([]string{"atom", input})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
	fileExists(t, expected)
}

func TestRun_AtomOutputFlag(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n    action:\n      label: Inbox\n")
	output := filepath.Join(t.TempDir(), "flag.xml")

	code := Run([]string{"atom", "--output", output, input})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
	fileExists(t, output)
}

func TestRun_AtomVerbose(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n    action:\n      label: Inbox\n")
	output := filepath.Join(t.TempDir(), "verbose.xml")

	code := Run([]string{"atom", "--verbose", input, output})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}

func TestRun_AtomMultipleFilters(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: a\n    expression: x\n    action:\n      label: Inbox\n  - name: b\n    expression: y\n    action:\n      label: Inbox\n")
	output := filepath.Join(t.TempDir(), "multi.xml")

	code := Run([]string{"atom", input, output})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestRun_AtomMissingInput(t *testing.T) {
	code := Run([]string{"atom"})
	if code != exitUsage {
		t.Fatalf("expected exit %d, got %d", exitUsage, code)
	}
}

func TestRun_AtomInvalidYAML(t *testing.T) {
	input := writeTestYAML(t, "{invalid")
	output := filepath.Join(t.TempDir(), "out.xml")

	code := Run([]string{"atom", input, output})
	if code != exitError {
		t.Fatalf("expected exit %d, got %d", exitError, code)
	}
}

func TestRun_AtomInvalidInputPath(t *testing.T) {
	code := Run([]string{"atom", "/nonexistent/file.yaml"})
	if code != exitError {
		t.Fatalf("expected exit %d, got %d", exitError, code)
	}
}

func TestRun_AtomHelpFlag(t *testing.T) {
	code := Run([]string{"atom", "--help"})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}

func TestRun_AtomUnknownFlag(t *testing.T) {
	code := Run([]string{"atom", "--bogus"})
	if code != exitUsage {
		t.Fatalf("expected exit %d, got %d", exitUsage, code)
	}
}

func TestRun_AtomOverwriteProtected(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n    action:\n      label: Inbox\n")
	output := filepath.Join(t.TempDir(), "prot.xml")

	if code := Run([]string{"atom", input, output}); code != exitOK {
		t.Fatalf("first run failed: %d", code)
	}
	if code := Run([]string{"atom", input, output}); code != exitError {
		t.Fatalf("expected exit %d without force, got %d", exitError, code)
	}
}

func TestRun_AtomForce(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n    action:\n      label: Inbox\n")
	output := filepath.Join(t.TempDir(), "force.xml")

	if code := Run([]string{"atom", input, output}); code != exitOK {
		t.Fatalf("first run failed: %d", code)
	}
	if code := Run([]string{"atom", "--force", input, output}); code != exitOK {
		t.Fatalf("expected exit %d with force, got %d", exitOK, code)
	}
}

func TestRun_HelpForSubcommand(t *testing.T) {
	for _, sub := range []string{"atom", "gmail", "version", "help"} {
		code := Run([]string{"help", sub})
		if code != exitOK {
			t.Fatalf("help %s: expected exit %d, got %d", sub, exitOK, code)
		}
	}
}

func TestRun_GmailBasic(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    criteria:\n      from: spam@example.com\n    action:\n      label: Spam\n      archive: true\n")
	output := filepath.Join(t.TempDir(), "out.gmail.xml")

	code := Run([]string{"gmail", input, output})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
	fileExists(t, output)
}

func TestRun_GmailDefaultOutput(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    criteria:\n      from: a@b.com\n    action:\n      label: Inbox\n")
	expected := strings.TrimSuffix(input, filepath.Ext(input)) + ".gmail.xml"
	t.Cleanup(func() { _ = os.Remove(expected) })

	code := Run([]string{"gmail", input})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
	fileExists(t, expected)
}

func TestRun_GmailExpressionFallback(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: subject:hello\n    action:\n      label: Inbox\n")
	output := filepath.Join(t.TempDir(), "expr.xml")

	code := Run([]string{"gmail", input, output})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}

func TestRun_GmailCriteriaAction(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: support\n    criteria:\n      subject: help\n    action:\n      label: Support\n      archive: true\n")
	output := filepath.Join(t.TempDir(), "support.xml")

	code := Run([]string{"gmail", input, output})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}

func TestRun_GmailHelpFlag(t *testing.T) {
	code := Run([]string{"gmail", "--help"})
	if code != exitOK {
		t.Fatalf("expected exit %d, got %d", exitOK, code)
	}
}
