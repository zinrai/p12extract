package main

import "os"

// Not left half-written: a stray file would make the next run refuse to
// overwrite it.
func writeAllOrNone(files []pemFile) error {
	var written []string
	for _, f := range files {
		if err := writeNewFile(f); err != nil {
			removeFiles(written)
			return err
		}
		written = append(written, f.path)
	}
	return nil
}

func removeFiles(paths []string) {
	for _, path := range paths {
		os.Remove(path)
	}
}

// Not O_TRUNC: overwriting an existing private key cannot be undone.
func writeNewFile(f pemFile) error {
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, f.perm)
	if err != nil {
		return err
	}
	if _, err := file.Write(f.data); err != nil {
		file.Close()
		os.Remove(f.path)
		return err
	}
	if err := file.Close(); err != nil {
		os.Remove(f.path)
		return err
	}
	return nil
}
