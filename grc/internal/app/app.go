package app

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/carlosrabelo/grc/grc/internal/generate"
	"github.com/carlosrabelo/grc/grc/internal/loader"
	"github.com/carlosrabelo/grc/grc/internal/save"
)

const helpText = `Usage: grc [options] <input.yaml> [output]

Generate an Atom XML feed from a YAML filter configuration.

Arguments:
  input.yaml            Path to the YAML configuration file
  output                Output path for the XML file (default: <input>.xml)

Options:
  --force               Overwrite existing output file
  --verbose             Enable verbose logging
  --version             Show version information
  --help                Show this help message
`

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, helpText)
		return 1
	}

	fs := flag.NewFlagSet("grc", flag.ContinueOnError)
	fs.Usage = func() { fmt.Fprint(os.Stderr, helpText) }

	showVersion := fs.Bool("version", false, "")
	showHelp := fs.Bool("help", false, "")
	force := fs.Bool("force", false, "")
	verbose := fs.Bool("verbose", false, "")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *showVersion {
		displayVersion()
		return 0
	}
	if *showHelp {
		fmt.Print(helpText)
		return 0
	}
	if fs.NArg() == 0 {
		fmt.Fprint(os.Stderr, helpText)
		return 1
	}

	input := fs.Arg(0)
	output := resolveOutputPath(input, fs.Arg(1))

	logger := createLogger(*verbose)

	logger("Loading configuration from %s", input)
	cfg, err := loader.LoadConfig(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	logger("Generating Atom feed with %d filters", len(cfg.Filters))
	feed, err := generate.GenerateFeed(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating feed: %v\n", err)
		return 1
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling XML: %v\n", err)
		return 1
	}
	data = append([]byte(xml.Header), data...)

	logger("Saving feed to %s", output)
	if err := save.SaveXML(data, output, *force); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	displaySuccessMessage(output)
	return 0
}

func resolveOutputPath(input, output string) string {
	if output != "" {
		return output
	}
	ext := filepath.Ext(input)
	base := strings.TrimSuffix(input, ext)
	return base + ".xml"
}

func createLogger(verbose bool) func(string, ...interface{}) {
	if !verbose {
		return func(string, ...interface{}) {}
	}
	return func(format string, args ...interface{}) {
		log.Printf(format, args...)
	}
}

func displayVersion() {
	fmt.Printf("grc version %s (built %s)\n", Version, BuildTime)
}

func displaySuccessMessage(path string) {
	fmt.Printf("Feed generated: %s\n", path)
}
