package generate

import (
	"fmt"
	"strings"
	"time"

	"github.com/carlosrabelo/grc/grc/internal/query"
	"github.com/carlosrabelo/grc/grc/pkg/types"
)

// GenerateFeed builds an Atom feed from the configuration without mutating it.
// Filter-level defaults are resolved into the produced entries, leaving the
// input cfg.Filters untouched.
func GenerateFeed(cfg *types.FiltersConfig) (*types.Feed, error) {
	feed := createBaseFeed()
	for i := range cfg.Filters {
		feed.Entries = append(feed.Entries, *createFeedEntry(&cfg.Filters[i], cfg.Defaults))
	}
	return feed, nil
}

func createBaseFeed() *types.Feed {
	now := time.Now().UTC().Format(time.RFC3339)
	return &types.Feed{
		XMLNS:   types.AtomNS,
		Title:   "GRC Feed",
		ID:      types.FeedID + now,
		Updated: now,
	}
}

func createFeedEntry(filter *types.Filter, defaults *types.Defaults) *types.Entry {
	author := filter.Author
	if author == nil && defaults != nil {
		author = defaults.Author
	}

	category := filter.Category
	if category == nil && defaults != nil {
		category = defaults.Category
	}

	entry := &types.Entry{
		Title:   filter.Name,
		ID:      types.FeedID + filter.Name,
		Updated: time.Now().UTC().Format(time.RFC3339),
		Summary: buildFilterSummary(filter),
	}
	if author != nil {
		entry.Author = author
	}
	if category != nil {
		entry.Categories = []types.Category{*category}
	}
	return entry
}

// buildFilterSummary renders a human-readable summary of the filter. Expression
// takes precedence; otherwise we reconstruct a Gmail search string from the
// structured criteria so the entry always carries useful context.
func buildFilterSummary(filter *types.Filter) string {
	var parts []string

	if expr := strings.TrimSpace(filter.Expression); expr != "" {
		parts = append(parts, fmt.Sprintf("Expression: %s", expr))
	} else if filter.Criteria != nil {
		if s := query.BuildSearchString(filter.Criteria); s != "" {
			parts = append(parts, fmt.Sprintf("Search: %s", s))
		}
	}

	for _, p := range filter.Properties {
		parts = append(parts, fmt.Sprintf("%s: %s", p.Name, p.Value))
	}

	return strings.Join(parts, "\n")
}
