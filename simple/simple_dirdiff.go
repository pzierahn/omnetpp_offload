package simple

import (
	"os"
	"path/filepath"
)

func ListDir(root string) (files map[string]bool, err error) {

	files = make(map[string]bool)

	err = filepath.Walk(root, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		files[path] = true

		return err
	})

	return
}

func DirDiff(ori, cha map[string]bool) (additions map[string]bool) {
	additions = make(map[string]bool)

	for file := range cha {
		if ori[file] {
			//
			// file is not new!
			//

			continue
		}

		//
		// file is new!
		//

		additions[file] = true
	}

	return
}
