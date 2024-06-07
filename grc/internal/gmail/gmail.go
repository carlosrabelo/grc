package gmail

import (
	"encoding/xml"
	"fmt"

	"github.com/carlosrabelo/grc/grc/internal/parser"
	"github.com/carlosrabelo/grc/grc/pkg/types"
)

const (
	atomNS         = "http://www.w3.org/2005/Atom"
	appsNS         = "http://schemas.google.com/apps/2006"
	filterCategory = "filter"
)

type gmailFeed struct {
	XMLName xml.Name     `xml:"feed"`
	XMLNS   string       `xml:"xmlns,attr"`
	AppsNS  string       `xml:"xmlns:apps,attr"`
	Entries []gmailEntry `xml:"entry"`
}

type gmailEntry struct {
	XMLName    xml.Name        `xml:"entry"`
	Category   gmailCategory   `xml:"category"`
	Title      string          `xml:"title"`
	Properties []gmailProperty `xml:"apps:property"`
}

type gmailCategory struct {
	Term string `xml:"term,attr"`
}

type gmailProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func GenerateGmailXML(cfg *types.FiltersConfig) ([]byte, error) {
	feed := gmailFeed{
		XMLNS:  atomNS,
		AppsNS: appsNS,
	}

	for _, f := range types.FiltersByLabel(cfg.Filters) {
		entry := gmailEntry{
			Category: gmailCategory{Term: filterCategory},
			Title:    f.Name,
		}
		criteria := f.Criteria
		if criteria == nil && f.Expression != "" {
			criteria = parser.Parse(f.Expression)
		}
		if criteria != nil {
			entry.Properties = append(entry.Properties, criteriaProperties(criteria)...)
		}
		if f.Action != nil {
			entry.Properties = append(entry.Properties, actionProperties(f.Action)...)
		}
		if len(entry.Properties) > 0 {
			feed.Entries = append(feed.Entries, entry)
		}
	}

	out, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling gmail xml: %w", err)
	}
	return append([]byte(xml.Header), out...), nil
}

func criteriaProperties(c *types.GmailCriteria) []gmailProperty {
	return nil
}

func actionProperties(a *types.GmailAction) []gmailProperty {
	return nil
}
