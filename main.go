package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Defaults struct {
	ShouldArchive              bool `yaml:"shouldArchive"`
	ShouldNeverSpam            bool `yaml:"shouldNeverSpam"`
	ShouldNeverMarkAsImportant bool `yaml:"shouldNeverMarkAsImportant"`
}

type Filter struct {
	From                       string `yaml:"from"`
	Label                      string `yaml:"label"`
	ShouldArchive              *bool  `yaml:"shouldArchive,omitempty"`
	ShouldNeverSpam            *bool  `yaml:"shouldNeverSpam,omitempty"`
	ShouldNeverMarkAsImportant *bool  `yaml:"shouldNeverMarkAsImportant,omitempty"`
}

type Author struct {
	Name  string `yaml:"name" xml:"name"`
	Email string `yaml:"email" xml:"email"`
}

type FiltersConfig struct {
	Author   Author   `yaml:"author"`
	Defaults Defaults `yaml:"default"`
	Filters  []Filter `yaml:"filters"`
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	XMLNS   string   `xml:"xmlns,attr"`
	Apps    string   `xml:"xmlns:apps,attr"`
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Author  Author   `xml:"author"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName    xml.Name   `xml:"entry"`
	Category   Category   `xml:"category"`
	Title      string     `xml:"title"`
	ID         string     `xml:"id"`
	Updated    string     `xml:"updated"`
	Content    string     `xml:"content"`
	Properties []Property `xml:"apps:property"`
}

type Category struct {
	Term string `xml:"term,attr"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <yaml_file>", os.Args[0])
	}

	// Input file name
	yamlFile := os.Args[1]

	// Read YAML file
	file, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var config FiltersConfig
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error decoding YAML: %v", err)
	}

	// Create the XML Feed
	feed := Feed{
		XMLNS:   "http://www.w3.org/2005/Atom",
		Apps:    "http://schemas.google.com/apps/2006",
		Title:   "Mail Filters",
		ID:      "tag:mail.google.com,2008:filters:1234567890",
		Updated: time.Now().UTC().Format(time.RFC3339),
		Author:  config.Author,
	}

	for i, filter := range config.Filters {
		// Set default values if missing in YAML
		if filter.ShouldArchive == nil {
			filter.ShouldArchive = &config.Defaults.ShouldArchive
		}
		if filter.ShouldNeverSpam == nil {
			filter.ShouldNeverSpam = &config.Defaults.ShouldNeverSpam
		}
		if filter.ShouldNeverMarkAsImportant == nil {
			filter.ShouldNeverMarkAsImportant = &config.Defaults.ShouldNeverMarkAsImportant
		}

		entry := Entry{
			Category: Category{Term: "filter"},
			Title:    "Mail Filter",
			ID:       fmt.Sprintf("tag:mail.google.com,2008:filter:z%016d", i+1),
			Updated:  time.Now().UTC().Format(time.RFC3339),
			Content:  "",
			Properties: []Property{
				{Name: "from", Value: filter.From},
				{Name: "label", Value: filter.Label},
				{Name: "shouldArchive", Value: fmt.Sprintf("%v", *filter.ShouldArchive)},
				{Name: "shouldNeverSpam", Value: fmt.Sprintf("%v", *filter.ShouldNeverSpam)},
				{Name: "shouldNeverMarkAsImportant", Value: fmt.Sprintf("%v", *filter.ShouldNeverMarkAsImportant)},
			},
		}
		feed.Entries = append(feed.Entries, entry)
	}

	// Output file name (same as input, with .xml extension)
	xmlFile := strings.TrimSuffix(yamlFile, filepath.Ext(yamlFile)) + ".xml"

	// Generate the XML
	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		log.Fatalf("Error generating XML: %v", err)
	}

	// Save the XML file
	err = os.WriteFile(xmlFile, output, 0644)
	if err != nil {
		log.Fatalf("Error saving XML file: %v", err)
	}

	fmt.Println("XML file successfully generated:", xmlFile)
}
