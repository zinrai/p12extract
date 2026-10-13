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

	files, err := pemFiles(filepath.Join(t.TempDir(), "client"), key, leaf, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeAllOrNone(files); err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.path, ".key") {
			info, err := os.Stat(f.path)
			if err != nil {
				t.Fatal(err)
			}
			if perm := info.Mode().Perm(); perm != 0o600 {
				t.Errorf("%s perm = %o, want 600", f.path, perm)
			}
			return
		}
	}
	t.Fatal("no key file")
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
