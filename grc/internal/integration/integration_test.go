package integration

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/carlosrabelo/grc/grc/internal/app"
	"github.com/carlosrabelo/grc/grc/internal/gmail"
	"github.com/carlosrabelo/grc/grc/internal/importer"
	"github.com/carlosrabelo/grc/grc/internal/loader"
	"github.com/carlosrabelo/grc/grc/internal/testutils"
)

type gmailFeed struct {
	XMLName xml.Name     `xml:"feed"`
	XMLNS   string       `xml:"xmlns,attr"`
	Entries []gmailEntry `xml:"entry"`
}

type gmailEntry struct {
	Category   gmailCategory   `xml:"category"`
	Title      string          `xml:"title"`
	Properties []gmailProperty `xml:"http://schemas.google.com/apps/2006 property"`
}

type gmailCategory struct {
	Term string `xml:"term,attr"`
}

type gmailProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func TestEndToEnd_GmailXMLStructure(t *testing.T) {
	yaml := `filters:
  - name: Support
    criteria:
      from: support@example.com
      subject: help
      hasAttachment: true
    action:
      archive: true
      label: Support
      star: true`

	path := writeTestYAML(t, yaml)
	cfg, err := loader.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	data, err := gmail.GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	var feed gmailFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("valid XML expected: %v", err)
	}

	if feed.XMLNS != "http://www.w3.org/2005/Atom" {
		t.Errorf("atom ns: got %q", feed.XMLNS)
	}
	if len(feed.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(feed.Entries))
	}

	entry := feed.Entries[0]
	if entry.Title != "Support" {
		t.Errorf("title: got %q", entry.Title)
	}
	if entry.Category.Term != "filter" {
		t.Errorf("category term: got %q", entry.Category.Term)
	}

	names := propertyNames(entry.Properties)
	checks := []string{"from", "subject", "hasAttachment", "shouldArchive", "label", "shouldStar"}
	for _, c := range checks {
		if !contains(names, c) {
			t.Errorf("missing property %q", c)
		}
	}
}

func TestEndToEnd_AllCriteria(t *testing.T) {
	yaml := `filters:
  - name: All
    criteria:
      from: a@b.com
      to: c@d.com
      subject: subj
      hasTheWord: free
      hasAttachment: true
      is: unread
      in: inbox
      cc: e@f.com
      bcc: g@h.com
      deliveredto: i@j.com
      list: k@l.com
      filename: doc.pdf
      after: "2024/01/01"
      before: "2024/12/31"
      older_than: 30d
      newer_than: 7d
      size: 10M
      larger: 5M
      smaller: 1M
      label: MyLabel
      category: social
    action:
      archive: true
      markRead: true
      star: true
      delete: false
      neverSpam: true
      alwaysSpam: false
      forward: fwd@example.com
      label: ActionLabel
      smartLabel: receipts
      category: updates
      moveToTrash: false`

	path := writeTestYAML(t, yaml)
	cfg, err := loader.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	data, err := gmail.GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	var feed gmailFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("valid XML expected: %v", err)
	}

	if len(feed.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(feed.Entries))
	}

	props := feed.Entries[0].Properties
	names := propertyNames(props)
	expectedCriteria := []string{
		"from", "to", "subject", "hasTheWord", "hasAttachment",
		"is", "in", "cc", "bcc", "deliveredto", "list", "filename",
		"after", "before", "olderThan", "newerThan",
		"size", "larger", "smaller", "label", "category",
	}
	expectedActions := []string{
		"shouldArchive", "shouldMarkAsRead", "shouldStar",
		"shouldDelete", "shouldNeverSpam", "shouldAlwaysMarkAsSpam",
		"forwardTo", "label", "smartLabelToApply", "category",
		"shouldMoveToTrash",
	}
	expected := append(expectedCriteria, expectedActions...)
	for _, e := range expected {
		if !contains(names, e) {
			t.Errorf("missing property %q", e)
		}
	}
	if len(props) != len(expected) {
		t.Errorf("expected %d properties, got %d", len(expected), len(props))
	}
}

