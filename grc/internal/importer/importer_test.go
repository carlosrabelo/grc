package importer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/carlosrabelo/grc/grc/internal/gmail"
	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func TestImportGmailXML_Basic(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>Spam Filter</title>
    <apps:property name="from" value="spam@example.com"/>
    <apps:property name="shouldArchive" value="true"/>
    <apps:property name="label" value="Spam"/>
  </entry>
</feed>`

	cfg, err := parseGmailXML([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(cfg.Filters))
	}

	f := cfg.Filters[0]
	if f.Name != "Spam Filter" {
		t.Errorf("got name %q", f.Name)
	}
	if f.Criteria == nil || f.Criteria.From != "spam@example.com" {
		t.Errorf("got from %v", f.Criteria)
	}
	if f.Action == nil || f.Action.Archive == nil || !*f.Action.Archive {
		t.Error("expected archive true")
	}
	if f.Action.Label != "Spam" {
		t.Errorf("got label %q", f.Action.Label)
	}
}

func TestImportGmailXML_MultipleFilters(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>Filter A</title>
    <apps:property name="from" value="a@example.com"/>
    <apps:property name="shouldArchive" value="true"/>
  </entry>
  <entry>
    <category term="filter"/>
    <title>Filter B</title>
    <apps:property name="subject" value="newsletter"/>
    <apps:property name="label" value="News"/>
  </entry>
</feed>`

	cfg, err := parseGmailXML([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(cfg.Filters))
	}
}

func TestImportGmailXML_AllCriteria(t *testing.T) {
	props := []string{
		`<apps:property name="from" value="from@example.com"/>`,
		`<apps:property name="to" value="to@example.com"/>`,
		`<apps:property name="subject" value="subject text"/>`,
		`<apps:property name="hasTheWord" value="important"/>`,
		`<apps:property name="hasAttachment" value="true"/>`,
		`<apps:property name="is" value="unread"/>`,
		`<apps:property name="in" value="inbox"/>`,
		`<apps:property name="cc" value="cc@example.com"/>`,
		`<apps:property name="bcc" value="bcc@example.com"/>`,
		`<apps:property name="deliveredto" value="delivered@example.com"/>`,
		`<apps:property name="list" value="list@example.com"/>`,
		`<apps:property name="filename" value="report.pdf"/>`,
		`<apps:property name="after" value="2024/01/01"/>`,
		`<apps:property name="before" value="2024/06/01"/>`,
		`<apps:property name="olderThan" value="30d"/>`,
		`<apps:property name="newerThan" value="7d"/>`,
		`<apps:property name="size" value="10M"/>`,
		`<apps:property name="larger" value="5M"/>`,
		`<apps:property name="smaller" value="1M"/>`,
	}
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>All Criteria</title>
    ` + strings.Join(props, "\n    ") + `
  </entry>
</feed>`

	cfg, err := parseGmailXML([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}

	c := cfg.Filters[0].Criteria
	if c == nil {
		t.Fatal("expected criteria")
	}
	if c.From != "from@example.com" { t.Errorf("From: got %q", c.From) }
	if c.To != "to@example.com" { t.Errorf("To: got %q", c.To) }
	if c.Subject != "subject text" { t.Errorf("Subject: got %q", c.Subject) }
	if c.HasTheWord != "important" { t.Errorf("HasTheWord: got %q", c.HasTheWord) }
	if c.HasAttachment == nil || !*c.HasAttachment { t.Error("HasAttachment expected true") }
	if c.Is != "unread" { t.Errorf("Is: got %q", c.Is) }
	if c.In != "inbox" { t.Errorf("In: got %q", c.In) }
	if c.Cc != "cc@example.com" { t.Errorf("Cc: got %q", c.Cc) }
	if c.Bcc != "bcc@example.com" { t.Errorf("Bcc: got %q", c.Bcc) }
	if c.DeliveredTo != "delivered@example.com" { t.Errorf("DeliveredTo: got %q", c.DeliveredTo) }
	if c.List != "list@example.com" { t.Errorf("List: got %q", c.List) }
	if c.Filename != "report.pdf" { t.Errorf("Filename: got %q", c.Filename) }
	if c.After != "2024/01/01" { t.Errorf("After: got %q", c.After) }
	if c.Before != "2024/06/01" { t.Errorf("Before: got %q", c.Before) }
	if c.OlderThan != "30d" { t.Errorf("OlderThan: got %q", c.OlderThan) }
	if c.NewerThan != "7d" { t.Errorf("NewerThan: got %q", c.NewerThan) }
	if c.Size != "10M" { t.Errorf("Size: got %q", c.Size) }
	if c.Larger != "5M" { t.Errorf("Larger: got %q", c.Larger) }
	if c.Smaller != "1M" { t.Errorf("Smaller: got %q", c.Smaller) }
}

func TestImportGmailXML_AllActions(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>All Actions</title>
    <apps:property name="from" value="x@y.com"/>
    <apps:property name="shouldArchive" value="true"/>
    <apps:property name="shouldMarkAsRead" value="true"/>
    <apps:property name="shouldStar" value="false"/>
    <apps:property name="shouldDelete" value="false"/>
    <apps:property name="shouldNeverSpam" value="true"/>
    <apps:property name="shouldAlwaysMarkAsSpam" value="false"/>
    <apps:property name="forwardTo" value="backup@example.com"/>
    <apps:property name="label" value="Auto-Label"/>
    <apps:property name="smartLabelToApply" value="Smart Label"/>
    <apps:property name="shouldMoveToTrash" value="true"/>
  </entry>
</feed>`

	cfg, err := parseGmailXML([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}

	a := cfg.Filters[0].Action
	if a == nil {
		t.Fatal("expected action")
	}
	if a.Archive == nil || !*a.Archive { t.Error("expected Archive true") }
	if a.MarkRead == nil || !*a.MarkRead { t.Error("expected MarkRead true") }
	if a.Star == nil || *a.Star { t.Error("expected Star false") }
	if a.Delete == nil || *a.Delete { t.Error("expected Delete false") }
	if a.NeverSpam == nil || !*a.NeverSpam { t.Error("expected NeverSpam true") }
	if a.AlwaysSpam == nil || *a.AlwaysSpam { t.Error("expected AlwaysSpam false") }
	if a.Forward != "backup@example.com" { t.Errorf("Forward: got %q", a.Forward) }
	if a.Label != "Auto-Label" { t.Errorf("Label: got %q", a.Label) }
	if a.SmartLabel != "Smart Label" { t.Errorf("SmartLabel: got %q", a.SmartLabel) }
	if a.MoveToTrash == nil || !*a.MoveToTrash { t.Error("expected MoveToTrash true") }
}

func TestImportGmailXML_Roundtrip(t *testing.T) {
	original := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name: "Newsletter",
				Criteria: &types.GmailCriteria{
					From:    "news@example.com",
					Subject: "weekly",
				},
				Action: &types.GmailAction{
					Archive:  boolPtr(true),
					Label:    "News",
					MarkRead: boolPtr(true),
				},
			},
			{
				Name: "Spam",
				Criteria: &types.GmailCriteria{
					From: "spam@example.com",
				},
				Action: &types.GmailAction{
					AlwaysSpam: boolPtr(true),
				},
			},
		},
	}

	xmlData, err := gmail.GenerateGmailXML(original)
	if err != nil {
		t.Fatal(err)
	}

	imported, err := parseGmailXML(xmlData)
	if err != nil {
		t.Fatal(err)
	}

	if len(imported.Filters) != len(original.Filters) {
		t.Fatalf("expected %d filters, got %d", len(original.Filters), len(imported.Filters))
	}

	for i, f := range imported.Filters {
		if f.Name != original.Filters[i].Name {
			t.Errorf("filter %d: name got %q, want %q", i, f.Name, original.Filters[i].Name)
		}
	}
}

