package parser

import (
	"strings"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

type operatorMap struct {
	field string
	setFn func(*types.GmailCriteria, string)
}

var operators = []operatorMap{
	{"from", func(c *types.GmailCriteria, v string) { c.From = v }},
	{"to", func(c *types.GmailCriteria, v string) { c.To = v }},
	{"subject", func(c *types.GmailCriteria, v string) { c.Subject = v }},
	{"has", func(c *types.GmailCriteria, v string) {
		if strings.EqualFold(v, "attachment") {
			t := true
			c.HasAttachment = &t
		}
	}},
	{"is", func(c *types.GmailCriteria, v string) { c.Is = v }},
	{"in", func(c *types.GmailCriteria, v string) { c.In = v }},
	{"cc", func(c *types.GmailCriteria, v string) { c.Cc = v }},
	{"bcc", func(c *types.GmailCriteria, v string) { c.Bcc = v }},
	{"deliveredto", func(c *types.GmailCriteria, v string) { c.DeliveredTo = v }},
	{"list", func(c *types.GmailCriteria, v string) { c.List = v }},
	{"filename", func(c *types.GmailCriteria, v string) { c.Filename = v }},
	{"after", func(c *types.GmailCriteria, v string) { c.After = v }},
	{"before", func(c *types.GmailCriteria, v string) { c.Before = v }},
	{"older_than", func(c *types.GmailCriteria, v string) { c.OlderThan = v }},
	{"newer_than", func(c *types.GmailCriteria, v string) { c.NewerThan = v }},
	{"size", func(c *types.GmailCriteria, v string) { c.Size = v }},
	{"larger", func(c *types.GmailCriteria, v string) { c.Larger = v }},
	{"smaller", func(c *types.GmailCriteria, v string) { c.Smaller = v }},
	{"label", func(c *types.GmailCriteria, v string) { c.Label = v }},
	{"category", func(c *types.GmailCriteria, v string) { c.Category = v }},
}

func Parse(expression string) *types.GmailCriteria {
	c := &types.GmailCriteria{}
	for _, token := range tokenize(expression) {
		applyToken(c, token)
	}
	return c
}

func tokenize(expr string) []string {
	var tokens []string
	for _, t := range strings.Fields(expr) {
		if t != "" {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

func applyToken(c *types.GmailCriteria, token string) {
	if strings.HasPrefix(token, "-") {
		appendHasTheWord(c, token)
		return
	}

	colonIdx := strings.Index(token, ":")
	if colonIdx == -1 {
		appendHasTheWord(c, token)
		return
	}

	op := strings.ToLower(token[:colonIdx])
	value := strings.TrimSpace(token[colonIdx+1:])
	if value == "" {
		return
	}

	for _, opEntry := range operators {
		if opEntry.field == op {
			opEntry.setFn(c, value)
			return
		}
	}

	appendHasTheWord(c, token)
}

func appendHasTheWord(c *types.GmailCriteria, token string) {
	token = strings.TrimSpace(token)
	if token == "" {
		return
	}
	if c.HasTheWord == "" {
		c.HasTheWord = token
	} else {
		c.HasTheWord = c.HasTheWord + " " + token
	}
}
