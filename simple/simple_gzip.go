package simple

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func TarGz(path, dirname string, exclude ...string) (buffer bytes.Buffer, err error) {

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

		if err != nil || info.IsDir() {
			return
		}

		isSymlink := (info.Mode() & os.ModeSymlink) == os.ModeSymlink

		if !info.Mode().IsRegular() && !isSymlink {
			logger.Printf("skipping '%s', it has unknown file type '%v'\n", walkPath, info.Mode())
			return
		}

		var link string

		if isSymlink {
			link, err = filepath.EvalSymlinks(walkPath)
			if err != nil {
				return
			}

			wDir, _ := filepath.Split(walkPath)
			lDir, linkFile := filepath.Split(link)

			link, err = filepath.Rel(wDir, lDir)
			if err != nil {
				return
			}
			link = filepath.Join(link, linkFile)
		}

		// generate tar header
		header, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = fmt.Errorf("error creating header for %s: %v", walkPath, err)
			return
		}

		//logger.Println("############", info.Name(), header.Linkname)

		relPath, err := filepath.Rel(path, walkPath)
		if err != nil {
			return
		}

		for _, ignore := range exclude {
			if regexp.MustCompile(ignore).MatchString(relPath) {
				logger.Println("exclude", relPath)
				return
			}
		}

		header.Name = filepath.Join(dirname, relPath)

		if err = tw.WriteHeader(header); err != nil {
			err = fmt.Errorf("error writing file %s: %v", walkPath, err)
			return
		}

		if info.Mode().IsRegular() {
			var input *os.File
			input, err = os.Open(walkPath)
			if err != nil {
				err = fmt.Errorf("error open file %s: %v", walkPath, err)
				return
			}

			_, err = io.Copy(tw, input)
			if err != nil {
				err = fmt.Errorf("error copying file %s: %v", walkPath, err)
				return
			}

			err = input.Close()
		}

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

		case tar.TypeSymlink:
			dir, _ := filepath.Split(target)
			err = os.MkdirAll(dir, 0755)

			err = os.Symlink(header.Linkname, target)
			if err != nil {
				break LOOP
			}
		}
	}

	return
}
