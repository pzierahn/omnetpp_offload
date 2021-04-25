package storage

import (
	"io"
)

const (
	// 1 megabyte buffer
	//bufferSize = 1024 * 1024
	// 1 megabyte buffer
	bufferSize = 1024 * 1024 * 2
	// 3 megabyte buffer
	//bufferSize = 1024 * 1024 * 3
	// 64 kilo byte buffer
	//bufferSize = 1024 * 64
	// 4 kilo byte buffer
	//bufferSize = 1024 * 4
)

type fileChunk struct {
	size    int
	offset  int
	payload []byte
}

func streamReader(reader io.Reader) (stream chan fileChunk) {

	stream = make(chan fileChunk)

	go func() {
		defer close(stream)

		offset := 0

		for {
			// TODO: Check for memory leaks
			buffer := make([]byte, bufferSize)

			var size int
			size, err := reader.Read(buffer)
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
		}
	}()

	return stream
}
