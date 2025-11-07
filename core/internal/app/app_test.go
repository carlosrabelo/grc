package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_BasicExecution(t *testing.T) {
	content := `author:
  name: "Test User"
  email: "test@example.com"
default:
  shouldArchive: false
filters:
  - from: "example@test.com"
    label: "Test"
    shouldArchive: true
`
	tmpFile := createTempYAMLFile(t, content)
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	var stdout, stderr bytes.Buffer
	ctx := context.Background()

	err := Run(ctx, "test-version", "2023-01-01T00:00:00Z", []string{tmpFile}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "XML file successfully generated:") {
		t.Errorf("Expected success message, got: %s", output)
	}

	// Verifica se o arquivo XML foi gerado
	expectedXML := strings.TrimSuffix(tmpFile, filepath.Ext(tmpFile)) + ".xml"
	if _, err := os.Stat(expectedXML); os.IsNotExist(err) {
		t.Errorf("Expected XML file %s to be created", expectedXML)
	}
	defer func() {
		_ = os.Remove(expectedXML)
	}()
}

func TestRun_CustomOutputFile(t *testing.T) {
	content := `author:
  name: "Test User"
  email: "test@example.com"
filters:
  - from: "example@test.com"
    label: "Test"
`
	tmpFile := createTempYAMLFile(t, content)
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	customOutput := filepath.Join(t.TempDir(), "custom.xml")
	var stdout, stderr bytes.Buffer
	ctx := context.Background()

	err := Run(ctx, "test-version", "2023-01-01T00:00:00Z", []string{"-output", customOutput, tmpFile}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, customOutput) {
		t.Errorf("Expected output to contain %s, got: %s", customOutput, output)
	}

	// Verifica se o arquivo XML foi criado no local customizado
	if _, err := os.Stat(customOutput); os.IsNotExist(err) {
		t.Errorf("Expected XML file %s to be created", customOutput)
	}
}

func TestRun_VerboseLogging(t *testing.T) {
	content := `author:
  name: "Test User"
  email: "test@example.com"
filters:
  - from: "example@test.com"
    label: "Test"
`
	tmpFile := createTempYAMLFile(t, content)
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	var stdout, stderr bytes.Buffer
	ctx := context.Background()

	err := Run(ctx, "test-version", "2023-01-01T00:00:00Z", []string{"-verbose", tmpFile}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	logOutput := stderr.String()
	expectedLogs := []string{
		"GRC version test-version",
		"Reading YAML file:",
		"Generating XML feed",
		"Saving XML to:",
	}

	for _, expected := range expectedLogs {
		if !strings.Contains(logOutput, expected) {
			t.Errorf("Expected log to contain %s, got: %s", expected, logOutput)
		}
	}
}

func TestRun_NoArguments(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ctx := context.Background()

	err := Run(ctx, "test-version", "2023-01-01T00:00:00Z", []string{}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Expected error when no arguments provided")
	}

	expectedError := "usage: grc [-output <xml_file>] [-verbose] <yaml_file>"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %s, got: %v", expectedError, err)
	}
}

func TestRun_InvalidYAMLFile(t *testing.T) {
	tmpFile := createTempFile(t, "invalid.yaml", "invalid: yaml: content:")
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	var stdout, stderr bytes.Buffer
	ctx := context.Background()

	err := Run(ctx, "test-version", "2023-01-01T00:00:00Z", []string{tmpFile}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Expected error when invalid YAML file provided")
	}

	if !strings.Contains(err.Error(), "loading configuration:") {
		t.Errorf("Expected configuration loading error, got: %v", err)
	}
}

func TestRun_OutputFileAlreadyExists(t *testing.T) {
	content := `author:
  name: "Test User"
  email: "test@example.com"
filters:
  - from: "example@test.com"
    label: "Test"
`
	tmpFile := createTempYAMLFile(t, content)
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	// Cria o arquivo XML de saída antes da execução
	outputFile := strings.TrimSuffix(tmpFile, filepath.Ext(tmpFile)) + ".xml"
	if err := os.WriteFile(outputFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create existing output file: %v", err)
	}
	defer func() {
		_ = os.Remove(outputFile)
	}()

	var stdout, stderr bytes.Buffer
	ctx := context.Background()

	err := Run(ctx, "test-version", "2023-01-01T00:00:00Z", []string{tmpFile}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Expected error when output file already exists")
	}

	if !strings.Contains(err.Error(), "saving XML:") {
		t.Errorf("Expected XML saving error, got: %v", err)
	}
}

func TestRun_ContextCancelled(t *testing.T) {
	content := `author:
  name: "Test User"
  email: "test@example.com"
filters:
  - from: "example@test.com"
    label: "Test"
`
	tmpFile := createTempYAMLFile(t, content)
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	var stdout, stderr bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancela o contexto imediatamente

	err := Run(ctx, "test-version", "2023-01-01T00:00:00Z", []string{tmpFile}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Expected error when context is cancelled")
	}

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}
}

func TestRun_ForceFlag(t *testing.T) {
	content := `author:
  name: "Test User"
  email: "test@example.com"
filters:
  - from: "example@test.com"
    label: "Test"
`
	tmpFile := createTempYAMLFile(t, content)
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	// Cria o arquivo XML de saída antes da execução
	outputFile := strings.TrimSuffix(tmpFile, filepath.Ext(tmpFile)) + ".xml"
	if err := os.WriteFile(outputFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create existing output file: %v", err)
	}
	defer func() {
		_ = os.Remove(outputFile)
	}()

	var stdout, stderr bytes.Buffer
	ctx := context.Background()

	// Test with --force flag
	err := Run(ctx, "test-version", "2023-01-01T00:00:00Z", []string{"--force", tmpFile}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Expected no error with --force flag, got: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "XML file successfully generated:") {
		t.Errorf("Expected success message, got: %s", output)
	}

	// Verify file was overwritten
	fileContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read overwritten file: %v", err)
	}
	if !strings.Contains(string(fileContent), "Mail Filters") {
		t.Errorf("Expected file to contain 'Mail Filters', got: %s", string(fileContent))
	}
}

func TestRun_ForceFlagVerbose(t *testing.T) {
	content := `author:
  name: "Test User"
  email: "test@example.com"
filters:
  - from: "example@test.com"
    label: "Test"
`
	tmpFile := createTempYAMLFile(t, content)
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	// Cria o arquivo XML de saída antes da execução
	outputFile := strings.TrimSuffix(tmpFile, filepath.Ext(tmpFile)) + ".xml"
	if err := os.WriteFile(outputFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create existing output file: %v", err)
	}
	defer func() {
		_ = os.Remove(outputFile)
	}()

	var stdout, stderr bytes.Buffer
	ctx := context.Background()

	err := Run(ctx, "test-version", "2023-01-01T00:00:00Z", []string{"--force", "--verbose", tmpFile}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Expected no error with --force and --verbose flags, got: %v", err)
	}

	logOutput := stderr.String()
	if !strings.Contains(logOutput, "Overwriting XML to:") {
		t.Errorf("Expected log to contain 'Overwriting XML to:', got: %s", logOutput)
	}
}

// Helper functions

func createTempYAMLFile(t *testing.T, content string) string {
	return createTempFile(t, "test.yaml", content)
}

func createTempFile(t *testing.T, name, content string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, name)
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tmpFile
}