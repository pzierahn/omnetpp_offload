package simple

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func TarGz(path, dirname string) (buffer bytes.Buffer, err error) {

	zr := gzip.NewWriter(&buffer)
	defer func() {
		_ = zr.Close()
	}()

	tw := tar.NewWriter(zr)
	defer func() {
		_ = tw.Close()
	}()

	walker := func(walkPath string, info os.FileInfo, inErr error) (err error) {

		err = inErr

		if err != nil ||
			strings.HasPrefix(info.Name(), ".") ||
			info.IsDir() {
			return
		}

		var input *os.File
		input, err = os.Open(walkPath)
		if err != nil {
			return
		}

		// generate tar header
		header := &tar.Header{
			Name:    dirname + "/" + strings.TrimPrefix(walkPath, path),
			Size:    info.Size(),
			Mode:    int64(info.Mode()),
			ModTime: info.ModTime(),
		}

		if err = tw.WriteHeader(header); err != nil {
			return
		}

		_, err = io.Copy(tw, input)
		if err != nil {
			return
		}

		err = input.Close()

		return
	}

	err = filepath.Walk(path, walker)

	return
}
