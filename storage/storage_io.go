package storage

import (
	"github.com/pzierahn/project.go.omnetpp/simple"
	"io"
)

const bufferSize = simple.MEGABYTE * 2

type fileChunk struct {
	size    int
	offset  int
	payload []byte
}

func streamReader(reader io.Reader) (stream chan *fileChunk) {

	stream = make(chan *fileChunk)

	go func() {
		defer close(stream)

		offset := 0

		for {
			// TODO: Check for memory leaks
			buffer := make([]byte, bufferSize)

			size, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln(err)
			}

			stream <- &fileChunk{
				size:    size,
				offset:  offset,
				payload: buffer[:size],
			}

			offset += size
		}
	}()

	return stream
}
