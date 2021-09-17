package simple

import (
	"bytes"
	"golang.org/x/crypto/blake2b"
	"io"
	"os"
	"path/filepath"
)

// FilesChangeDetector detects and bundles which files were modified.
type FilesChangeDetector struct {
	Root string
	snap map[string][]byte
}

// Snapshot creates and stores a list of files and their checksums.
func (cfiles *FilesChangeDetector) Snapshot() (err error) {
	cfiles.snap, err = ListDir(cfiles.Root)
	return
}

// ZipChanges compresses the changed files since the last snapshot.
func (cfiles *FilesChangeDetector) ZipChanges(dirname string) (buffer bytes.Buffer, err error) {
	files, err := ListDir(cfiles.Root)
	if err != nil {
		return
	}

	diff := DirDiff(cfiles.snap, files)

	return TarGzFiles(cfiles.Root, dirname, diff)
}

// ListDir lists all files in the given directory with its blake2b checksum.
func ListDir(root string) (files map[string][]byte, err error) {

	files = make(map[string][]byte)

	err = filepath.WalkDir(root, func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() {
			return err
		}

		hash, err := blake2b.New256(nil)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer func() {
			_ = file.Close()
		}()

		if _, err = io.Copy(hash, file); err != nil {
			return err
		}

		files[path] = hash.Sum(nil)

		return nil
	})

	return
}

// DirDiff returns which files changed.
func DirDiff(ori, cha map[string][]byte) (changed map[string]bool) {
	changed = make(map[string]bool)

	for file := range cha {
		changed[file] = !bytes.Equal(ori[file], cha[file])
	}

	return
}
