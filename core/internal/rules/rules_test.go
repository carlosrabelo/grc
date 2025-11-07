package rules

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/carlosrabelo/grc/core/internal/testutils"
)

func TestLoadConfig_ValidConfig(t *testing.T) {
	content := `author:
  name: "Test User"
  email: "test@example.com"
default:
  shouldArchive: false
filters:
  - from: "example@test.com"
    label: "Test"
    shouldArchive: true
`
	tmpFile := testutils.CreateTempYAMLFile(t, content)
	defer testutils.CleanupFile(tmpFile)

	config, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.Author.Name != "Test User" {
		t.Errorf("Expected author name 'Test User', got '%s'", config.Author.Name)
	}
	if config.Author.Email != "test@example.com" {
		t.Errorf("Expected author email 'test@example.com', got '%s'", config.Author.Email)
	}
	if len(config.Filters) != 1 {
		t.Errorf("Expected 1 filter, got %d", len(config.Filters))
	}
}

func TestLoadConfig_InvalidExtension(t *testing.T) {
	tmpFile := testutils.CreateTempFile(t, "test.txt", "content")
	defer testutils.CleanupFile(tmpFile)

	_, err := LoadConfig(tmpFile)
	if err == nil || !strings.Contains(err.Error(), "must have .yaml or .yml extension") {
		t.Errorf("Expected extension error, got: %v", err)
	}
}

func TestLoadConfig_MissingAuthor(t *testing.T) {
	content := `filters:
  - from: "example@test.com"
    label: "Test"
`
	tmpFile := testutils.CreateTempYAMLFile(t, content)
	defer testutils.CleanupFile(tmpFile)

	_, err := LoadConfig(tmpFile)
	if err == nil || !strings.Contains(err.Error(), "author name is required") {
		t.Errorf("Expected author validation error, got: %v", err)
	}
}

func TestLoadConfig_EmptyFilters(t *testing.T) {
	content := `author:
  name: "Test User"
  email: "test@example.com"
filters: []
`
	tmpFile := testutils.CreateTempYAMLFile(t, content)
	defer testutils.CleanupFile(tmpFile)

	_, err := LoadConfig(tmpFile)
	if err == nil || !strings.Contains(err.Error(), "at least one filter is required") {
		t.Errorf("Expected empty filters validation error, got: %v", err)
	}
}

func TestGenerateFeed(t *testing.T) {
	config := FiltersConfig{
		Author: Author{
			Name:  "Test User",
			Email: "test@example.com",
		},
		Filters: []Filter{
			{
				From:  "example@test.com",
				Label: "Test",
			},
		},
	}

	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	feed := GenerateFeed(config, now)

	if feed.Title != "Mail Filters" {
		t.Errorf("Expected title 'Mail Filters', got '%s'", feed.Title)
	}
	if feed.Author.Name != "Test User" {
		t.Errorf("Expected author name 'Test User', got '%s'", feed.Author.Name)
	}
	if len(feed.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(feed.Entries))
	}
}

func TestGenerateFeed_DefaultBooleans(t *testing.T) {
	config := FiltersConfig{
		Author: Author{
			Name:  "Test User",
			Email: "test@example.com",
		},
		Defaults: Defaults{
			ShouldArchive:    true,
			ShouldMarkAsRead: false,
			ShouldStar:       true,
		},
		Filters: []Filter{
			{
				From:  "example@test.com",
				Label: "Test",
			},
		},
	}

	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	feed := GenerateFeed(config, now)
	props := feed.Entries[0].Properties

	if !hasProperty(props, "shouldArchive", "true") {
		t.Fatalf("Expected shouldArchive property with value true")
	}
	if !hasProperty(props, "shouldStar", "true") {
		t.Fatalf("Expected shouldStar property with value true")
	}
	if hasProperty(props, "shouldMarkAsRead", "false") {
		t.Fatalf("Did not expect shouldMarkAsRead property when default is false")
	}
}

func TestGenerateFeed_ExplicitFalsePreserved(t *testing.T) {
	config := FiltersConfig{
		Author: Author{
			Name:  "Test User",
			Email: "test@example.com",
		},
		Defaults: Defaults{
			ShouldArchive: true,
		},
		Filters: []Filter{
			{
				From:          "example@test.com",
				Label:         "Test",
				ShouldArchive: testutils.BoolPtr(false),
			},
		},
	}

	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	feed := GenerateFeed(config, now)
	props := feed.Entries[0].Properties

	if !hasProperty(props, "shouldArchive", "false") {
		t.Fatalf("Expected shouldArchive property with value false when explicitly provided")
	}
}

