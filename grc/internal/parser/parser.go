package parser

import (
	"strings"
	"unicode"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

type operatorMap struct {
	field    string
	setFn    func(*types.GmailCriteria, string)
}

var operators = []operatorMap{
	{"from", func(c *types.GmailCriteria, v string) { c.From = v }},
	{"to", func(c *types.GmailCriteria, v string) { c.To = v }},
	{"subject", func(c *types.GmailCriteria, v string) { c.Subject = v }},
	{"has", func(c *types.GmailCriteria, v string) {
		if strings.ToLower(v) == "attachment" {
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
	tokens := tokenize(expression)
	for _, token := range tokens {
		applyToken(c, token)
	}
	return c
}

func tokenize(expr string) []string {
	var tokens []string
	var buf strings.Builder
	inQuote := false
	quoteChar := byte(0)
	depth := 0

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if inQuote {
			buf.WriteByte(ch)
			if ch == quoteChar {
				inQuote = false
			}
			continue
		}
		if ch == '"' || ch == '\'' {
			inQuote = true
			quoteChar = ch
			buf.WriteByte(ch)
			continue
		}
		if ch == '(' || ch == ')' {
			if ch == '(' {
				depth++
			} else {
				depth--
			}
			buf.WriteByte(ch)
			continue
		}
		if unicode.IsSpace(rune(ch)) && depth == 0 {
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
			continue
		}
		buf.WriteByte(ch)
	}
	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}
	return tokens
}

func applyToken(c *types.GmailCriteria, token string) {
	// Remove negation prefix
	token = strings.TrimPrefix(token, "-")

	// Check for operator:value pattern
	colonIdx := strings.Index(token, ":")
	if colonIdx == -1 {
		appendHasTheWord(c, token)
		return
	}

	op := strings.ToLower(token[:colonIdx])
	rawValue := token[colonIdx+1:]

	// Unwrap parentheses from value: subject:(a OR b) → a OR b
	value := unwrapParens(rawValue)

	// Unwrap quotes from value
	value = strings.TrimFunc(value, func(r rune) bool {
		return r == '"' || r == '\''
	})

	value = strings.TrimSpace(value)
	if value == "" {
		return
	}

	for _, opEntry := range operators {
		if opEntry.field == op {
			opEntry.setFn(c, value)
			return
		}
	}

	// Unknown operator — treat as hasTheWord with full token
	appendHasTheWord(c, token)
}

func unwrapParens(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' {
		inner := s[1 : len(s)-1]
		depth := 0
		for i := 0; i < len(s); i++ {
			if s[i] == '(' {
				depth++
			} else if s[i] == ')' {
				depth--
			}
			if depth == 0 && i < len(s)-1 {
				return s
			}
		}
		return inner
	}
	return s
}

func appendHasTheWord(c *types.GmailCriteria, token string) {
	token = strings.TrimFunc(token, func(r rune) bool {
		return r == '"' || r == '\''
	})
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
