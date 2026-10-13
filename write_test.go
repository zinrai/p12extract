package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAllOrNoneRemovesWrittenFilesWhenALaterOneExists(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "existing")
	if err := os.WriteFile(existing, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	fresh := filepath.Join(dir, "fresh")

	err := writeAllOrNone([]pemFile{
		{fresh, 0o600, []byte("new")},
		{existing, 0o600, []byte("new")},
	})
	if err == nil {
		t.Fatal("writeAllOrNone succeeded although a file exists")
	}
	if _, err := os.Stat(fresh); !os.IsNotExist(err) {
		t.Errorf("%s left behind", fresh)
	}
	if data, _ := os.ReadFile(existing); string(data) != "keep" {
		t.Errorf("existing file overwritten with %q", data)
	}
}
