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
	"github.com/carlosrabelo/grc/grc/internal/gmail"
	"github.com/carlosrabelo/grc/grc/internal/importer"
	"github.com/carlosrabelo/grc/grc/internal/loader"
	"github.com/carlosrabelo/grc/grc/internal/save"
	"gopkg.in/yaml.v3"
)

const helpText = `Usage: grc [options] <input.yaml> [output]
       grc --import <input.xml> [output]

Generate an Atom XML feed or Gmail filter XML from a YAML configuration,
or import a Gmail filter XML back to YAML.

Arguments:
  input                 YAML config file or Gmail XML file (with --import)
  output                Output path (default: <input>.xml for atom,
                        <input>.gmail.xml for gmail, <input>.yaml for import)

Options:
  --import              Import Gmail filter XML to YAML config
  --format atom|gmail   Output format (default: atom)
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
	format := fs.String("format", "atom", "")
	importMode := fs.Bool("import", false, "")

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
	logger := createLogger(*verbose)

	if *importMode {
		return importPipeline(input, fs.Arg(1), *force, logger)
	}

	output := resolveOutputPath(input, fs.Arg(1), *format)

	logger("Loading configuration from %s", input)
	cfg, err := loader.LoadConfig(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	var data []byte

	switch *format {
	case "atom":
		logger("Generating Atom feed with %d filters", len(cfg.Filters))
		feed, err := generate.GenerateFeed(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating feed: %v\n", err)
			return 1
		}
		data, err = xml.MarshalIndent(feed, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling XML: %v\n", err)
			return 1
		}
		data = append([]byte(xml.Header), data...)

	case "gmail":
		logger("Generating Gmail filter XML with %d filters", len(cfg.Filters))
		data, err = gmail.GenerateGmailXML(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating Gmail XML: %v\n", err)
			return 1
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown format %q (use atom or gmail)\n", *format)
		return 1
	}

	logger("Saving to %s", output)
	if err := save.SaveXML(data, output, *force); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	displaySuccessMessage(output)
	return 0
}

func importPipeline(input, output string, force bool, logger func(string, ...interface{})) int {
	logger("Parsing Gmail XML from %s", input)
	cfg, err := importer.ImportGmailXML(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error importing Gmail XML: %v\n", err)
		return 1
	}

	logger("Imported %d filters", len(cfg.Filters))

	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling YAML: %v\n", err)
		return 1
	}

	outputPath := output
	if outputPath == "" {
		ext := filepath.Ext(input)
		base := strings.TrimSuffix(input, ext)
		outputPath = base + ".yaml"
	}

	logger("Saving to %s", outputPath)
	if _, err := os.Stat(outputPath); err == nil && !force {
		fmt.Fprintf(os.Stderr, "Error: file %q already exists (use --force to overwrite)\n", outputPath)
		return 1
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		return 1
	}
	if err := os.WriteFile(outputPath, yamlData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		return 1
	}

	displaySuccessMessage(outputPath)
	return 0
}

func resolveOutputPath(input, output, format string) string {
	if output != "" {
		return output
	}
	ext := filepath.Ext(input)
	base := strings.TrimSuffix(input, ext)
	suffix := ".xml"
	if format == "gmail" {
		suffix = ".gmail.xml"
	}
	return base + suffix
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
	fmt.Printf("Generated: %s\n", path)
}
