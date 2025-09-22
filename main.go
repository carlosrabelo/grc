package main

import (
	"bytes"
	"encoding/xml"
	"errors"
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
	ShouldArchive               bool `yaml:"shouldArchive"`
	ShouldMarkAsRead            bool `yaml:"shouldMarkAsRead"`
	ShouldStar                  bool `yaml:"shouldStar"`
	ShouldNeverSpam             bool `yaml:"shouldNeverSpam"`
	ShouldAlwaysMarkAsImportant bool `yaml:"shouldAlwaysMarkAsImportant"`
	ShouldNeverMarkAsImportant  bool `yaml:"shouldNeverMarkAsImportant"`
	ShouldTrash                 bool `yaml:"shouldTrash"`
	HasAttachment               bool `yaml:"hasAttachment"`
}

type Filter struct {
	From                        string `yaml:"from,omitempty"`
	To                          string `yaml:"to,omitempty"`
	Subject                     string `yaml:"subject,omitempty"`
	HasTheWord                  string `yaml:"hasTheWord,omitempty"`
	DoesNotHaveTheWord          string `yaml:"doesNotHaveTheWord,omitempty"`
	List                        string `yaml:"list,omitempty"`
	Query                       string `yaml:"query,omitempty"`
	Label                       string `yaml:"label,omitempty"`
	SmartLabel                  string `yaml:"smartLabel,omitempty"`
	ForwardTo                   string `yaml:"forwardTo,omitempty"`
	HasAttachment               *bool  `yaml:"hasAttachment,omitempty"`
	ShouldArchive               *bool  `yaml:"shouldArchive,omitempty"`
	ShouldMarkAsRead            *bool  `yaml:"shouldMarkAsRead,omitempty"`
	ShouldStar                  *bool  `yaml:"shouldStar,omitempty"`
	ShouldNeverSpam             *bool  `yaml:"shouldNeverSpam,omitempty"`
	ShouldAlwaysMarkAsImportant *bool  `yaml:"shouldAlwaysMarkAsImportant,omitempty"`
	ShouldNeverMarkAsImportant  *bool  `yaml:"shouldNeverMarkAsImportant,omitempty"`
	ShouldTrash                 *bool  `yaml:"shouldTrash,omitempty"`
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
	Content    string     `xml:"content,omitempty"`
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

func boolPtr(v bool) *bool {
	return &v
}

func hasCriteria(filter Filter) bool {
	if filter.From != "" || filter.To != "" || filter.Subject != "" || filter.HasTheWord != "" || filter.DoesNotHaveTheWord != "" || filter.List != "" || filter.Query != "" {
		return true
	}
	if filter.HasAttachment != nil && *filter.HasAttachment {
		return true
	}
	return false
}

func hasAction(filter Filter) bool {
	if filter.Label != "" || filter.SmartLabel != "" || filter.ForwardTo != "" {
		return true
	}
	boolActions := []*bool{
		filter.ShouldArchive,
		filter.ShouldMarkAsRead,
		filter.ShouldStar,
		filter.ShouldNeverSpam,
		filter.ShouldAlwaysMarkAsImportant,
		filter.ShouldNeverMarkAsImportant,
		filter.ShouldTrash,
	}
	for _, action := range boolActions {
		if action != nil && *action {
			return true
		}
	}
	return false
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
	decoder := yaml.NewDecoder(bytes.NewReader(file))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		return FiltersConfig{}, fmt.Errorf("decoding YAML: %w", err)
	}

	if err := validateConfig(config); err != nil {
		return FiltersConfig{}, err
	}

	return config, nil
}

func validateConfig(config FiltersConfig) error {
	if strings.TrimSpace(config.Author.Name) == "" {
		return fmt.Errorf("author name is required")
	}
	if strings.TrimSpace(config.Author.Email) == "" {
		return fmt.Errorf("author email is required")
	}

	for i, filter := range config.Filters {
		normalized := normalizeFilter(filter, config.Defaults)
		if !hasCriteria(normalized) {
			return fmt.Errorf("filter %d must define at least one condition", i)
		}
		if !hasAction(normalized) {
			return fmt.Errorf("filter %d must define at least one action", i)
		}
	}

	return nil
}