func TestSaveXML_FileExists(t *testing.T) {
	tmpFile := testutils.CreateTempFile(t, "existing.xml", "existing content")
	defer testutils.CleanupFile(tmpFile)

	feed := Feed{Title: "Test"}
	err := SaveXML(tmpFile, feed, false)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected file exists error, got: %v", err)
	}
}

func TestSaveXML_FileExistsWithForce(t *testing.T) {
	tmpFile := testutils.CreateTempFile(t, "existing.xml", "existing content")
	defer testutils.CleanupFile(tmpFile)

	feed := Feed{Title: "Test"}
	err := SaveXML(tmpFile, feed, true)
	if err != nil {
		t.Errorf("Expected no error with force=true, got: %v", err)
	}

	// Verify file was overwritten
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read overwritten file: %v", err)
	}
	if !strings.Contains(string(content), "Test") {
		t.Errorf("Expected file to contain 'Test', got: %s", string(content))
	}
}

func TestHasCriteria(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		expected bool
	}{
		{"Empty filter", Filter{}, false},
		{"Has from", Filter{From: "test@example.com"}, true},
		{"Has subject", Filter{Subject: "Test"}, true},
		{"Has attachment", Filter{HasAttachment: testutils.BoolPtr(true)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasCriteria(tt.filter)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHasAction(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		expected bool
	}{
		{"Empty filter", Filter{}, false},
		{"Has label", Filter{Label: "TestLabel"}, true},
		{"Should archive", Filter{ShouldArchive: testutils.BoolPtr(true)}, true},
		{"Should mark as read", Filter{ShouldMarkAsRead: testutils.BoolPtr(true)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasAction(tt.filter)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGenerateFeed_HasAttachmentDefaults(t *testing.T) {
	config := FiltersConfig{
		Author: Author{
			Name:  "Test User",
			Email: "test@example.com",
		},
		Defaults: Defaults{
			HasAttachment: true,
		},
		Filters: []Filter{
			{
				From:  "example@test.com",
				Label: "Test",
			},
		},
	}

	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	feed := GenerateFeed(config, now)
	props := feed.Entries[0].Properties

	if !hasProperty(props, "hasAttachment", "true") {
		t.Fatalf("Expected hasAttachment property with value true from default")
	}
}

func TestGenerateFeed_HasAttachmentExplicit(t *testing.T) {
	config := FiltersConfig{
		Author: Author{
			Name:  "Test User",
			Email: "test@example.com",
		},
		Defaults: Defaults{
			HasAttachment: true,
		},
		Filters: []Filter{
			{
				From:          "example@test.com",
				Label:         "Test",
				HasAttachment: testutils.BoolPtr(false),
			},
		},
	}

	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	feed := GenerateFeed(config, now)
	props := feed.Entries[0].Properties

	if !hasProperty(props, "hasAttachment", "false") {
		t.Fatalf("Expected hasAttachment property with value false when explicitly set")
	}
}

func TestGenerateFeed_HasAttachmentNotInDefaults(t *testing.T) {
	config := FiltersConfig{
		Author: Author{
			Name:  "Test User",
			Email: "test@example.com",
		},
		Defaults: Defaults{
			ShouldArchive: true,
		},
		Filters: []Filter{
			{
				From:          "example@test.com",
				Label:         "Test",
				HasAttachment: testutils.BoolPtr(true),
			},
		},
	}

	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	feed := GenerateFeed(config, now)
	props := feed.Entries[0].Properties

	if !hasProperty(props, "hasAttachment", "true") {
		t.Fatalf("Expected hasAttachment property with value true when explicitly set")
	}
}

func TestHasCriteria_WithHasAttachment(t *testing.T) {
	tests := []struct {
		name     string
		filter   Filter
		expected bool
	}{
		{"HasAttachment true", Filter{HasAttachment: testutils.BoolPtr(true)}, true},
		{"HasAttachment false", Filter{HasAttachment: testutils.BoolPtr(false)}, false},
		{"HasAttachment nil", Filter{HasAttachment: nil}, false},
		{"HasAttachment true with other criteria", Filter{From: "test@example.com", HasAttachment: testutils.BoolPtr(true)}, true},
		{"HasAttachment false with other criteria", Filter{From: "test@example.com", HasAttachment: testutils.BoolPtr(false)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasCriteria(tt.filter)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Helper functions make tests stay readable.

func hasProperty(props []Property, name, value string) bool {
	for _, p := range props {
		if p.Name == name && p.Value == value {
			return true
		}
	}
	return false
}
