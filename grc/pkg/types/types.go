package types

import (
	"encoding/xml"
	"regexp"
)

const (
	AtomNS    = "http://www.w3.org/2005/Atom"
	AppsNS    = "http://www.w3.org/2007/app"
	FeedID    = "urn:uuid:"
	XMLHeader = xml.Header
)

var (
	EmailRegex  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	DomainRegex = regexp.MustCompile(`^[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

type FiltersConfig struct {
	Filters  []Filter  `yaml:"filters"`
	Defaults *Defaults `yaml:"defaults,omitempty"`
}

type Filter struct {
	Name       string     `yaml:"name"`
	Expression string     `yaml:"expression,omitempty"`
	Category   *Category  `yaml:"category,omitempty"`
	Properties []Property `yaml:"properties,omitempty"`
	Author     *Author    `yaml:"author,omitempty"`
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	XMLNS   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`
	ID      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title      string     `xml:"title"`
	ID         string     `xml:"id"`
	Updated    string     `xml:"updated"`
	Summary    string     `xml:"summary"`
	Author     *Author    `xml:"author,omitempty"`
	Categories []Category `xml:"category,omitempty"`
}

type Category struct {
	Term   string `xml:"term,attr" yaml:"term"`
	Scheme string `xml:"scheme,attr,omitempty" yaml:"scheme,omitempty"`
	Label  string `xml:"label,attr,omitempty" yaml:"label,omitempty"`
}

type Property struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Author struct {
	Name  string `xml:"name" yaml:"name"`
	Email string `xml:"email,omitempty" yaml:"email,omitempty"`
	URI   string `xml:"uri,omitempty" yaml:"uri,omitempty"`
}

type Defaults struct {
	Author   *Author   `yaml:"author,omitempty"`
	Category *Category `yaml:"category,omitempty"`
}
