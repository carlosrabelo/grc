package query

import (
	"fmt"
	"net/url"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func BuildSearchString(criteria *types.GmailCriteria) string {
	return ""
}

func BuildSearchURL(criteria *types.GmailCriteria) string {
	q := BuildSearchString(criteria)
	if q == "" {
		return ""
	}
	return fmt.Sprintf("https://mail.google.com/mail/u/0/#search/%s", url.QueryEscape(q))
}
