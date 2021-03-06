package gconfig

import (
	"os"
	"path/filepath"
)

const (
	projectName = "omnetpp-offload"
)

func CacheDir() (dir string) {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	dir = filepath.Join(dir, projectName)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	return
}

func ConfigDir() (dir string) {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	dir = filepath.Join(dir, projectName)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	return
}
