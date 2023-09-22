package parser

import (
	"testing"
)

func checkStr(t *testing.T, name, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", name, got, want)
	}
}

func TestParse_Subject(t *testing.T) {
	c := Parse("subject:newsletter")
	checkStr(t, "Subject", c.Subject, "newsletter")
}

func TestParse_From(t *testing.T) {
	c := Parse("from:sender@example.com")
	checkStr(t, "From", c.From, "sender@example.com")
}

func TestParse_To(t *testing.T) {
	c := Parse("to:me@example.com")
	checkStr(t, "To", c.To, "me@example.com")
}

func TestParse_HasAttachment(t *testing.T) {
	c := Parse("has:attachment")
	if c.HasAttachment == nil || !*c.HasAttachment {
		t.Error("expected HasAttachment true")
	}
}

func TestParse_Is(t *testing.T) {
	c := Parse("is:unread")
	checkStr(t, "Is", c.Is, "unread")
}

func TestParse_In(t *testing.T) {
	c := Parse("in:trash")
	checkStr(t, "In", c.In, "trash")
}

func TestParse_Cc(t *testing.T) {
	c := Parse("cc:team@example.com")
	checkStr(t, "Cc", c.Cc, "team@example.com")
}

func TestParse_Bcc(t *testing.T) {
	c := Parse("bcc:me@example.com")
	checkStr(t, "Bcc", c.Bcc, "me@example.com")
}

func TestParse_DeliveredTo(t *testing.T) {
	c := Parse("deliveredto:me@example.com")
	checkStr(t, "DeliveredTo", c.DeliveredTo, "me@example.com")
}

func TestParse_List(t *testing.T) {
	c := Parse("list:newsletter@example.com")
	checkStr(t, "List", c.List, "newsletter@example.com")
}

func TestParse_Filename(t *testing.T) {
	c := Parse("filename:report.pdf")
	checkStr(t, "Filename", c.Filename, "report.pdf")
}

func TestParse_After(t *testing.T) {
	c := Parse("after:2024/01/01")
	checkStr(t, "After", c.After, "2024/01/01")
}

func TestParse_Before(t *testing.T) {
	c := Parse("before:2024/06/01")
	checkStr(t, "Before", c.Before, "2024/06/01")
}

func TestParse_OlderThan(t *testing.T) {
	c := Parse("older_than:7d")
	checkStr(t, "OlderThan", c.OlderThan, "7d")
}

func TestParse_NewerThan(t *testing.T) {
	c := Parse("newer_than:2d")
	checkStr(t, "NewerThan", c.NewerThan, "2d")
}

func TestParse_Size(t *testing.T) {
	c := Parse("size:10M")
	checkStr(t, "Size", c.Size, "10M")
}

func TestParse_Larger(t *testing.T) {
	c := Parse("larger:5M")
	checkStr(t, "Larger", c.Larger, "5M")
}

func TestParse_Smaller(t *testing.T) {
	c := Parse("smaller:1M")
	checkStr(t, "Smaller", c.Smaller, "1M")
}

func TestParse_Label(t *testing.T) {
	c := Parse("label:Important")
	checkStr(t, "Label", c.Label, "Important")
}

func TestParse_Category(t *testing.T) {
	c := Parse("category:updates")
	checkStr(t, "Category", c.Category, "updates")
}

func TestParse_MultipleOperators(t *testing.T) {
	c := Parse(`from:sender@example.com subject:newsletter has:attachment is:unread`)
	checkStr(t, "From", c.From, "sender@example.com")
	checkStr(t, "Subject", c.Subject, "newsletter")
	if c.HasAttachment == nil || !*c.HasAttachment {
		t.Error("expected HasAttachment true")
	}
	checkStr(t, "Is", c.Is, "unread")
}

func TestParse_OperatorWithParens(t *testing.T) {
	c := Parse("subject:(newsletter OR weekly)")
	checkStr(t, "Subject", c.Subject, "newsletter OR weekly")
}

func TestParse_OperatorWithParensAndOr(t *testing.T) {
	c := Parse("from:(@example.com OR @test.com)")
	checkStr(t, "From", c.From, "@example.com OR @test.com")
}

func TestParse_QuotedValue(t *testing.T) {
	c := Parse(`subject:"hello world"`)
	checkStr(t, "Subject", c.Subject, "hello world")
}

func TestParse_SingleQuotedValue(t *testing.T) {
	c := Parse("subject:'hello world'")
	checkStr(t, "Subject", c.Subject, "hello world")
}

func TestParse_PlainText(t *testing.T) {
	c := Parse("some random text")
	checkStr(t, "HasTheWord", c.HasTheWord, "some random text")
}

func TestParse_MixedPlainAndOperators(t *testing.T) {
	c := Parse("important subject:invoice from:vendor@example.com")
	checkStr(t, "HasTheWord", c.HasTheWord, "important")
	checkStr(t, "Subject", c.Subject, "invoice")
	checkStr(t, "From", c.From, "vendor@example.com")
}

func TestParse_Negation(t *testing.T) {
	c := Parse("-from:spam@example.com")
	// GmailCriteria has no negation field, so negated tokens must be preserved
	// verbatim in HasTheWord for correct reconstruction by query.BuildSearchString.
	checkStr(t, "From", c.From, "")
	checkStr(t, "HasTheWord", c.HasTheWord, "-from:spam@example.com")
}

func TestParse_Empty(t *testing.T) {
	c := Parse("")
	if c.HasTheWord != "" {
		t.Errorf("expected empty, got %q", c.HasTheWord)
	}
}

func TestParse_ComplexQuery(t *testing.T) {
	c := Parse(`from:(newsletter@example.com OR weekly@example.com) subject:"weekly digest" has:attachment after:2024/01/01`)
	checkStr(t, "From", c.From, "newsletter@example.com OR weekly@example.com")
	checkStr(t, "Subject", c.Subject, "weekly digest")
	if c.HasAttachment == nil || !*c.HasAttachment {
		t.Error("expected HasAttachment true")
	}
	checkStr(t, "After", c.After, "2024/01/01")
}

func TestParse_InPlusIs(t *testing.T) {
	c := Parse("in:anywhere is:starred")
	checkStr(t, "In", c.In, "anywhere")
	checkStr(t, "Is", c.Is, "starred")
}

func TestParse_CcAndBcc(t *testing.T) {
	c := Parse("cc:support@example.com bcc:archive@example.com")
	checkStr(t, "Cc", c.Cc, "support@example.com")
	checkStr(t, "Bcc", c.Bcc, "archive@example.com")
}

func TestParse_FullExample(t *testing.T) {
	c := Parse(`from:(@internal.example.com) subject:(deploy OR release) has:attachment filename:deploy.sh`)
	checkStr(t, "From", c.From, "@internal.example.com")
	checkStr(t, "Subject", c.Subject, "deploy OR release")
	if c.HasAttachment == nil || !*c.HasAttachment {
		t.Error("expected HasAttachment true")
	}
	checkStr(t, "Filename", c.Filename, "deploy.sh")
}

func TestTokenize(t *testing.T) {
	tokens := tokenize(`from:a subject:"hello world" has:attachment`)
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d: %v", len(tokens), tokens)
	}
	if tokens[0] != "from:a" {
		t.Errorf("token 0: got %q", tokens[0])
	}
	if tokens[1] != `subject:"hello world"` {
		t.Errorf("token 1: got %q", tokens[1])
	}
}

func TestTokenizeWithParens(t *testing.T) {
	tokens := tokenize(`from:(a OR b) subject:x`)
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d: %v", len(tokens), tokens)
	}
}
