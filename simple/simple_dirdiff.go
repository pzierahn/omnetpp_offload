package simple

import (
	"bytes"
	"golang.org/x/crypto/blake2b"
	"io"
	"os"
	"path/filepath"
)

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

func DirDiff(ori, cha map[string][]byte) (changed map[string]bool) {
	changed = make(map[string]bool)

	for file := range cha {
		if bytes.Equal(ori[file], cha[file]) {
			//
			// file didn't change!
			//

			continue
		}

		//
		// file changed!
		//

		changed[file] = true
	}

	return
}
