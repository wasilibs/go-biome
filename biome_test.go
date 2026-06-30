package biome

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/wasilibs/go-biome/v2/internal/runner"
)

//go:embed testdata/in
var inFiles embed.FS

//go:embed testdata/exp
var expFiles embed.FS

// TestFormat formats and compares with golden data.
// Run with UPDATE_GOLDEN=1 to regenerate testdata/exp.
func TestFormat(t *testing.T) {
	inFS, err := fs.Sub(inFiles, "testdata/in")
	if err != nil {
		t.Fatal(err)
	}
	expFS, err := fs.Sub(expFiles, "testdata/exp")
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	if err := fs.WalkDir(inFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("biome: reading testdata: %w", err)
		}
		if d.IsDir() {
			return nil
		}
		c, _ := fs.ReadFile(inFS, path)
		return os.WriteFile(filepath.Join(dir, path), c, 0o644)
	}); err != nil {
		t.Fatal(err)
	}

	stdin := bytes.Buffer{}
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	ret := runner.Run("biome", []string{"format", "--write", "."}, &stdin, &stdout, &stderr, dir)
	if want := 0; ret != want {
		t.Fatalf("unexpected return code: have %d, want %d\nstdout:\n%s\nstderr:\n%s", ret, want, stdout.String(), stderr.String())
	}

	update := os.Getenv("UPDATE_GOLDEN") == "1"
	err = fs.WalkDir(inFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		got, err := os.ReadFile(filepath.Join(dir, path))
		if err != nil {
			return fmt.Errorf("biome: reading formatted file: %w", err)
		}
		if update {
			if err := os.WriteFile(filepath.Join("testdata/exp", path), got, 0o644); err != nil {
				return fmt.Errorf("biome: writing golden: %w", err)
			}
			return nil
		}
		want, err := fs.ReadFile(expFS, path)
		if err != nil {
			return fmt.Errorf("missing golden for %s: %w", path, err)
		}
		if string(got) != string(want) {
			t.Errorf("%s mismatch:\n--- got ---\n%s\n--- want ---\n%s", path, got, want)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if update {
		t.Log("updated testdata/exp")
	}
}
