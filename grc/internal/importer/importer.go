package importer

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

type importFeed struct {
	XMLName xml.Name      `xml:"feed"`
	Entries []importEntry `xml:"entry"`
}

type importEntry struct {
	Category   importCategory   `xml:"category"`
	Title      string           `xml:"title"`
	Properties []importProperty `xml:"http://schemas.google.com/apps/2006 property"`
}

type importCategory struct {
	Term string `xml:"term,attr"`
}

type importProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

var criteriaFields = map[string]func(*types.GmailCriteria, string){
	"from":        func(c *types.GmailCriteria, v string) { c.From = v },
	"to":          func(c *types.GmailCriteria, v string) { c.To = v },
	"subject":     func(c *types.GmailCriteria, v string) { c.Subject = v },
	"hasTheWord":  func(c *types.GmailCriteria, v string) { c.HasTheWord = v },
	"is":          func(c *types.GmailCriteria, v string) { c.Is = v },
	"in":          func(c *types.GmailCriteria, v string) { c.In = v },
	"cc":          func(c *types.GmailCriteria, v string) { c.Cc = v },
	"bcc":         func(c *types.GmailCriteria, v string) { c.Bcc = v },
	"deliveredto": func(c *types.GmailCriteria, v string) { c.DeliveredTo = v },
	"list":        func(c *types.GmailCriteria, v string) { c.List = v },
	"filename":    func(c *types.GmailCriteria, v string) { c.Filename = v },
	"after":       func(c *types.GmailCriteria, v string) { c.After = v },
	"before":      func(c *types.GmailCriteria, v string) { c.Before = v },
	"olderThan":   func(c *types.GmailCriteria, v string) { c.OlderThan = v },
	"newerThan":   func(c *types.GmailCriteria, v string) { c.NewerThan = v },
	"size":        func(c *types.GmailCriteria, v string) { c.Size = v },
	"larger":      func(c *types.GmailCriteria, v string) { c.Larger = v },
	"smaller":     func(c *types.GmailCriteria, v string) { c.Smaller = v },
}

func ImportGmailXML(path string) (*types.FiltersConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return parseGmailXML(data)
}

func parseGmailXML(data []byte) (*types.FiltersConfig, error) {
	var feed importFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("parsing XML: %w", err)
	}

	cfg := &types.FiltersConfig{}
	for _, entry := range feed.Entries {
		if entry.Category.Term != "filter" {
			continue
		}
		filter := types.Filter{
			Name:     entry.Title,
			Criteria: &types.GmailCriteria{},
		}
		hasCriteria := false
		for _, prop := range entry.Properties {
			if fn, ok := criteriaFields[prop.Name]; ok {
				fn(filter.Criteria, prop.Value)
				hasCriteria = true
				continue
			}
			if prop.Name == "hasAttachment" {
				t := parseBool(prop.Value)
				filter.Criteria.HasAttachment = &t
				hasCriteria = true
			}
		}
		if !hasCriteria {
			filter.Criteria = nil
		}
		cfg.Filters = append(cfg.Filters, filter)
	}
	return cfg, nil
}

func parseBool(v string) bool {
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}
