package app

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/carlosrabelo/grc/grc/internal/generate"
	"github.com/carlosrabelo/grc/grc/internal/loader"
	"github.com/carlosrabelo/grc/grc/internal/save"
)

const (
	exitOK    = 0
	exitError = 1
	exitUsage = 2
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

const rootHelp = `grc — Gmail Rule Converter

Usage:
  grc <command> [flags] <input> [output]

Commands:
  atom      Generate an Atom XML feed from a YAML config
  help      Show help (use "grc help <command>" for details)
  version   Show version information

Common flags (available on atom):
  --output PATH   Explicit output path (overrides the positional argument)
  --force         Overwrite an existing output file
  --verbose       Enable verbose logging
  -h, --help      Show help for the current command

Default output paths:
  atom      <input>.xml          (cfg/x.yaml -> cfg/x.xml)

Examples:
  grc atom cfg/example.yaml

Run "grc help <command>" for command-specific help.
`

const atomHelp = `Usage: grc atom [flags] <input.yaml> [output]

Generate an Atom XML feed from a YAML configuration. The feed contains one
entry per filter, with optional author, category and properties.

Arguments:
  input     YAML configuration file (.yaml or .yml). Required.
  output    Output path. Default: <input>.xml.

Flags:
  --output PATH   Explicit output path (overrides the positional output)
  --force         Overwrite an existing output file
  --verbose       Enable verbose logging
  -h, --help      Show this help message

Examples:
  grc atom cfg/example.yaml
  grc atom cfg/example.yaml --output /tmp/feed.xml --force
`

const versionHelp = `Usage: grc version

Print the grc version and build metadata.
`

type options struct {
	output  string
	force   bool
	verbose bool
}

func Run(args []string) int {
	if len(args) == 0 {
		fmt.Print(rootHelp)
		return exitOK
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "help":
		return runHelp(rest)
	case "version", "--version", "-v":
		return runVersion(rest)
	case "-h", "--help":
		fmt.Print(rootHelp)
		return exitOK
	case "atom":
		return runGenerate(rest, "atom", atomHelp)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\n", cmd)
		fmt.Fprint(os.Stderr, rootHelp)
		return exitUsage
	}
}

func runHelp(args []string) int {
	if len(args) == 0 {
		fmt.Print(rootHelp)
		return exitOK
	}
	switch args[0] {
	case "atom":
		fmt.Print(atomHelp)
		return exitOK
	case "version":
		fmt.Print(versionHelp)
		return exitOK
	case "help":
		fmt.Print(rootHelp)
		return exitOK
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\n", args[0])
		fmt.Fprint(os.Stderr, rootHelp)
		return exitUsage
	}
}

func runVersion(args []string) int {
	if hasHelpFlag(args) {
		fmt.Print(versionHelp)
		return exitOK
	}
	fmt.Printf("grc version %s (built %s)\n", Version, BuildTime)
	return exitOK
}

func runGenerate(args []string, mode, help string) int {
	opts, positional, code := parseFlags(mode, args, help)
	if opts == nil {
		return code
	}
	if len(positional) == 0 {
		fmt.Fprintln(os.Stderr, "Error: input file required")
		fmt.Fprint(os.Stderr, help)
		return exitUsage
	}

	input := positional[0]
	output := resolveOutput(input, positionalAt(positional, 1), opts.output, mode)
	logger := createLogger(opts.verbose)

	logger("Loading configuration from %s", input)
	cfg, err := loader.LoadConfig(input)
	if err != nil {
		return fail("Error: %v", err)
	}

	logger("Generating Atom feed with %d filters", len(cfg.Filters))
	feed, ferr := generate.GenerateFeed(cfg)
	if ferr != nil {
		return fail("Error generating feed: %v", ferr)
	}
	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fail("Error marshaling XML: %v", err)
	}
	data = append([]byte(xml.Header), data...)

	logger("Saving to %s", output)
	if _, err := save.SaveFile(data, output, ".xml", opts.force); err != nil {
		return fail("Error: %v", err)
	}

	fmt.Printf("Generated: %s\n", output)
	return exitOK
}

func parseFlags(name string, args []string, help string) (*options, []string, int) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	var opts options
	var showHelp bool
	fs.StringVar(&opts.output, "output", "", "Explicit output path")
	fs.StringVar(&opts.output, "o", "", "Shortcut for --output")
	fs.BoolVar(&opts.force, "force", false, "Overwrite existing output file")
	fs.BoolVar(&opts.verbose, "verbose", false, "Enable verbose logging")
	fs.BoolVar(&showHelp, "help", false, "Show help for this command")
	fs.BoolVar(&showHelp, "h", false, "Shortcut for --help")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprint(os.Stderr, help)
		return nil, nil, exitUsage
	}
	if showHelp {
		fmt.Print(help)
		return nil, nil, exitOK
	}
	return &opts, fs.Args(), exitOK
}

func resolveOutput(input, positionalOutput, flagOutput, mode string) string {
	if flagOutput != "" {
		return flagOutput
	}
	if positionalOutput != "" {
		return positionalOutput
	}
	base := strings.TrimSuffix(input, filepath.Ext(input))
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

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}

func positionalAt(args []string, i int) string {
	if i < 0 || i >= len(args) {
		return ""
	}
	return args[i]
}

func fail(format string, args ...interface{}) int {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	return exitError
}
