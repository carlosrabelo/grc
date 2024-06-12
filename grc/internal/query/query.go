package query

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func BuildSearchString(criteria *types.GmailCriteria) string {
	if criteria == nil {
		return ""
	}
	var parts []string

	add := func(operator, value string) {
		if value == "" {
			return
		}
		if strings.ContainsAny(value, " ") {
			value = `"` + value + `"`
		}
		if operator == "" {
			parts = append(parts, value)
		} else {
			parts = append(parts, operator+":"+value)
		}
	}

	add("from", criteria.From)
	add("to", criteria.To)
	add("subject", criteria.Subject)
	add("", criteria.HasTheWord)
	add("is", criteria.Is)
	add("in", criteria.In)
	add("cc", criteria.Cc)
	add("bcc", criteria.Bcc)
	add("deliveredto", criteria.DeliveredTo)
	add("list", criteria.List)
	add("filename", criteria.Filename)
	add("label", criteria.Label)
	add("category", criteria.Category)
	add("after", criteria.After)
	add("before", criteria.Before)
	add("older_than", criteria.OlderThan)
	add("newer_than", criteria.NewerThan)
	add("size", criteria.Size)
	add("larger", criteria.Larger)
	add("smaller", criteria.Smaller)

	if criteria.HasAttachment != nil && *criteria.HasAttachment {
		parts = append(parts, "has:attachment")
	}

	return strings.Join(parts, " ")
}

func BuildSearchURL(criteria *types.GmailCriteria) string {
	q := BuildSearchString(criteria)
	if q == "" {
		return ""
	}
	return fmt.Sprintf("https://mail.google.com/mail/u/0/#search/%s", url.QueryEscape(q))
}
