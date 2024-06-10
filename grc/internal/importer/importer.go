package importer

import (
	"encoding/xml"
	"fmt"
	"os"

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
		cfg.Filters = append(cfg.Filters, types.Filter{Name: entry.Title})
	}
	return cfg, nil
}
