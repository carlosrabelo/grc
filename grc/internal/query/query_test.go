package query

import (
	"net/url"
	"testing"

	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func boolPtr(b bool) *bool { return &b }

func TestBuildSearchString_Empty(t *testing.T) {
	if s := BuildSearchString(nil); s != "" {
		t.Errorf("expected empty, got %q", s)
	}
	if s := BuildSearchString(&types.GmailCriteria{}); s != "" {
		t.Errorf("expected empty, got %q", s)
	}
}

func TestBuildSearchString_From(t *testing.T) {
	c := &types.GmailCriteria{From: "newsletter@example.com"}
	if s := BuildSearchString(c); s != "from:newsletter@example.com" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Subject(t *testing.T) {
	c := &types.GmailCriteria{Subject: "weekly report"}
	if s := BuildSearchString(c); s != `subject:"weekly report"` {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_HasTheWord(t *testing.T) {
	c := &types.GmailCriteria{HasTheWord: "important"}
	if s := BuildSearchString(c); s != "important" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_HasAttachment(t *testing.T) {
	c := &types.GmailCriteria{HasAttachment: boolPtr(true)}
	if s := BuildSearchString(c); s != "has:attachment" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_IsUnread(t *testing.T) {
	c := &types.GmailCriteria{Is: "unread"}
	if s := BuildSearchString(c); s != "is:unread" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_InInbox(t *testing.T) {
	c := &types.GmailCriteria{In: "inbox"}
	if s := BuildSearchString(c); s != "in:inbox" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Cc(t *testing.T) {
	c := &types.GmailCriteria{Cc: "team@example.com"}
	if s := BuildSearchString(c); s != "cc:team@example.com" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Bcc(t *testing.T) {
	c := &types.GmailCriteria{Bcc: "hidden@example.com"}
	if s := BuildSearchString(c); s != "bcc:hidden@example.com" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_DeliveredTo(t *testing.T) {
	c := &types.GmailCriteria{DeliveredTo: "me@example.com"}
	if s := BuildSearchString(c); s != "deliveredto:me@example.com" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_List(t *testing.T) {
	c := &types.GmailCriteria{List: "list@example.com"}
	if s := BuildSearchString(c); s != "list:list@example.com" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Filename(t *testing.T) {
	c := &types.GmailCriteria{Filename: "report.pdf"}
	if s := BuildSearchString(c); s != "filename:report.pdf" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Label(t *testing.T) {
	c := &types.GmailCriteria{Label: "Work"}
	if s := BuildSearchString(c); s != "label:Work" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Category(t *testing.T) {
	c := &types.GmailCriteria{Category: "updates"}
	if s := BuildSearchString(c); s != "category:updates" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_After(t *testing.T) {
	c := &types.GmailCriteria{After: "2024/01/01"}
	if s := BuildSearchString(c); s != "after:2024/01/01" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Before(t *testing.T) {
	c := &types.GmailCriteria{Before: "2024/12/31"}
	if s := BuildSearchString(c); s != "before:2024/12/31" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_OlderThan(t *testing.T) {
	c := &types.GmailCriteria{OlderThan: "30d"}
	if s := BuildSearchString(c); s != "older_than:30d" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_NewerThan(t *testing.T) {
	c := &types.GmailCriteria{NewerThan: "7d"}
	if s := BuildSearchString(c); s != "newer_than:7d" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Size(t *testing.T) {
	c := &types.GmailCriteria{Size: "10M"}
	if s := BuildSearchString(c); s != "size:10M" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Larger(t *testing.T) {
	c := &types.GmailCriteria{Larger: "5M"}
	if s := BuildSearchString(c); s != "larger:5M" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_Smaller(t *testing.T) {
	c := &types.GmailCriteria{Smaller: "1M"}
	if s := BuildSearchString(c); s != "smaller:1M" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_MultipleFields(t *testing.T) {
	c := &types.GmailCriteria{
		From:          "boss@example.com",
		Subject:       "urgent",
		HasAttachment: boolPtr(true),
		Is:            "unread",
	}
	s := BuildSearchString(c)
	if !contains(s, "from:boss@example.com") {
		t.Errorf("missing from, got %q", s)
	}
	if !contains(s, "subject:urgent") {
		t.Errorf("missing subject, got %q", s)
	}
	if !contains(s, "has:attachment") {
		t.Errorf("missing has:attachment, got %q", s)
	}
	if !contains(s, "is:unread") {
		t.Errorf("missing is:unread, got %q", s)
	}
}

func TestBuildSearchString_To(t *testing.T) {
	c := &types.GmailCriteria{To: "recipient@example.com"}
	if s := BuildSearchString(c); s != "to:recipient@example.com" {
		t.Errorf("got %q", s)
	}
}

func TestBuildSearchString_HasAttachmentFalse(t *testing.T) {
	c := &types.GmailCriteria{HasAttachment: boolPtr(false)}
	if s := BuildSearchString(c); s != "" {
		t.Errorf("expected empty when HasAttachment=false, got %q", s)
	}
}

func TestBuildSearchString_ValueWithSpaces(t *testing.T) {
	c := &types.GmailCriteria{
		From:    "mail sender",
		Subject: "very important message",
	}
	s := BuildSearchString(c)
	if !contains(s, `from:"mail sender"`) {
		t.Errorf("missing quoted from, got %q", s)
	}
	if !contains(s, `subject:"very important message"`) {
		t.Errorf("missing quoted subject, got %q", s)
	}
}

func TestBuildSearchURL_Empty(t *testing.T) {
	if u := BuildSearchURL(nil); u != "" {
		t.Errorf("expected empty, got %q", u)
	}
	if u := BuildSearchURL(&types.GmailCriteria{}); u != "" {
		t.Errorf("expected empty, got %q", u)
	}
}

func TestBuildSearchURL_Format(t *testing.T) {
	c := &types.GmailCriteria{From: "test@example.com"}
	u := BuildSearchURL(c)
	expectedPrefix := "https://mail.google.com/mail/u/0/#search/"
	if !contains(u, expectedPrefix) {
		t.Errorf("expected prefix %q, got %q", expectedPrefix, u)
	}
}

func TestBuildSearchURL_Encoded(t *testing.T) {
	c := &types.GmailCriteria{Subject: "hello world"}
	u := BuildSearchURL(c)
	expected := "https://mail.google.com/mail/u/0/#search/subject%3A%22hello+world%22"
	if u != expected {
		t.Errorf("got %q, want %q", u, expected)
	}
}

func TestBuildSearchURL_Roundtrip(t *testing.T) {
	c := &types.GmailCriteria{
		From:    "user@example.com",
		Subject: "important",
	}
	u := BuildSearchURL(c)
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Scheme != "https" {
		t.Errorf("scheme: got %q", parsed.Scheme)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestBuildSearchString_AllFields(t *testing.T) {
	c := &types.GmailCriteria{
		From:          "a@b.com",
		To:            "c@d.com",
		Subject:       "subj",
		HasTheWord:    "free",
		HasAttachment: boolPtr(true),
		Is:            "starred",
		In:            "inbox",
		Cc:            "e@f.com",
		Bcc:           "g@h.com",
		DeliveredTo:   "i@j.com",
		List:          "k@l.com",
		Filename:      "doc.pdf",
		Label:         "MyLabel",
		Category:      "social",
		After:         "2024/01/01",
		Before:        "2024/12/31",
		OlderThan:     "30d",
		NewerThan:     "7d",
		Size:          "10M",
		Larger:        "5M",
		Smaller:       "1M",
	}
	s := BuildSearchString(c)
	checks := []string{
		"from:a@b.com",
		"to:c@d.com",
		"subject:subj",
		"free",
		"has:attachment",
		"is:starred",
		"in:inbox",
		"cc:e@f.com",
		"bcc:g@h.com",
		"deliveredto:i@j.com",
		"list:k@l.com",
		"filename:doc.pdf",
		"label:MyLabel",
		"category:social",
		"after:2024/01/01",
		"before:2024/12/31",
		"older_than:30d",
		"newer_than:7d",
		"size:10M",
		"larger:5M",
		"smaller:1M",
	}
	for _, check := range checks {
		if !contains(s, check) {
			t.Errorf("missing %q in %q", check, s)
		}
	}
}
