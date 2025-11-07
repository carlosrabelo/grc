package rules

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ============================================================================
// Constants
// ============================================================================

const (
	// AtomNS defines the XML namespace for Atom feeds
	AtomNS = "http://www.w3.org/2005/Atom"
	// AppsNS defines the XML namespace for Gmail filters
	AppsNS = "http://schemas.google.com/apps/2006"
	// FeedID template used to build unique ID values
	FeedID = "tag:mail.google.com,2008:filters:%d"
	// XMLHeader contains the standard XML declaration
	XMLHeader = `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
)

// ============================================================================
// Data Types - YAML Configuration
// ============================================================================

// Defaults defines default options that can be reused in filters
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

// Filter represents a Gmail filter coming from the YAML file
type Filter struct {
	// Filtering criteria
	From               string `yaml:"from,omitempty"`
	To                 string `yaml:"to,omitempty"`
	Subject            string `yaml:"subject,omitempty"`
	HasTheWord         string `yaml:"hasTheWord,omitempty"`
	DoesNotHaveTheWord string `yaml:"doesNotHaveTheWord,omitempty"`
	List               string `yaml:"list,omitempty"`
	Query              string `yaml:"query,omitempty"`
	HasAttachment      *bool  `yaml:"hasAttachment,omitempty"`

	// Filter actions
	Label                       string `yaml:"label,omitempty"`
	SmartLabel                  string `yaml:"smartLabel,omitempty"`
	ForwardTo                   string `yaml:"forwardTo,omitempty"`
	ShouldArchive               *bool  `yaml:"shouldArchive,omitempty"`
	ShouldMarkAsRead            *bool  `yaml:"shouldMarkAsRead,omitempty"`
	ShouldStar                  *bool  `yaml:"shouldStar,omitempty"`
	ShouldNeverSpam             *bool  `yaml:"shouldNeverSpam,omitempty"`
	ShouldAlwaysMarkAsImportant *bool  `yaml:"shouldAlwaysMarkAsImportant,omitempty"`
	ShouldNeverMarkAsImportant  *bool  `yaml:"shouldNeverMarkAsImportant,omitempty"`
	ShouldTrash                 *bool  `yaml:"shouldTrash,omitempty"`
}

// Author represents the author block used in Gmail export
type Author struct {
	Name  string `yaml:"name" xml:"name"`
	Email string `yaml:"email" xml:"email"`
}

// FiltersConfig defines how to build the Gmail filters feed
type FiltersConfig struct {
	Author   Author   `yaml:"author"`
	Defaults Defaults `yaml:"default"`
	Filters  []Filter `yaml:"filters"`
}

// ============================================================================
// Data Types - XML Generation
// ============================================================================

// Feed represents the root Atom feed for Gmail filters
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

// Entry models a Gmail filter entry within the feed
type Entry struct {
	XMLName    xml.Name   `xml:"entry"`
	Category   Category   `xml:"category"`
	Title      string     `xml:"title"`
	ID         string     `xml:"id"`
	Updated    string     `xml:"updated"`
	Content    string     `xml:"content,omitempty"`
	Properties []Property `xml:"apps:property"`
}

// Category identifies the filter category
type Category struct {
	Term string `xml:"term,attr"`
}

// Property carries name and value from Gmail filters export
type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// ============================================================================
// Main Public API
// ============================================================================

// LoadConfig reads and validates a YAML configuration file
func LoadConfig(filePath string) (FiltersConfig, error) {
	if err := validateYAMLExtension(filePath); err != nil {
		return FiltersConfig{}, err
	}

	fileContent, err := readFileContent(filePath)
	if err != nil {
		return FiltersConfig{}, err
	}

	config, err := parseYAMLContent(fileContent)
	if err != nil {
		return FiltersConfig{}, err
	}

	if err := validateConfiguration(config); err != nil {
		return FiltersConfig{}, err
	}

	return config, nil
}

// GenerateFeed builds the Atom feed ready for XML serialization
func GenerateFeed(config FiltersConfig, now time.Time) Feed {
	updated := now.Format(time.RFC3339)
	
	feed := createBaseFeed(config.Author, updated, now)
	
	for i, filterConfig := range config.Filters {
		normalizedFilter := applyDefaults(filterConfig, config.Defaults)
		entry := createFeedEntry(normalizedFilter, i, updated)
		feed.Entries = append(feed.Entries, entry)
	}

	return feed
}

// SaveXML writes the feed to disk and refuses to overwrite files unless force is true
func SaveXML(filePath string, feed Feed, force bool) error {
	normalizedPath := ensureXMLExtension(filePath)

	if err := validateFileOverwrite(normalizedPath, force); err != nil {
		return err
	}

	return writeXMLFile(normalizedPath, feed)
}

// ============================================================================
// Reading and Validation Functions
// ============================================================================

// validateYAMLExtension checks if the file has a valid YAML extension
func validateYAMLExtension(filePath string) error {
	lowerPath := strings.ToLower(filePath)
	if !strings.HasSuffix(lowerPath, ".yaml") && !strings.HasSuffix(lowerPath, ".yml") {
		return fmt.Errorf("input file %s must have .yaml or .yml extension", filePath)
	}
	return nil
}

// readFileContent reads the file content
func readFileContent(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return file, nil
}

// parseYAMLContent decodes the YAML content to FiltersConfig structure
func parseYAMLContent(fileContent []byte) (FiltersConfig, error) {
	var config FiltersConfig
	decoder := yaml.NewDecoder(bytes.NewReader(fileContent))
	decoder.KnownFields(true)
	if err := decoder.Decode(&config); err != nil {
		return FiltersConfig{}, fmt.Errorf("decoding YAML: %w", err)
	}
	return config, nil
}

// validateConfiguration validates the complete configuration
func validateConfiguration(config FiltersConfig) error {
	if err := validateAuthorData(config.Author); err != nil {
		return err
	}

	if len(config.Filters) == 0 {
		return fmt.Errorf("at least one filter is required")
	}

	return validateAllFilters(config.Filters, config.Defaults)
}

// validateAuthorData validates the author data
func validateAuthorData(author Author) error {
	if strings.TrimSpace(author.Name) == "" {
		return fmt.Errorf("author name is required")
	}
	if strings.TrimSpace(author.Email) == "" {
		return fmt.Errorf("author email is required")
	}
	return nil
}

// validateAllFilters validates all filters in the configuration
func validateAllFilters(filters []Filter, defaults Defaults) error {
	for i, filter := range filters {
		normalized := applyDefaults(filter, defaults)
		if !hasCriteria(normalized) {
			return fmt.Errorf("filter %d must define at least one condition", i)
		}
		if !hasAction(normalized) {
			return fmt.Errorf("filter %d must define at least one action", i)
		}
	}
	return nil
}

// ============================================================================
// XML Feed Generation Functions
// ============================================================================

// createBaseFeed creates the base feed structure
func createBaseFeed(author Author, updated string, now time.Time) Feed {
	return Feed{
		XMLNS:   AtomNS,
		Apps:    AppsNS,
		Title:   "Mail Filters",
		ID:      fmt.Sprintf(FeedID, now.UnixNano()),
		Updated: updated,
		Author:  author,
		Entries: make([]Entry, 0), // Will be filled dynamically
	}
}

// createFeedEntry creates a filter entry for the feed
func createFeedEntry(filter Filter, index int, updated string) Entry {
	props := buildFilterProperties(filter)

	return Entry{
		Category:   Category{Term: "filter"},
		Title:      "Mail Filter",
		ID:         fmt.Sprintf("tag:mail.google.com,2008:filter:z%016d", index+1),
		Updated:    updated,
		Content:    "",
		Properties: props,
	}
}

// buildFilterProperties builds the filter properties list
func buildFilterProperties(filter Filter) []Property {
	props := make([]Property, 0, 16)
	boolToString := map[bool]string{true: "true", false: "false"}

	addStringProperty := func(name, value string) {
		if value != "" {
			props = append(props, Property{Name: name, Value: value})
		}
	}

	addBoolProperty := func(name string, value *bool) {
		if value != nil {
			props = append(props, Property{Name: name, Value: boolToString[*value]})
		}
	}

	// Filtering criteria
	addStringProperty("from", filter.From)
	addStringProperty("to", filter.To)
	addStringProperty("subject", filter.Subject)
	addStringProperty("hasTheWord", filter.HasTheWord)
	addStringProperty("doesNotHaveTheWord", filter.DoesNotHaveTheWord)
	addStringProperty("list", filter.List)
	addStringProperty("query", filter.Query)
	addBoolProperty("hasAttachment", filter.HasAttachment)

	// Filter actions
	addBoolProperty("shouldArchive", filter.ShouldArchive)
	addBoolProperty("shouldMarkAsRead", filter.ShouldMarkAsRead)
	addBoolProperty("shouldStar", filter.ShouldStar)
	addBoolProperty("shouldNeverSpam", filter.ShouldNeverSpam)
	addBoolProperty("shouldAlwaysMarkAsImportant", filter.ShouldAlwaysMarkAsImportant)
	addBoolProperty("shouldNeverMarkAsImportant", filter.ShouldNeverMarkAsImportant)
	addBoolProperty("shouldTrash", filter.ShouldTrash)

	// String actions
	addStringProperty("label", filter.Label)
	addStringProperty("smartLabelToApply", filter.SmartLabel)
	addStringProperty("forwardTo", filter.ForwardTo)

	return props
}

// ============================================================================
// Persistence and Utility Functions
// ============================================================================

// validateFileOverwrite checks if the file can be overwritten
func validateFileOverwrite(filePath string, force bool) error {
	if _, err := os.Stat(filePath); err == nil {
		if !force {
			return fmt.Errorf("output file %s already exists", filePath)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("checking output file: %w", err)
	}
	return nil
}

// writeXMLFile serializes and writes the XML feed to disk
func writeXMLFile(filePath string, feed Feed) error {
	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("generating XML: %w", err)
	}

	finalOutput := []byte(XMLHeader + string(output))
	return os.WriteFile(filePath, finalOutput, 0o644)
}

// ensureXMLExtension ensures the file has .xml extension
func ensureXMLExtension(filePath string) string {
	if filepath.Ext(filePath) == "" {
		filePath += ".xml"
	}
	return filePath
}

// applyDefaults applies default values to the filter
func applyDefaults(filter Filter, defaults Defaults) Filter {
	applyDefault := func(target **bool, defaultValue bool) {
		if *target == nil && defaultValue {
			*target = boolPtr(true)
		}
	}

	applyDefault(&filter.ShouldArchive, defaults.ShouldArchive)
	applyDefault(&filter.ShouldMarkAsRead, defaults.ShouldMarkAsRead)
	applyDefault(&filter.ShouldStar, defaults.ShouldStar)
	applyDefault(&filter.ShouldNeverSpam, defaults.ShouldNeverSpam)
	applyDefault(&filter.ShouldAlwaysMarkAsImportant, defaults.ShouldAlwaysMarkAsImportant)
	applyDefault(&filter.ShouldNeverMarkAsImportant, defaults.ShouldNeverMarkAsImportant)
	applyDefault(&filter.ShouldTrash, defaults.ShouldTrash)
	applyDefault(&filter.HasAttachment, defaults.HasAttachment)

	return filter
}

// hasCriteria checks if the filter has at least one criterion defined
func hasCriteria(filter Filter) bool {
	// Checks string criteria
	if filter.From != "" || filter.To != "" || filter.Subject != "" ||
		filter.HasTheWord != "" || filter.DoesNotHaveTheWord != "" ||
		filter.List != "" || filter.Query != "" {
		return true
	}

	// Checks boolean criterion
	if filter.HasAttachment != nil && *filter.HasAttachment {
		return true
	}

	return false
}

// hasAction checks if the filter has at least one action defined
func hasAction(filter Filter) bool {
	// Checks string actions
	if filter.Label != "" || filter.SmartLabel != "" || filter.ForwardTo != "" {
		return true
	}

	// Checks boolean actions
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

// boolPtr creates a pointer to boolean
func boolPtr(v bool) *bool {
	return &v
}
