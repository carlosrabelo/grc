package types

import (
	"encoding/xml"
	"regexp"
	"sort"
	"strings"
)

const (
	AtomNS    = "http://www.w3.org/2005/Atom"
	AppsNS    = "http://www.w3.org/2007/app"
	FeedID    = "urn:uuid:"
	XMLHeader = xml.Header
)

var (
	EmailRegex  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	DomainRegex = regexp.MustCompile(`^[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

type FiltersConfig struct {
	Filters  []Filter  `yaml:"filters"`
	Defaults *Defaults `yaml:"defaults,omitempty"`
}

type Filter struct {
	Name       string         `yaml:"name"`
	Expression string         `yaml:"expression,omitempty"`
	Category   *Category      `yaml:"category,omitempty"`
	Properties []Property     `yaml:"properties,omitempty"`
	Author     *Author        `yaml:"author,omitempty"`
	Criteria   *GmailCriteria `yaml:"criteria,omitempty"`
	Action     *GmailAction   `yaml:"action,omitempty"`
}

type GmailCriteria struct {
	From          string `yaml:"from,omitempty"`
	To            string `yaml:"to,omitempty"`
	Subject       string `yaml:"subject,omitempty"`
	HasTheWord    string `yaml:"hasTheWord,omitempty"`
	HasAttachment *bool  `yaml:"hasAttachment,omitempty"`
	Is            string `yaml:"is,omitempty"`
	In            string `yaml:"in,omitempty"`
	Cc            string `yaml:"cc,omitempty"`
	Bcc           string `yaml:"bcc,omitempty"`
	DeliveredTo   string `yaml:"deliveredto,omitempty"`
	List          string `yaml:"list,omitempty"`
	Filename      string `yaml:"filename,omitempty"`
	After         string `yaml:"after,omitempty"`
	Before        string `yaml:"before,omitempty"`
	OlderThan     string `yaml:"older_than,omitempty"`
	NewerThan     string `yaml:"newer_than,omitempty"`
	Size          string `yaml:"size,omitempty"`
	Larger        string `yaml:"larger,omitempty"`
	Smaller       string `yaml:"smaller,omitempty"`
	Label         string `yaml:"label,omitempty"`
	Category      string `yaml:"category,omitempty"`
}

type GmailAction struct {
	Archive     *bool  `yaml:"archive,omitempty"`
	MarkRead    *bool  `yaml:"markRead,omitempty"`
	Star        *bool  `yaml:"star,omitempty"`
	Delete      *bool  `yaml:"delete,omitempty"`
	NeverSpam   *bool  `yaml:"neverSpam,omitempty"`
	AlwaysSpam  *bool  `yaml:"alwaysSpam,omitempty"`
	Forward     string `yaml:"forward,omitempty"`
	Label       string `yaml:"label,omitempty"`
	SmartLabel  string `yaml:"smartLabel,omitempty"`
	Category    string `yaml:"category,omitempty"`
	MoveToTrash *bool  `yaml:"moveToTrash,omitempty"`
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	XMLNS   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title      string     `xml:"title"`
	ID         string     `xml:"id"`
	Updated    string     `xml:"updated"`
	Summary    string     `xml:"summary"`
	Author     *Author    `xml:"author,omitempty"`
	Categories []Category `xml:"category,omitempty"`
}

type Category struct {
	Term   string `xml:"term,attr" yaml:"term"`
	Scheme string `xml:"scheme,attr,omitempty" yaml:"scheme,omitempty"`
	Label  string `xml:"label,attr,omitempty" yaml:"label,omitempty"`
}

type Property struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Author struct {
	Name  string `xml:"name" yaml:"name"`
	Email string `xml:"email,omitempty" yaml:"email,omitempty"`
	URI   string `xml:"uri,omitempty" yaml:"uri,omitempty"`
}

type Defaults struct {
	Author   *Author   `yaml:"author,omitempty"`
	Category *Category `yaml:"category,omitempty"`
}

// FiltersByLabel returns a copy of filters sorted alphabetically by the
// destination label (action.label), case-insensitive and stable. The input
// slice is not modified. Callers should ensure every filter has a non-empty
// action.label via validate.Filter before calling.
func FiltersByLabel(filters []Filter) []Filter {
	out := make([]Filter, len(filters))
	copy(out, filters)
	sort.SliceStable(out, func(i, j int) bool {
		li, lj := labelOrEmpty(out[i]), labelOrEmpty(out[j])
		return strings.ToLower(li) < strings.ToLower(lj)
	})
	return out
}

func labelOrEmpty(f Filter) string {
	if f.Action == nil {
		return ""
	}
	return f.Action.Label
}
