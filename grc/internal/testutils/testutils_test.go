package testutils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateTempFile(t *testing.T) {
	path, err := CreateTempFile("", "test-*.txt", "hello world")
	if err != nil {
		t.Fatal(err)
	}
	defer CleanupFile(path)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello world" {
		t.Errorf("got %q, want %q", string(data), "hello world")
	}
}

func TestCreateTempFileInDir(t *testing.T) {
	dir := t.TempDir()
	path, err := CreateTempFile(dir, "test-*.txt", "content")
	if err != nil {
		t.Fatal(err)
	}
	defer CleanupFile(path)

	if !strings.HasPrefix(path, dir) {
		t.Errorf("expected path in %s, got %s", dir, path)
	}
}

func TestCreateTempYAMLFile(t *testing.T) {
	path, err := CreateTempYAMLFile("", "key: value\n")
	if err != nil {
		t.Fatal(err)
	}
	defer CleanupFile(path)

	if !strings.HasSuffix(path, ".yaml") {
		t.Errorf("expected .yaml extension, got %s", path)
	}
}

func TestCreateTempYAMLFileInDir(t *testing.T) {
	dir := t.TempDir()
	path, err := CreateTempYAMLFile(dir, "foo: bar\n")
	if err != nil {
		t.Fatal(err)
	}
	defer CleanupFile(path)

	if !strings.HasPrefix(path, dir) {
		t.Errorf("expected path in %s, got %s", dir, path)
	}
}

func TestCleanupFile(t *testing.T) {
	f, err := os.CreateTemp("", "cleanup-*")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	CleanupFile(f.Name())
	if _, err := os.Stat(f.Name()); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestCleanupFileEmpty(t *testing.T) {
	CleanupFile("")
}

func TestCleanupFileNotExist(t *testing.T) {
	CleanupFile(filepath.Join(os.TempDir(), "nonexistent-"+strings.Repeat("x", 10)))
}

func TestBoolPtr(t *testing.T) {
	b := BoolPtr(true)
	if b == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *b != true {
		t.Errorf("got %v, want true", *b)
	}

	b = BoolPtr(false)
	if b == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *b != false {
		t.Errorf("got %v, want false", *b)
	}
}

func TestBoolPtrUnique(t *testing.T) {
	b1 := BoolPtr(true)
	b2 := BoolPtr(true)
	if b1 == b2 {
		t.Error("expected different pointers")
	}
}
