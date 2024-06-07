package generate

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func TestGenerateFeed_Basic(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "newsletter", Expression: "subject:newsletter"},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if feed.XMLNS != types.AtomNS {
		t.Errorf("got xmlns %q, want %q", feed.XMLNS, types.AtomNS)
	}
	if feed.Title != "GRC Feed" {
		t.Errorf("got title %q, want %q", feed.Title, "GRC Feed")
	}
	if feed.ID == "" {
		t.Error("feed ID should not be empty")
	}
	if feed.Updated == "" {
		t.Error("feed updated should not be empty")
	}
	if len(feed.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(feed.Entries))
	}
}

func TestGenerateFeed_MultipleFilters(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "a", Expression: "x"},
			{Name: "b", Expression: "y"},
			{Name: "c", Expression: "z"},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(feed.Entries))
	}
}

func TestGenerateFeed_EntryFields(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "test-filter", Expression: "subject:hello"},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	entry := feed.Entries[0]
	if entry.Title != "test-filter" {
		t.Errorf("got title %q", entry.Title)
	}
	if !strings.HasPrefix(entry.ID, types.FeedID) {
		t.Errorf("entry ID should start with %q, got %q", types.FeedID, entry.ID)
	}
	if entry.Updated == "" {
		t.Error("entry updated should not be empty")
	}
	if !strings.Contains(entry.Summary, "subject:hello") {
		t.Errorf("summary should contain expression, got %q", entry.Summary)
	}
}

func TestGenerateFeed_WithAuthor(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:       "test",
				Expression: "x",
				Author:     &types.Author{Name: "John", Email: "john@example.com"},
			},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	entry := feed.Entries[0]
	if entry.Author == nil {
		t.Fatal("expected author")
	}
	if entry.Author.Name != "John" {
		t.Errorf("got author name %q", entry.Author.Name)
	}
}

func TestGenerateFeed_WithCategory(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:       "test",
				Expression: "x",
				Category:   &types.Category{Term: "newsletter"},
			},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	entry := feed.Entries[0]
	if len(entry.Categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(entry.Categories))
	}
	if entry.Categories[0].Term != "newsletter" {
		t.Errorf("got term %q", entry.Categories[0].Term)
	}
}

func TestGenerateFeed_WithProperties(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:       "test",
				Expression: "x",
				Properties: []types.Property{
					{Name: "color", Value: "red"},
					{Name: "size", Value: "large"},
				},
			},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	summary := feed.Entries[0].Summary
	if !strings.Contains(summary, "color: red") {
		t.Errorf("summary should contain property, got %q", summary)
	}
	if !strings.Contains(summary, "size: large") {
		t.Errorf("summary should contain property, got %q", summary)
	}
}

func TestGenerateFeed_DefaultsAuthor(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "test", Expression: "x"},
		},
		Defaults: &types.Defaults{
			Author: &types.Author{Name: "Default Author"},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	entry := feed.Entries[0]
	if entry.Author == nil {
		t.Fatal("expected author from defaults")
	}
	if entry.Author.Name != "Default Author" {
		t.Errorf("got %q", entry.Author.Name)
	}
}

func TestGenerateFeed_DefaultsCategory(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "test", Expression: "x"},
		},
		Defaults: &types.Defaults{
			Category: &types.Category{Term: "default-cat"},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	entry := feed.Entries[0]
	if len(entry.Categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(entry.Categories))
	}
	if entry.Categories[0].Term != "default-cat" {
		t.Errorf("got %q", entry.Categories[0].Term)
	}
}

func TestGenerateFeed_ExplicitOverridesDefaults(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:       "test",
				Expression: "x",
				Author:     &types.Author{Name: "Explicit Author"},
				Category:   &types.Category{Term: "explicit-cat"},
			},
		},
		Defaults: &types.Defaults{
			Author:   &types.Author{Name: "Default Author"},
			Category: &types.Category{Term: "default-cat"},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	entry := feed.Entries[0]
	if entry.Author.Name != "Explicit Author" {
		t.Errorf("explicit author should override default, got %q", entry.Author.Name)
	}
	if entry.Categories[0].Term != "explicit-cat" {
		t.Errorf("explicit category should override default, got %q", entry.Categories[0].Term)
	}
}

func TestGenerateFeed_NoDefaults(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "test", Expression: "x"},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	entry := feed.Entries[0]
	if entry.Author != nil {
		t.Error("expected nil author when no defaults")
	}
}

func TestGenerateFeed_XMLMarshaling(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:       "test",
				Expression: "x",
				Author:     &types.Author{Name: "John"},
				Category:   &types.Category{Term: "cat1"},
			},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	out, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	output := string(out)
	if !strings.Contains(output, "<feed") {
		t.Error("output should contain <feed")
	}
	if !strings.Contains(output, "<entry>") {
		t.Error("output should contain <entry>")
	}
	if !strings.Contains(output, "<author>") {
		t.Error("output should contain <author>")
	}
	if !strings.Contains(output, "<category") {
		t.Error("output should contain <category")
	}
}

func TestGenerateFeed_ExplicitBoolProperties(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{
				Name:       "test",
				Expression: "x",
				Properties: []types.Property{
					{Name: "enabled", Value: "true"},
					{Name: "active", Value: "false"},
				},
			},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	summary := feed.Entries[0].Summary
	if !strings.Contains(summary, "enabled: true") {
		t.Errorf("expected bool property in summary, got %q", summary)
	}
	if !strings.Contains(summary, "active: false") {
		t.Errorf("expected bool property in summary, got %q", summary)
	}
}

func TestGenerateFeed_EmptyConfig(t *testing.T) {
	cfg := &types.FiltersConfig{}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(feed.Entries))
	}
}

func TestGenerateFeed_SortsByActionLabel(t *testing.T) {
	cfg := &types.FiltersConfig{
		Filters: []types.Filter{
			{Name: "zeta", Expression: "x", Action: &types.GmailAction{Label: "Zebra"}},
			{Name: "alpha", Expression: "x", Action: &types.GmailAction{Label: "Apple"}},
			{Name: "mid", Expression: "x", Action: &types.GmailAction{Label: "Mango"}},
		},
	}

	feed, err := GenerateFeed(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(feed.Entries))
	}
	wantOrder := []string{"alpha", "mid", "zeta"}
	for i, want := range wantOrder {
		if feed.Entries[i].Title != want {
			t.Errorf("entry %d: got %q, want %q (not sorted by label)", i, feed.Entries[i].Title, want)
		}
	}
}
