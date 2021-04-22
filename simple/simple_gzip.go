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

func UnTarGz(dst string, buffer io.Reader) (err error) {

	zip, err := gzip.NewReader(buffer)
	if err != nil {
		return
	}

	defer func() {
		_ = zip.Close()
	}()

	tarReader := tar.NewReader(zip)

LOOP:
	for {
		var header *tar.Header
		header, err = tarReader.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			err = nil
			break LOOP

		// return any other error
		case err != nil:
			break LOOP

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue LOOP
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err = os.Stat(target); err != nil {
				if err = os.MkdirAll(target, 0755); err != nil {
					break LOOP
				}
			}

		// if it's a file create it
		case tar.TypeReg:

			dir, _ := filepath.Split(target)
			err = os.MkdirAll(dir, 0755)

			var file *os.File
			file, err = os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				break LOOP
			}

			// copy over contents
			if _, err = io.Copy(file, tarReader); err != nil {
				break LOOP
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			if err = file.Close(); err != nil {
				break LOOP
			}
		}
	}

	return
}