func normalizeFilter(filter Filter, defaults Defaults) Filter {
	if filter.ShouldArchive == nil {
		filter.ShouldArchive = boolPtr(defaults.ShouldArchive)
	}
	if filter.ShouldMarkAsRead == nil {
		filter.ShouldMarkAsRead = boolPtr(defaults.ShouldMarkAsRead)
	}
	if filter.ShouldStar == nil {
		filter.ShouldStar = boolPtr(defaults.ShouldStar)
	}
	if filter.ShouldNeverSpam == nil {
		filter.ShouldNeverSpam = boolPtr(defaults.ShouldNeverSpam)
	}
	if filter.ShouldAlwaysMarkAsImportant == nil {
		filter.ShouldAlwaysMarkAsImportant = boolPtr(defaults.ShouldAlwaysMarkAsImportant)
	}
	if filter.ShouldNeverMarkAsImportant == nil {
		filter.ShouldNeverMarkAsImportant = boolPtr(defaults.ShouldNeverMarkAsImportant)
	}
	if filter.ShouldTrash == nil {
		filter.ShouldTrash = boolPtr(defaults.ShouldTrash)
	}
	if filter.HasAttachment == nil && defaults.HasAttachment {
		filter.HasAttachment = boolPtr(true)
	}
	return filter
}


func generateFeed(config FiltersConfig, now time.Time) Feed {
	updated := now.Format(time.RFC3339)
	feed := Feed{
		XMLNS:   AtomNS,
		Apps:    AppsNS,
		Title:   "Mail Filters",
		ID:      fmt.Sprintf(FeedID, now.UnixNano()),
		Updated: updated,
		Author:  config.Author,
		Entries: make([]Entry, 0, len(config.Filters)),
	}

	boolToString := map[bool]string{true: "true", false: "false"}

	for i, f := range config.Filters {
		filter := normalizeFilter(f, config.Defaults)
		props := make([]Property, 0, 16)
		addString := func(name, value string) {
			if value != "" {
				props = append(props, Property{Name: name, Value: value})
			}
		}
		addBool := func(name string, value *bool) {
			if value != nil {
				props = append(props, Property{Name: name, Value: boolToString[*value]})
			}
		}

		addString("from", filter.From)
		addString("to", filter.To)
		addString("subject", filter.Subject)
		addString("hasTheWord", filter.HasTheWord)
		addString("doesNotHaveTheWord", filter.DoesNotHaveTheWord)
		addString("list", filter.List)
		addString("query", filter.Query)
		addBool("hasAttachment", filter.HasAttachment)

		addBool("shouldArchive", filter.ShouldArchive)
		addBool("shouldMarkAsRead", filter.ShouldMarkAsRead)
		addBool("shouldStar", filter.ShouldStar)
		addBool("shouldNeverSpam", filter.ShouldNeverSpam)
		addBool("shouldAlwaysMarkAsImportant", filter.ShouldAlwaysMarkAsImportant)
		addBool("shouldNeverMarkAsImportant", filter.ShouldNeverMarkAsImportant)
		addBool("shouldTrash", filter.ShouldTrash)

		addString("label", filter.Label)
		addString("smartLabelToApply", filter.SmartLabel)
		addString("forwardTo", filter.ForwardTo)

		entry := Entry{
			Category:   Category{Term: "filter"},
			Title:      "Mail Filter",
			ID:         fmt.Sprintf("tag:mail.google.com,2008:filter:z%016d", i+1),
			Updated:    updated,
			Content:    "",
			Properties: props,
		}
		feed.Entries = append(feed.Entries, entry)
	}

	return feed
}

func saveXML(filePath string, feed Feed) error {
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("output file %s already exists", filePath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("checking output file: %w", err)
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
	feed := generateFeed(config, time.Now().UTC())

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
