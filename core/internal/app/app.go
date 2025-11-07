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

	"github.com/carlosrabelo/grc/core/internal/rules"
)

// CLIFlags stores parsed command line flags
type CLIFlags struct {
	outputFile    string
	verbose       bool
	force         bool
	remainingArgs []string
}

// Run executes the main CLI flow
func Run(ctx context.Context, version, buildTime string, args []string, stdout, stderr io.Writer) error {
	if err := checkContextCancellation(ctx); err != nil {
		return err
	}

	flags, err := parseCLIArgs(args)
	if err != nil {
		return err
	}

	if err := validateRequiredArgs(flags); err != nil {
		return err
	}

	yamlFile := flags.remainingArgs[0]
	logger := createLogger(flags.verbose, stderr)

	logApplicationInfo(logger, version, buildTime, yamlFile)

	config, err := loadConfiguration(yamlFile)
	if err != nil {
		return err
	}

	feed, err := generateXMLFeed(config, logger)
	if err != nil {
		return err
	}

	outputFile := resolveOutputPath(yamlFile, flags.outputFile)

	if err := persistXMLFile(logger, outputFile, feed, flags.force); err != nil {
		return err
	}

	return displaySuccessMessage(stdout, outputFile)
}

// ============================================================================
// Configuration and Argument Parsing Functions
// ============================================================================

// parseCLIArgs parses command line flags
func parseCLIArgs(args []string) (*CLIFlags, error) {
	flagSet := flag.NewFlagSet("grc", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	flags := &CLIFlags{}

	flagSet.StringVar(&flags.outputFile, "output", "", "output XML file name")
	flagSet.BoolVar(&flags.verbose, "verbose", false, "enable verbose logging")
	flagSet.BoolVar(&flags.force, "force", false, "overwrite existing XML file")

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	flags.remainingArgs = flagSet.Args()
	return flags, nil
}

// validateRequiredArgs checks if required arguments were provided
func validateRequiredArgs(flags *CLIFlags) error {
	if len(flags.remainingArgs) == 0 {
		return errors.New("usage: grc [-output <xml_file>] [-verbose] [-force] <yaml_file>")
	}
	return nil
}

// createLogger configures the logger based on verbose flag
func createLogger(verbose bool, stderr io.Writer) *log.Logger {
	logger := log.New(io.Discard, "", log.LstdFlags)
	if verbose {
		logger.SetOutput(stderr)
	}
	return logger
}

// checkContextCancellation checks if the context was cancelled
func checkContextCancellation(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// ============================================================================
// Loading and Processing Functions
// ============================================================================

// loadConfiguration loads and validates the YAML configuration file
func loadConfiguration(yamlFile string) (rules.FiltersConfig, error) {
	config, err := rules.LoadConfig(yamlFile)
	if err != nil {
		return rules.FiltersConfig{}, fmt.Errorf("loading configuration: %w", err)
	}
	return config, nil
}

// generateXMLFeed generates the XML feed from configuration
func generateXMLFeed(config rules.FiltersConfig, logger *log.Logger) (rules.Feed, error) {
	logVerboseMessage(logger, "Generating XML feed")
	feed := rules.GenerateFeed(config, time.Now().UTC())
	return feed, nil
}

// resolveOutputPath determines the output file path
func resolveOutputPath(yamlFile, outputFile string) string {
	if outputFile == "" {
		ext := filepath.Ext(yamlFile)
		return strings.TrimSuffix(yamlFile, ext) + ".xml"
	}
	return outputFile
}

// persistXMLFile saves the XML file to disk
func persistXMLFile(logger *log.Logger, outputFile string, feed rules.Feed, force bool) error {
	logFileOperation(logger, outputFile, force)

	if err := rules.SaveXML(outputFile, feed, force); err != nil {
		return fmt.Errorf("saving XML: %w", err)
	}

	return nil
}

// ============================================================================
// Logging and Output Functions
// ============================================================================

// logApplicationInfo displays application information in verbose mode
func logApplicationInfo(logger *log.Logger, version, buildTime, yamlFile string) {
	if logger.Writer() == io.Discard {
		return
	}

	logger.Printf("GRC version %s (build %s)", version, buildTime)
	logger.Printf("Reading YAML file: %s", yamlFile)
}

// logVerboseMessage displays message only in verbose mode
func logVerboseMessage(logger *log.Logger, message string) {
	if logger.Writer() != io.Discard {
		logger.Println(message)
	}
}

// logFileOperation logs file operation in verbose mode
func logFileOperation(logger *log.Logger, outputFile string, force bool) {
	if logger.Writer() == io.Discard {
		return
	}

	if force {
		logger.Printf("Overwriting XML to: %s", outputFile)
	} else {
		logger.Printf("Saving XML to: %s", outputFile)
	}
}

// displaySuccessMessage displays success message on standard output
func displaySuccessMessage(stdout io.Writer, outputFile string) error {
	if _, err := fmt.Fprintf(stdout, "XML file successfully generated: %s\n", outputFile); err != nil {
		return fmt.Errorf("writing output message: %w", err)
	}
	return nil
}
