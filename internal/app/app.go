package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/carlosrabelo/grc/internal/rules"
)

// Run execute the CLI flow in a very simple way.
func Run(ctx context.Context, version, buildTime string, args []string, stdout, stderr io.Writer) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fs := flag.NewFlagSet("grc", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var (
		outputFile string
		verbose    bool
	)
	fs.StringVar(&outputFile, "output", "", "output XML file name")
	fs.BoolVar(&verbose, "verbose", false, "enable verbose logging")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	if len(remaining) == 0 {
		return errors.New("usage: grc [-output <xml_file>] [-verbose] <yaml_file>")
	}

	yamlFile := remaining[0]
	logger := log.New(io.Discard, "", log.LstdFlags)
	if verbose {
		logger.SetOutput(stderr)
		logger.Printf("GRC version %s (build %s)", version, buildTime)
		logger.Printf("Reading YAML file: %s", yamlFile)
	}

	config, err := rules.LoadConfig(yamlFile)
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	if verbose {
		logger.Println("Generating XML feed")
	}
	feed := rules.GenerateFeed(config, time.Now().UTC())

	if outputFile == "" {
		ext := filepath.Ext(yamlFile)
		outputFile = strings.TrimSuffix(yamlFile, ext) + ".xml"
	}

	if verbose {
		logger.Printf("Saving XML to: %s", outputFile)
	}
	if err := rules.SaveXML(outputFile, feed); err != nil {
		return fmt.Errorf("saving XML: %w", err)
	}

	if _, err := fmt.Fprintf(stdout, "XML file successfully generated: %s\n", outputFile); err != nil {
		return fmt.Errorf("writing output message: %w", err)
	}

	return nil
}
