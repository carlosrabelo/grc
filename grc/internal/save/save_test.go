package save

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureExtension(t *testing.T) {
	tests := []struct {
		input string
		ext   string
		want  string
	}{
		{"feed", ".xml", "feed.xml"},
		{"feed.xml", ".xml", "feed.xml"},
		{"feed.XML", ".xml", "feed.XML"},
		{"feed.atom", ".xml", "feed.atom.xml"},
		{"path/to/feed", ".xml", "path/to/feed.xml"},
		{"", ".xml", ".xml"},
		{"filters", ".yaml", "filters.yaml"},
		{"filters.yaml", ".yaml", "filters.yaml"},
		{"filters.xml", "", "filters.xml"},
	}
	for _, tt := range tests {
		t.Run(tt.input+"|"+tt.ext, func(t *testing.T) {
			got := ensureExtension(tt.input, tt.ext)
			if got != tt.want {
				t.Errorf("ensureExtension(%q, %q) = %q, want %q", tt.input, tt.ext, got, tt.want)
			}
		})
	}
}

func TestCheckOverwrite_NewFile(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test-new-"+strings.Repeat("x", 10)+".xml")
	_ = os.Remove(path)
	defer func() { _ = os.Remove(path) }()

	err := checkOverwrite(path, false)
	if err != nil {
		t.Errorf("expected no error for new file, got %v", err)
	}
}

func TestCheckOverwrite_ExistingNoForce(t *testing.T) {
	f, err := os.CreateTemp("", "save-test-*.xml")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	defer func() { _ = os.Remove(f.Name()) }()

	err = checkOverwrite(f.Name(), false)
	if err == nil {
		t.Fatal("expected error for existing file without force")
	}
}

func TestCheckOverwrite_ExistingWithForce(t *testing.T) {
	f, err := os.CreateTemp("", "save-test-*.xml")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	defer func() { _ = os.Remove(f.Name()) }()

	err = checkOverwrite(f.Name(), true)
	if err != nil {
		t.Errorf("expected no error with force, got %v", err)
	}
}

func TestSaveXML_NewFile(t *testing.T) {
	path := filepath.Join(os.TempDir(), "save-new-"+strings.Repeat("y", 10))
	defer func() { _ = os.Remove(path + ".xml") }()

	err := SaveXML([]byte("<feed></feed>"), path, false)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path + ".xml")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "<feed></feed>" {
		t.Errorf("got %q", string(data))
	}
}

func TestSaveXML_OverwriteProtected(t *testing.T) {
	f, err := os.CreateTemp("", "save-protected-*.xml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("<old></old>"); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	defer func() { _ = os.Remove(f.Name()) }()

	err = SaveXML([]byte("<new></new>"), f.Name(), false)
	if err == nil {
		t.Fatal("expected error for overwrite without force")
	}
}

func TestSaveXML_ForceOverwrite(t *testing.T) {
	f, err := os.CreateTemp("", "save-force-*.xml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("<old></old>"); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	defer func() { _ = os.Remove(f.Name()) }()

	err = SaveXML([]byte("<new></new>"), f.Name(), true)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "<new></new>" {
		t.Errorf("got %q", string(data))
	}
}

func TestSaveFile_AddsXMLExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "feed")
	defer func() { _ = os.Remove(path + ".xml") }()

	written, err := SaveFile([]byte("<feed></feed>"), path, ".xml", false)
	if err != nil {
		t.Fatal(err)
	}
	if written != path+".xml" {
		t.Errorf("expected written path %q, got %q", path+".xml", written)
	}

	if _, err := os.Stat(path + ".xml"); os.IsNotExist(err) {
		t.Error("expected file with .xml extension")
	}
}

func TestSaveFile_CreatesDirectories(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "save-nested-"+strings.Repeat("z", 8), "sub")
	defer func() { _ = os.RemoveAll(filepath.Dir(dir)) }()

	path := filepath.Join(dir, "feed.xml")
	_, err := SaveFile([]byte("<feed></feed>"), path, ".xml", false)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to exist")
	}
}

func TestSaveFile_NilData(t *testing.T) {
	path := filepath.Join(os.TempDir(), "save-nil-"+strings.Repeat("w", 8)+".xml")
	defer func() { _ = os.Remove(path) }()

	_, err := SaveFile([]byte{}, path, ".xml", false)
	if err != nil {
		t.Fatal(err)
	}
}
