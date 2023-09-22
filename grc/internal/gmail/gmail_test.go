package gmail

import (
	"strings"
	"testing"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func boolPtr(b bool) *bool { return &b }

func TestGenerateGmailXML_Basic(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name: "Test Filter",
				Criteria: &types.GmailCriteria{
					From:    "spam@example.com",
					Subject: "newsletter",
				},
				Action: &types.GmailAction{
					Archive: boolPtr(true),
					Label:   "Newsletters",
				},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	output := string(data)
	if !strings.Contains(output, "<feed") {
		t.Error("expected feed element")
	}
	if !strings.Contains(output, "<entry>") {
		t.Error("expected entry element")
	}
	if !strings.Contains(output, "xmlns:apps=") {
		t.Error("expected apps namespace")
	}
	if !strings.Contains(output, `term="filter"`) {
		t.Error("expected category term=filter")
	}
}

func TestGenerateGmailXML_CriteriaFrom(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name: "From Filter",
				Criteria: &types.GmailCriteria{
					From: "spam@example.com",
				},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), `name="from" value="spam@example.com"`) {
		t.Error("expected from property")
	}
}

func TestGenerateGmailXML_AllCriteria(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name: "All Criteria",
				Criteria: &types.GmailCriteria{
					From:          "from@example.com",
					To:            "to@example.com",
					Subject:       "subject text",
					HasTheWord:    "important",
					HasAttachment: boolPtr(true),
					Is:            "unread",
					In:            "inbox",
					Cc:            "cc@example.com",
					Bcc:           "bcc@example.com",
					DeliveredTo:   "delivered@example.com",
					List:          "list@example.com",
					Filename:      "report.pdf",
					After:         "2024/01/01",
					Before:        "2024/06/01",
					OlderThan:     "30d",
					NewerThan:     "7d",
					Size:          "10M",
					Larger:        "5M",
					Smaller:       "1M",
					Label:         "Important",
					Category:      "updates",
				},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	output := string(data)
	checks := []string{
		`name="from" value="from@example.com"`,
		`name="to" value="to@example.com"`,
		`name="subject" value="subject text"`,
		`name="hasTheWord" value="important"`,
		`name="hasAttachment" value="true"`,
		`name="is" value="unread"`,
		`name="in" value="inbox"`,
		`name="cc" value="cc@example.com"`,
		`name="bcc" value="bcc@example.com"`,
		`name="deliveredto" value="delivered@example.com"`,
		`name="list" value="list@example.com"`,
		`name="filename" value="report.pdf"`,
		`name="after" value="2024/01/01"`,
		`name="before" value="2024/06/01"`,
		`name="olderThan" value="30d"`,
		`name="newerThan" value="7d"`,
		`name="size" value="10M"`,
		`name="larger" value="5M"`,
		`name="smaller" value="1M"`,
		`name="label" value="Important"`,
		`name="category" value="updates"`,
	}
	for _, c := range checks {
		if !strings.Contains(output, c) {
			t.Errorf("missing: %s", c)
		}
	}
}

func TestGenerateGmailXML_AllActions(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name: "All Actions",
				Criteria: &types.GmailCriteria{
					From: "test@example.com",
				},
				Action: &types.GmailAction{
					Archive:     boolPtr(true),
					MarkRead:    boolPtr(true),
					Star:        boolPtr(false),
					Delete:      boolPtr(false),
					NeverSpam:   boolPtr(true),
					AlwaysSpam:  boolPtr(false),
					Forward:     "backup@example.com",
					Label:       "Auto-Label",
					SmartLabel:  "Smart Label",
					Category:    "updates",
					MoveToTrash: boolPtr(true),
				},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	output := string(data)
	checks := []string{
		`name="shouldArchive" value="true"`,
		`name="shouldMarkAsRead" value="true"`,
		`name="shouldStar" value="false"`,
		`name="shouldDelete" value="false"`,
		`name="shouldNeverSpam" value="true"`,
		`name="shouldAlwaysMarkAsSpam" value="false"`,
		`name="forwardTo" value="backup@example.com"`,
		`name="label" value="Auto-Label"`,
		`name="smartLabelToApply" value="Smart Label"`,
		`name="category" value="updates"`,
		`name="shouldMoveToTrash" value="true"`,
	}
	for _, c := range checks {
		if !strings.Contains(output, c) {
			t.Errorf("missing: %s", c)
		}
	}
}

func TestGenerateGmailXML_MultipleFilters(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:     "Filter A",
				Criteria: &types.GmailCriteria{From: "a@example.com"},
				Action:   &types.GmailAction{Archive: boolPtr(true)},
			},
			{
				Name:     "Filter B",
				Criteria: &types.GmailCriteria{Subject: "newsletter"},
				Action:   &types.GmailAction{Label: "Newsletters"},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	output := string(data)
	if !strings.Contains(output, "Filter A") {
		t.Error("expected Filter A title")
	}
	if !strings.Contains(output, "Filter B") {
		t.Error("expected Filter B title")
	}

	entries := strings.Count(output, "<entry>")
	if entries != 2 {
		t.Errorf("expected 2 entries, got %d", entries)
	}
}

func TestGenerateGmailXML_NoCriteria(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "Empty"},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(data), "<entry>") {
		t.Error("expected no entries for filters without criteria")
	}
}

func TestGenerateGmailXML_ActionOnly(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name: "Action Only",
				Action: &types.GmailAction{
					Label: "Test",
				},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "<entry>") {
		t.Error("expected entry for action-only filter")
	}
	if !strings.Contains(string(data), `name="label" value="Test"`) {
		t.Error("expected label property")
	}
}

func TestGenerateGmailXML_EmptyConfig(t *testing.T) {
	cfg := &types.FiltersConfig{}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(data), "<entry>") {
		t.Error("expected no entries for empty config")
	}
}

func TestGenerateGmailXML_Header(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:     "Test",
				Criteria: &types.GmailCriteria{From: "x@y.com"},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(string(data), "<?xml ") {
		t.Error("expected XML header")
	}
}

func TestGenerateGmailXML_XMLNamespace(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:     "Test",
				Criteria: &types.GmailCriteria{From: "x@y.com"},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	output := string(data)
	if !strings.Contains(output, atomNS) {
		t.Error("expected Atom namespace")
	}
	if !strings.Contains(output, appsNS) {
		t.Error("expected Apps namespace")
	}
}

func TestGenerateGmailXML_FilterNameAsTitle(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:     "My Custom Filter",
				Criteria: &types.GmailCriteria{From: "x@y.com"},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "<title>My Custom Filter</title>") {
		t.Error("expected filter name as title")
	}
}

func TestGenerateGmailXML_CategoryTerm(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:     "Test",
				Criteria: &types.GmailCriteria{From: "x@y.com"},
			},
		},
	}

	data, err := GenerateGmailXML(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), `term="filter"`) {
		t.Error("expected category term=filter")
	}
}
