package parser

import (
	"strings"
	"unicode"
	"unicode/utf8"

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

// Parse converts a Gmail search expression into structured GmailCriteria.
// Negated tokens (e.g. "-from:spam") are preserved verbatim in HasTheWord,
// since the GmailCriteria struct has no dedicated negation field.
func Parse(expression string) *types.GmailCriteria {
	c := &types.GmailCriteria{}
	for _, token := range tokenize(expression) {
		applyToken(c, token)
	}
	return c
}

func tokenize(expr string) []string {
	var tokens []string
	var buf strings.Builder
	inQuote := false
	quoteRune := rune(0)
	depth := 0

	for i := 0; i < len(expr); {
		r, size := utf8.DecodeRuneInString(expr[i:])
		i += size

		if inQuote {
			buf.WriteRune(r)
			if r == quoteRune {
				inQuote = false
			}
			continue
		}
		if r == '"' || r == '\'' {
			inQuote = true
			quoteRune = r
			buf.WriteRune(r)
			continue
		}
		switch r {
		case '(':
			depth++
			buf.WriteRune(r)
			continue
		case ')':
			if depth > 0 {
				depth--
			}
			buf.WriteRune(r)
			continue
		}
		if unicode.IsSpace(r) && depth == 0 {
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
			continue
		}
		buf.WriteRune(r)
	}
	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}
	return tokens
}

func applyToken(c *types.GmailCriteria, token string) {
	// Detect but preserve negation: "-from:spam" stays intact so it can be
	// reconstructed correctly via query.BuildSearchString.
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
	rawValue := token[colonIdx+1:]

	value := strings.TrimSpace(unwrapQuotes(unwrapParens(rawValue)))
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

// unwrapParens removes a single pair of outer parentheses only when they
// span the entire string (e.g. "(a OR b)" → "a OR b", but "(a) OR (b)" is unchanged).
func unwrapParens(s string) string {
	s = strings.TrimSpace(s)
	runes := []rune(s)
	if len(runes) < 2 || runes[0] != '(' || runes[len(runes)-1] != ')' {
		return s
	}
	depth := 0
	for i, r := range runes {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		}
		if depth == 0 && i < len(runes)-1 {
			return s
		}
	}
	return string(runes[1 : len(runes)-1])
}

// unwrapQuotes removes a matching pair of surrounding quotes (single or double).
func unwrapQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	first := s[0]
	last := s[len(s)-1]
	if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
		return s[1 : len(s)-1]
	}
	return s
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
