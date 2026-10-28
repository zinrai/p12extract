package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestKeyFileIsWrittenOwnerOnly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file mode bits do not apply on Windows")
	}
	key := newKey(t)
	leaf := newCertificate(t, "leaf", key)

	dir := t.TempDir()
	files, err := pemFiles(filepath.Join(dir, "client"), key, leaf, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeAllOrNone(files); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(dir, "client.key"))
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("client.key perm = %o, want 600", perm)
	}
}

func TestPemFilesOmitsCAFileWhenThereAreNoCAs(t *testing.T) {
	key := newKey(t)
	leaf := newCertificate(t, "leaf", key)

	files, err := pemFiles("client", key, leaf, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.path, "-ca.crt") {
			t.Errorf("%s returned with no CA certificates", f.path)
		}
	}
}
