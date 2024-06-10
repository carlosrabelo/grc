package gmail

import (
	"encoding/xml"
	"fmt"
	"strconv"

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
	var props []gmailProperty
	add := func(name, value string) {
		if value != "" {
			props = append(props, gmailProperty{Name: name, Value: value})
		}
	}
	addBool := func(name string, val *bool) {
		if val != nil {
			props = append(props, gmailProperty{Name: name, Value: strconv.FormatBool(*val)})
		}
	}
	add("from", c.From)
	add("to", c.To)
	add("subject", c.Subject)
	add("hasTheWord", c.HasTheWord)
	addBool("hasAttachment", c.HasAttachment)
	add("is", c.Is)
	add("in", c.In)
	add("cc", c.Cc)
	add("bcc", c.Bcc)
	add("deliveredto", c.DeliveredTo)
	add("list", c.List)
	add("filename", c.Filename)
	add("after", c.After)
	add("before", c.Before)
	add("olderThan", c.OlderThan)
	add("newerThan", c.NewerThan)
	add("size", c.Size)
	add("larger", c.Larger)
	add("smaller", c.Smaller)
	add("label", c.Label)
	add("category", c.Category)
	return props
}

func actionProperties(a *types.GmailAction) []gmailProperty {
	var props []gmailProperty
	add := func(name, value string) {
		if value != "" {
			props = append(props, gmailProperty{Name: name, Value: value})
		}
	}
	addBool := func(name string, val *bool) {
		if val != nil {
			props = append(props, gmailProperty{Name: name, Value: strconv.FormatBool(*val)})
		}
	}
	addBool("shouldArchive", a.Archive)
	addBool("shouldMarkAsRead", a.MarkRead)
	addBool("shouldStar", a.Star)
	addBool("shouldDelete", a.Delete)
	addBool("shouldNeverSpam", a.NeverSpam)
	addBool("shouldAlwaysMarkAsSpam", a.AlwaysSpam)
	add("forwardTo", a.Forward)
	add("label", a.Label)
	add("smartLabelToApply", a.SmartLabel)
	add("category", a.Category)
	addBool("shouldMoveToTrash", a.MoveToTrash)
	return props
}
