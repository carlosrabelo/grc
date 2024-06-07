package parser

import (
	"strings"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

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
