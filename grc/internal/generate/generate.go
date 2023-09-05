package generate

import (
	"fmt"
	"strings"
	"time"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func GenerateFeed(cfg *types.FiltersConfig) (*types.Feed, error) {
	feed := createBaseFeed(cfg)
	for i := range cfg.Filters {
		filter := &cfg.Filters[i]
		applyDefaults(filter, cfg.Defaults)
		entry := createFeedEntry(filter)
		feed.Entries = append(feed.Entries, *entry)
	}
	return feed, nil
}

func createBaseFeed(cfg *types.FiltersConfig) *types.Feed {
	now := time.Now().UTC().Format(time.RFC3339)
	return &types.Feed{
		XMLNS:   types.AtomNS,
		Title:   "GRC Feed",
		ID:      types.FeedID + now,
		Updated: now,
	}
}

func createFeedEntry(filter *types.Filter) *types.Entry {
	summary := buildFilterProperties(filter)
	entry := &types.Entry{
		Title:   filter.Name,
		ID:      types.FeedID + filter.Name,
		Updated: time.Now().UTC().Format(time.RFC3339),
		Summary: summary,
	}
	if filter.Author != nil {
		entry.Author = filter.Author
	}
	if filter.Category != nil {
		entry.Categories = []types.Category{*filter.Category}
	}
	return entry
}

func buildFilterProperties(filter *types.Filter) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Expression: %s", filter.Expression))
	for _, p := range filter.Properties {
		parts = append(parts, fmt.Sprintf("%s: %s", p.Name, p.Value))
	}
	return strings.Join(parts, "\n")
}

func applyDefaults(filter *types.Filter, defaults *types.Defaults) {
	if defaults == nil {
		return
	}
	if filter.Author == nil && defaults.Author != nil {
		filter.Author = defaults.Author
	}
	if filter.Category == nil && defaults.Category != nil {
		filter.Category = defaults.Category
	}
}
