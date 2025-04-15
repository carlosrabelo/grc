package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	AtomNS    = "http://www.w3.org/2005/Atom"
	AppsNS    = "http://schemas.google.com/apps/2006"
	FeedID    = "tag:mail.google.com,2008:filters:%d"
	XMLHeader = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
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

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func loadConfig(filePath string) (FiltersConfig, error) {
	if !strings.HasSuffix(strings.ToLower(filePath), ".yaml") && !strings.HasSuffix(strings.ToLower(filePath), ".yml") {
		return FiltersConfig{}, fmt.Errorf("input file %s must have .yaml or .yml extension", filePath)
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return FiltersConfig{}, fmt.Errorf("reading file: %w", err)
	}

	var config FiltersConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		return FiltersConfig{}, fmt.Errorf("decoding YAML: %w", err)
	}

	for i, filter := range config.Filters {
		if filter.From == "" || filter.Label == "" {
			return FiltersConfig{}, fmt.Errorf("filter %d missing required fields: from=%q, label=%q", i, filter.From, filter.Label)
		}
	}

	return config, nil
}

func normalizeFilter(filter Filter, defaults Defaults) Filter {
	if filter.ShouldArchive == nil {
		filter.ShouldArchive = &defaults.ShouldArchive
	}
	if filter.ShouldNeverSpam == nil {
		filter.ShouldNeverSpam = &defaults.ShouldNeverSpam
	}
	if filter.ShouldNeverMarkAsImportant == nil {
		filter.ShouldNeverMarkAsImportant = &defaults.ShouldNeverMarkAsImportant
	}
	return filter
}

func generateFeed(config FiltersConfig) Feed {
	feed := Feed{
		XMLNS:   AtomNS,
		Apps:    AppsNS,
		Title:   "Mail Filters",
		ID:      fmt.Sprintf(FeedID, time.Now().UnixNano()),
		Updated: time.Now().UTC().Format(time.RFC3339),
		Author:  config.Author,
		Entries: make([]Entry, 0, len(config.Filters)),
	}

	boolToString := map[bool]string{true: "true", false: "false"}

	for i, filter := range config.Filters {
		filter = normalizeFilter(filter, config.Defaults)
		entry := Entry{
			Category: Category{Term: "filter"},
			Title:    "Mail Filter",
			ID:       fmt.Sprintf("tag:mail.google.com,2008:filter:z%016d", i+1),
			Updated:  time.Now().UTC().Format(time.RFC3339),
			Content:  "",
			Properties: []Property{
				{Name: "from", Value: filter.From},
				{Name: "label", Value: filter.Label},
				{Name: "shouldArchive", Value: boolToString[*filter.ShouldArchive]},
				{Name: "shouldNeverSpam", Value: boolToString[*filter.ShouldNeverSpam]},
				{Name: "shouldNeverMarkAsImportant", Value: boolToString[*filter.ShouldNeverMarkAsImportant]},
			},
		}
		feed.Entries = append(feed.Entries, entry)
	}

	return feed
}

func saveXML(filePath string, feed Feed) error {
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("output file %s already exists", filePath)
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("generating XML: %w", err)
	}

	finalOutput := []byte(XMLHeader + string(output))
	return os.WriteFile(filePath, finalOutput, 0644)
}

func main() {
	var outputFile string
	var verbose bool
	flag.StringVar(&outputFile, "output", "", "output XML file name")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatalf("Usage: %s [-output <xml_file>] [-verbose] <yaml_file>", os.Args[0])
	}

	yamlFile := flag.Args()[0]
	if verbose {
		log.Printf("Reading YAML file: %s", yamlFile)
	}

	config, err := loadConfig(yamlFile)
	checkError(err, "loading configuration")

	if verbose {
		log.Println("Generating XML feed")
	}
	feed := generateFeed(config)

	if outputFile == "" {
		outputFile = strings.TrimSuffix(yamlFile, filepath.Ext(yamlFile)) + ".xml"
	}

	if verbose {
		log.Printf("Saving XML to: %s", outputFile)
	}
	err = saveXML(outputFile, feed)
	checkError(err, "saving XML")

	fmt.Println("XML file successfully generated:", outputFile)
}