func TestEndToEnd_Roundtrip(t *testing.T) {
	yaml := `filters:
  - name: Newsletter
    criteria:
      from: newsletter@example.com
      subject: weekly
    action:
      archive: true
      label: News
      markRead: true
  - name: Spam
    criteria:
      from: spam@example.com
    action:
      alwaysSpam: true`

	path := writeTestYAML(t, yaml)
	orig, err := loader.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	xmlData, err := gmail.GenerateGmailXML(orig)
	if err != nil {
		t.Fatal(err)
	}

	xmlPath := filepath.Join(t.TempDir(), "roundtrip.xml")
	if err := os.WriteFile(xmlPath, xmlData, 0644); err != nil {
		t.Fatal(err)
	}

	imported, err := importer.ImportGmailXML(xmlPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(imported.Filters) != len(orig.Filters) {
		t.Fatalf("expected %d filters, got %d", len(orig.Filters), len(imported.Filters))
	}

	for i, f := range imported.Filters {
		if f.Name != orig.Filters[i].Name {
			t.Errorf("filter %d: name got %q, want %q", i, f.Name, orig.Filters[i].Name)
		}
	}
}

func TestEndToEnd_CLIRoundtrip(t *testing.T) {
	yaml := `filters:
  - name: Support
    criteria:
      from: support@example.com
      subject: help
    action:
      label: Support
      archive: true`

	yamlPath := writeTestYAML(t, yaml)
	gmailPath := filepath.Join(t.TempDir(), "out.gmail.xml")

	code := app.Run([]string{"--format", "gmail", yamlPath, gmailPath})
	if code != 0 {
		t.Fatalf("gmail generation failed: %d", code)
	}
	if _, err := os.Stat(gmailPath); os.IsNotExist(err) {
		t.Fatal("gmail XML not created")
	}

	importPath := filepath.Join(t.TempDir(), "reimport.yaml")
	code = app.Run([]string{"--import", gmailPath, importPath})
	if code != 0 {
		t.Fatalf("import failed: %d", code)
	}
	if _, err := os.Stat(importPath); os.IsNotExist(err) {
		t.Fatal("imported YAML not created")
	}

	imported, err := loader.LoadConfig(importPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(imported.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(imported.Filters))
	}
}

func TestEndToEnd_ExampleGmailYAML(t *testing.T) {
	path := filepath.Join("..", "..", "..", "cfg", "example-gmail.yaml")
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		t.Skip("example-gmail.yaml not found, skipping")
	}

	cfg, err := loader.LoadConfig(abs)
	if err != nil {
		t.Fatalf("loading example-gmail.yaml: %v", err)
	}

	if len(cfg.Filters) == 0 {
		t.Fatal("expected at least one filter")
	}

	data, err := gmail.GenerateGmailXML(cfg)
	if err != nil {
		t.Fatalf("generating gmail xml: %v", err)
	}

	var feed gmailFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("valid XML: %v", err)
	}

	if len(feed.Entries) == 0 {
		t.Fatal("expected entries in XML")
	}

	hasCriteria := false
	hasAction := false
	for _, entry := range feed.Entries {
		for _, p := range entry.Properties {
			if isCriteriaProperty(p.Name) {
				hasCriteria = true
			}
			if isActionProperty(p.Name) {
				hasAction = true
			}
		}
	}
	if !hasCriteria {
		t.Error("no criteria properties found in output")
	}
	if !hasAction {
		t.Error("no action properties found in output")
	}
}

func TestEndToEnd_ExampleYAML_ExpressionFallback(t *testing.T) {
	path := filepath.Join("..", "..", "..", "cfg", "example.yaml")
	abs, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		t.Skip("example.yaml not found, skipping")
	}

	cfg, err := loader.LoadConfig(abs)
	if err != nil {
		t.Fatalf("loading example.yaml: %v", err)
	}

	data, err := gmail.GenerateGmailXML(cfg)
	if err != nil {
		t.Fatalf("generating gmail xml from expression-based config: %v", err)
	}

	var feed gmailFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("valid XML: %v", err)
	}

	if len(feed.Entries) == 0 {
		t.Fatal("expected entries")
	}
}

