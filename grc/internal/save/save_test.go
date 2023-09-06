package save

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnsureXMLExtension(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"feed", "feed.xml"},
		{"feed.xml", "feed.xml"},
		{"feed.XML", "feed.XML"},
		{"feed.atom", "feed.atom.xml"},
		{"path/to/feed", "path/to/feed.xml"},
		{"", ".xml"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ensureXMLExtension(tt.input)
			if got != tt.want {
				t.Errorf("ensureXMLExtension(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateFileOverwrite_NewFile(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test-new-"+strings.Repeat("x", 10)+".xml")
	os.Remove(path)
	defer os.Remove(path)

	err := validateFileOverwrite(path, false)
	if err != nil {
		t.Errorf("expected no error for new file, got %v", err)
	}
}

func TestValidateFileOverwrite_ExistingNoForce(t *testing.T) {
	f, err := os.CreateTemp("", "save-test-*.xml")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	err = validateFileOverwrite(f.Name(), false)
	if err == nil {
		t.Fatal("expected error for existing file without force")
	}
}

func TestValidateFileOverwrite_ExistingWithForce(t *testing.T) {
	f, err := os.CreateTemp("", "save-test-*.xml")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	err = validateFileOverwrite(f.Name(), true)
	if err != nil {
		t.Errorf("expected no error with force, got %v", err)
	}
}

func TestSaveXML_NewFile(t *testing.T) {
	path := filepath.Join(os.TempDir(), "save-new-"+strings.Repeat("y", 10))
	defer os.Remove(path + ".xml")

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
	f.WriteString("<old></old>")
	f.Close()
	defer os.Remove(f.Name())

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
	f.WriteString("<old></old>")
	f.Close()
	defer os.Remove(f.Name())

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

func TestSaveXML_AddsXMLExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "feed")
	defer os.Remove(path + ".xml")

	err := SaveXML([]byte("<feed></feed>"), path, false)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path + ".xml"); os.IsNotExist(err) {
		t.Error("expected file with .xml extension")
	}
}

func TestSaveXML_CreatesDirectories(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "save-nested-"+strings.Repeat("z", 8), "sub")
	defer os.RemoveAll(filepath.Dir(dir))

	path := filepath.Join(dir, "feed.xml")
	err := SaveXML([]byte("<feed></feed>"), path, false)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to exist")
	}
}

func TestSaveXML_NilData(t *testing.T) {
	path := filepath.Join(os.TempDir(), "save-nil-"+strings.Repeat("w", 8)+".xml")
	defer os.Remove(path)

	err := SaveXML([]byte{}, path, false)
	if err != nil {
		t.Fatal(err)
	}
}
