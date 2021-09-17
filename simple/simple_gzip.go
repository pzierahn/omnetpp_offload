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

// TarGzFiles tars and compresses files.
func TarGzFiles(path, dirname string, files map[string]bool) (buffer bytes.Buffer, err error) {
	zr, _ := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
	defer func() {
		_ = zr.Close()
	}()

	tw := tar.NewWriter(zr)
	defer func() {
		_ = tw.Close()
	}()

	for walkPath, include := range files {

		if !include {
			continue
		}

		info, err := os.Stat(walkPath)
		if err != nil || info.IsDir() {
			continue
		}

		isSymlink := (info.Mode() & os.ModeSymlink) == os.ModeSymlink

		if !info.Mode().IsRegular() && !isSymlink {
			logger.Printf("skipping '%s', it has unknown file type '%v'\n", walkPath, info.Mode())
			continue
		}

		var link string

		if isSymlink {
			link, err = filepath.EvalSymlinks(walkPath)
			if err != nil {
				break
			}

			wDir, _ := filepath.Split(walkPath)
			lDir, linkFile := filepath.Split(link)

			link, err = filepath.Rel(wDir, lDir)
			if err != nil {
				break
			}
			link = filepath.Join(link, linkFile)
		}

		// generate tar header
		header, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = fmt.Errorf("error creating header for %s: %v", walkPath, err)
			break
		}

		relPath, err := filepath.Rel(path, walkPath)
		if err != nil {
			break
		}

		header.Name = filepath.Join(dirname, relPath)

		if err = tw.WriteHeader(header); err != nil {
			err = fmt.Errorf("error writing file %s: %v", walkPath, err)
			break
		}

		if info.Mode().IsRegular() {
			var input *os.File
			input, err = os.Open(walkPath)
			if err != nil {
				err = fmt.Errorf("error open file %s: %v", walkPath, err)
				break
			}

			_, err = io.Copy(tw, input)
			if err != nil {
				err = fmt.Errorf("error copying file %s: %v", walkPath, err)
				break
			}

			err = input.Close()
		}
	}

	return
}

// TarGz tars and compresses all files in a directory.
// To exclude files you can define regex.
func TarGz(path, dirname string, exclude ...string) (buffer bytes.Buffer, err error) {

	files := make(map[string]bool)

	walker := func(walkPath string, _ os.DirEntry, inErr error) (err error) {

		if inErr != nil {
			return inErr
		}

		relPath, err := filepath.Rel(path, walkPath)
		if err != nil {
			return
		}

		for _, ignore := range exclude {
			if regexp.MustCompile(ignore).MatchString(relPath) {
				logger.Printf("exclude pattern='%s' file=%s", ignore, relPath)
				return
			}
		}

		files[walkPath] = true

		return
	}

	err = filepath.WalkDir(path, walker)

	buffer, err = TarGzFiles(path, dirname, files)

	return
}

// ExtractTarGz extracts a tar gzip archive to the desired destination.
// Source: https://gist.github.com/mislav/ca62231f776526729b5d4ddd74ad6657
func ExtractTarGz(dst string, byt []byte) (err error) {

	buffer := bytes.NewReader(byt)
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

		case err == io.EOF:
			err = nil
			break LOOP

		case err != nil:
			break LOOP

		case header == nil:
			continue LOOP
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, err = os.Stat(target); err != nil {
				if err = os.MkdirAll(target, 0755); err != nil {
					break LOOP
				}
			}

		case tar.TypeReg:

			dir := filepath.Dir(target)
			err = os.MkdirAll(dir, 0755)

			var file *os.File
			file, err = os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				break LOOP
			}

			if _, err = io.Copy(file, tarReader); err != nil {
				break LOOP
			}

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