func TestEndToEnd_CLIFormatAtom(t *testing.T) {
	yaml := `filters:
  - name: Test
    criteria:
      from: a@b.com`

	yamlPath := writeTestYAML(t, yaml)
	output := filepath.Join(t.TempDir(), "atom.xml")

	code := app.Run([]string{"--format", "atom", yamlPath, output})
	if code != 0 {
		t.Fatalf("atom generation failed: %d", code)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "<feed") {
		t.Error("expected <feed> in atom output")
	}
	if !strings.Contains(string(data), "<entry>") {
		t.Error("expected <entry> in atom output")
	}
}

func TestEndToEnd_CLIEmptyFilters(t *testing.T) {
	yaml := `filters:`
	yamlPath := writeTestYAML(t, yaml)
	output := filepath.Join(t.TempDir(), "empty.gmail.xml")

	code := app.Run([]string{"--format", "gmail", yamlPath, output})
	if code != 0 {
		t.Fatalf("expected exit 0 for empty filters, got %d", code)
	}
}

func TestEndToEnd_CLIInvalidFormat(t *testing.T) {
	yamlPath := writeTestYAML(t, "filters:\n  - name: test\n    expression: x\n")
	code := app.Run([]string{"--format", "bogus", yamlPath})
	if code != 1 {
		t.Fatalf("expected exit 1 for invalid format, got %d", code)
	}
}

func TestEndToEnd_ActionOnlyFilter(t *testing.T) {
	yaml := `filters:
  - name: AutoLabel
    action:
      label: Auto
      markRead: true`

	yamlPath := writeTestYAML(t, yaml)
	_, err := loader.LoadConfig(yamlPath)
	if err == nil {
		t.Fatal("expected error for action-only filter (no criteria/expression)")
	}
}

func TestEndToEnd_CriteriaOnlyFilter(t *testing.T) {
	yaml := `filters:
  - name: MatchOnly
    criteria:
      from: match@example.com
      subject: alert`

	yamlPath := writeTestYAML(t, yaml)
	cfg, err := loader.LoadConfig(yamlPath)
	if err != nil {
		t.Fatal(err)
	}

	data, err := gmail.GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	var feed gmailFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("valid XML: %v", err)
	}

	names := propertyNames(feed.Entries[0].Properties)
	if !contains(names, "from") {
		t.Error("expected from")
	}
	if !contains(names, "subject") {
		t.Error("expected subject")
	}
}

func propertyNames(props []gmailProperty) []string {
	names := make([]string, len(props))
	for i, p := range props {
		names[i] = p.Name
	}
	return names
}

func contains(list []string, item string) bool {
	for _, s := range list {
		if s == item {
			return true
		}
	}
	return false
}

func isCriteriaProperty(name string) bool {
	criteria := []string{
		"from", "to", "subject", "hasTheWord", "hasAttachment",
		"is", "in", "cc", "bcc", "deliveredto", "list", "filename",
		"after", "before", "olderThan", "newerThan",
		"size", "larger", "smaller",
	}
	for _, c := range criteria {
		if name == c {
			return true
		}
	}
	return false
}

func isActionProperty(name string) bool {
	actions := []string{
		"shouldArchive", "shouldMarkAsRead", "shouldStar",
		"shouldDelete", "shouldNeverSpam", "shouldAlwaysMarkAsSpam",
		"forwardTo", "smartLabelToApply", "shouldMoveToTrash",
	}
	for _, a := range actions {
		if name == a {
			return true
		}
	}
	return false
}

func writeTestYAML(t *testing.T, content string) string {
	t.Helper()
	path, err := testutils.CreateTempYAMLFile("", content)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { testutils.CleanupFile(path) })
	return path
}
