package types

import (
	"encoding/xml"
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"AtomNS", AtomNS, "http://www.w3.org/2005/Atom"},
		{"AppsNS", AppsNS, "http://www.w3.org/2007/app"},
		{"FeedID", FeedID, "urn:uuid:"},
		{"XMLHeader", XMLHeader, xml.Header},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestEmailRegex(t *testing.T) {
	valid := []string{
		"user@example.com",
		"a.b@sub.example.co.uk",
		"user+tag@example.org",
		"user_name@example.com.br",
	}
	invalid := []string{
		"",
		"notanemail",
		"@example.com",
		"user@",
		"user@.com",
		"user@exam ple.com",
		"user@example,com",
	}
	for _, e := range valid {
		if !EmailRegex.MatchString(e) {
			t.Errorf("expected %q to be valid", e)
		}
	}
	for _, e := range invalid {
		if EmailRegex.MatchString(e) {
			t.Errorf("expected %q to be invalid", e)
		}
	}
}

func TestDomainRegex(t *testing.T) {
	valid := []string{
		"example.com",
		"sub.example.co.uk",
		"example.com.br",
		"a-b.example.org",
	}
	invalid := []string{
		"",
		"notadomain",
		".com",
		"example",
		"exam ple.com",
	}
	for _, d := range valid {
		if !DomainRegex.MatchString(d) {
			t.Errorf("expected %q to be valid", d)
		}
	}
	for _, d := range invalid {
		if DomainRegex.MatchString(d) {
			t.Errorf("expected %q to be invalid", d)
		}
	}
}

func TestFeedXMLMarshaling(t *testing.T) {
	f := Feed{
		XMLNS: AtomNS,
		Title: "Test Feed",
		ID:    FeedID + "test",
		Entries: []Entry{
			{
				Title: "Entry 1",
				ID:    FeedID + "entry1",
			},
		},
	}
	out, err := xml.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty XML output")
	}
}

func TestFiltersConfigDefaults(t *testing.T) {
	cfg := FiltersConfig{
		Filters: []Filter{
			{Name: "test", Expression: "subject:hello"},
		},
	}
	if len(cfg.Filters) != 1 {
		t.Fatal("expected 1 filter")
	}
	if cfg.Defaults != nil {
		t.Fatal("expected nil defaults")
	}
}

func TestAuthorValidation(t *testing.T) {
	a := Author{Name: "John Doe", Email: "john@example.com"}
	if a.Name == "" {
		t.Error("author name should not be empty")
	}
	if !EmailRegex.MatchString(a.Email) {
		t.Error("expected valid email")
	}
}

func TestCategoryTags(t *testing.T) {
	c := Category{Term: "test", Scheme: "http://example.com", Label: "Test"}
	if c.Term != "test" {
		t.Errorf("got term %q, want %q", c.Term, "test")
	}
}

func TestProperties(t *testing.T) {
	p := Property{Name: "size", Value: "large"}
	if p.Name != "size" || p.Value != "large" {
		t.Errorf("unexpected Property: %+v", p)
	}
}

func TestDefaults(t *testing.T) {
	author := &Author{Name: "Default Author"}
	d := Defaults{Author: author}
	if d.Author == nil {
		t.Fatal("expected non-nil author")
	}
	if d.Author.Name != "Default Author" {
		t.Errorf("got %q, want %q", d.Author.Name, "Default Author")
	}
}
