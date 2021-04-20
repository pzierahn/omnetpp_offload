package storage

import (
	"io"
	"os"
)

const (
	// 1 megabyte buffer
	bufferSize = 1024 * 1024
)

type fileChunk struct {
	size    int
	offset  int
	payload []byte
}

func fileStream(filename string) (stream chan fileChunk) {

	stream = make(chan fileChunk)

	go func() {
		defer close(stream)

		file, err := os.Open(filename)
		if err != nil {
			logger.Fatalln(err)
		}
		defer func() { _ = file.Close() }()

		//stat, err := file.Stat()
		//if err != nil {
		//	logger.Fatalln(err)
		//}
		//
		//fileSize := stat.Size()

		offset := 0

		for {
			// TODO: Check for memory leaks
			buffer := make([]byte, bufferSize)

			var size int
			size, err = file.Read(buffer)
			if err == io.EOF {
				break
			}

			if err != nil {
				logger.Fatalln(err)
			}

			stream <- fileChunk{
				size:    size,
				offset:  offset,
				payload: buffer,
			}

			offset += size

			//percent := float64(offset) / float64(fileSize) * 100.0
			//logger.Println("Percent", percent)
		}
	}()

	return stream
}
