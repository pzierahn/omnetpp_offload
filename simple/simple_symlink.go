package simple

import (
	"os"
	"path/filepath"
)

// SymbolicCopy
//
// This function creates a new directory (target). Afterwards it
// will create symlinks that point to all files in source.
// Think of this as a copy function that creates a symlink instead of a copy.
func SymbolicCopy(source, target string, ignoreDirs map[string]bool) (err error) {

	_ = os.RemoveAll(target)
	err = os.MkdirAll(target, 0755)
	if err != nil {
		return
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if dir, _ := filepath.Split(rel); ignoreDirs[dir] {
			return err
		}

		path, err = filepath.Abs(path)
		if err != nil {
			return err
		}

		mirrorPath := filepath.Join(target, rel)

		if info.IsDir() {
			err = os.MkdirAll(mirrorPath, 0755)
		} else {
			err = os.Symlink(path, mirrorPath)
		}

		return err
	})

	return
}
