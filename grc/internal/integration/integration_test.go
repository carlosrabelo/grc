package integration

import (
	"encoding/xml"
	"testing"

	"github.com/carlosrabelo/grc/grc/internal/gmail"
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
