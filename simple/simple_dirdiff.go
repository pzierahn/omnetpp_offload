package simple

import (
	"encoding/hex"
	"golang.org/x/crypto/blake2b"
	"io"
	"os"
	"path/filepath"
)

func ListDir(root string) (files map[string]string, err error) {

	files = make(map[string]string)

	err = filepath.WalkDir(root, func(path string, dir os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() {
			return err
		}

		hasher, err := blake2b.New256(nil)
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

		if _, err = io.Copy(hasher, file); err != nil {
			return err
		}

		hash := hasher.Sum(nil)
		files[path] = hex.EncodeToString(hash[:])

		return nil
	})

	return
}

func DirDiff(ori, cha map[string]string) (additions map[string]bool) {
	additions = make(map[string]bool)

	for file := range cha {
		if ori[file] == cha[file] {
			//
			// file didn't change!
			//

			continue
		}

		//
		// file is changed!
		//

		additions[file] = true
	}

	return
}
