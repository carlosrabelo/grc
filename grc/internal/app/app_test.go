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

func TestRun_FormatGmail(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    criteria:\n      from: spam@example.com\n    action:\n      archive: true\n")
	output := filepath.Join(t.TempDir(), "out.gmail.xml")

	code := Run([]string{"--format", "gmail", input, output})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}

func TestRun_FormatGmailDefaultOutput(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    criteria:\n      from: a@b.com\n")
	expected := input[:len(input)-5] + ".gmail.xml"
	t.Cleanup(func() { os.Remove(expected) })

	code := Run([]string{"--format", "gmail", input})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Error("expected .gmail.xml output file")
	}
}

func TestRun_FormatAtomExplicit(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n")
	output := filepath.Join(t.TempDir(), "explicit.xml")

	code := Run([]string{"--format", "atom", input, output})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("expected output file to exist")
	}
}

func TestRun_FormatUnknown(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n")

	code := Run([]string{"--format", "pdf", input})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestRun_GmailCriteriaAction(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: support\n    criteria:\n      subject: help\n    action:\n      label: Support\n      archive: true\n")
	output := filepath.Join(t.TempDir(), "support.xml")

	code := Run([]string{"--format", "gmail", input, output})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
}

func TestRun_FormatGmailExpressionFallsBack(t *testing.T) {
	input := writeTestYAML(t, "filters:\n  - name: test\n    expression: subject:hello\n")
	output := filepath.Join(t.TempDir(), "expr.xml")

	code := Run([]string{"--format", "gmail", input, output})
	if code != 0 {
		t.Fatalf("expected exit 0 for expression-based filter in gmail mode, got %d", code)
	}
}

func writeTestXML(t *testing.T, content string) string {
	t.Helper()
	path, err := testutils.CreateTempFile("", "*.xml", content)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testutils.CleanupFile(path) })
	return path
}

func TestRun_Import(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>Spam</title>
    <apps:property name="from" value="spam@example.com"/>
    <apps:property name="shouldArchive" value="true"/>
  </entry>
</feed>`
	input := writeTestXML(t, xml)
	output := filepath.Join(t.TempDir(), "result.yaml")

	code := Run([]string{"--import", input, output})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty YAML output")
	}
}

func TestRun_ImportDefaultOutput(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>Test</title>
    <apps:property name="from" value="a@b.com"/>
  </entry>
</feed>`
	input := writeTestXML(t, xml)
	expected := input[:len(input)-4] + ".yaml"
	t.Cleanup(func() { os.Remove(expected) })

	code := Run([]string{"--import", input})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Error("expected .yaml output file")
	}
}

func TestRun_ImportInvalidPath(t *testing.T) {
	code := Run([]string{"--import", "/nonexistent/file.xml"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestRun_ImportNonFilterEntry(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="other"/>
    <title>Not a filter</title>
  </entry>
</feed>`
	input := writeTestXML(t, xml)
	output := filepath.Join(t.TempDir(), "out.yaml")

	code := Run([]string{"--import", input, output})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected output")
	}
}

func TestRun_ImportOverwriteProtected(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>Test</title>
  </entry>
</feed>`
	input := writeTestXML(t, xml)
	output := filepath.Join(t.TempDir(), "prot.yaml")

	code := Run([]string{"--import", input, output})
	if code != 0 {
		t.Fatalf("first run failed: %d", code)
	}

	code = Run([]string{"--import", input, output})
	if code == 0 {
		t.Fatal("expected error for overwrite without --force")
	}
}

func TestRun_ImportForce(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>Test</title>
  </entry>
</feed>`
	input := writeTestXML(t, xml)
	output := filepath.Join(t.TempDir(), "force.yaml")

	code := Run([]string{"--import", input, output})
	if code != 0 {
		t.Fatalf("first run failed: %d", code)
	}

	code = Run([]string{"--import", "--force", input, output})
	if code != 0 {
		t.Fatalf("expected --force to succeed, got %d", code)
	}
}

func TestRun_ImportRoundtrip(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>Support</title>
    <apps:property name="subject" value="help"/>
    <apps:property name="shouldMarkAsRead" value="true"/>
    <apps:property name="label" value="Support"/>
  </entry>
</feed>`
	input := writeTestXML(t, xml)
	yamlPath := filepath.Join(t.TempDir(), "roundtrip.yaml")

	code := Run([]string{"--import", input, yamlPath})
	if code != 0 {
		t.Fatalf("import failed: %d", code)
	}

	gmailPath := filepath.Join(t.TempDir(), "roundtrip.gmail.xml")
	code = Run([]string{"--format", "gmail", yamlPath, gmailPath})
	if code != 0 {
		t.Fatalf("generate failed: %d", code)
	}

	code = Run([]string{"--import", "--force", gmailPath, yamlPath})
	if code != 0 {
		t.Fatalf("re-import failed: %d", code)
	}
}