func TestImportGmailXML_FromFile(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>Test</title>
    <apps:property name="from" value="test@example.com"/>
  </entry>
</feed>`

	path := filepath.Join(t.TempDir(), "filters.xml")
	if err := os.WriteFile(path, []byte(xml), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := ImportGmailXML(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(cfg.Filters))
	}
}

func TestImportGmailXML_InvalidFile(t *testing.T) {
	_, err := ImportGmailXML("/nonexistent/file.xml")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestImportGmailXML_InvalidXML(t *testing.T) {
	_, err := parseGmailXML([]byte("{invalid"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestImportGmailXML_NonFilterEntry(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="other"/>
    <title>Not a filter</title>
    <apps:property name="from" value="x@y.com"/>
  </entry>
</feed>`

	cfg, err := parseGmailXML([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Filters) != 0 {
		t.Errorf("expected 0 filters, got %d", len(cfg.Filters))
	}
}

func TestImportGmailXML_EmptyFeed(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
</feed>`

	cfg, err := parseGmailXML([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Filters) != 0 {
		t.Errorf("expected 0 filters, got %d", len(cfg.Filters))
	}
}

func TestParseGmailXML_CategoryAndLabelAmbiguous(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:apps="http://schemas.google.com/apps/2006">
  <entry>
    <category term="filter"/>
    <title>Ambiguous</title>
    <apps:property name="label" value="MyLabel"/>
    <apps:property name="category" value="updates"/>
  </entry>
</feed>`

	cfg, err := parseGmailXML([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}

	f := cfg.Filters[0]
	if f.Criteria.Label != "MyLabel" {
		t.Errorf("criteria label: got %q", f.Criteria.Label)
	}
	if f.Action.Label != "MyLabel" {
		t.Errorf("action label: got %q", f.Action.Label)
	}
	if f.Criteria.Category != "updates" {
		t.Errorf("criteria category: got %q", f.Criteria.Category)
	}
	if f.Action.Category != "updates" {
		t.Errorf("action category: got %q", f.Action.Category)
	}
}
