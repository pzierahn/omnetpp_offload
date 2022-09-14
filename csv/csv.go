package csv

import (
	"encoding/csv"
	"google.golang.org/protobuf/reflect/protoreflect"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Writer struct {
	mu             sync.Mutex
	file           *os.File
	csv            *csv.Writer
	headersWritten bool
}

func NewWriter(dir, filename string) (writer *Writer) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("couldn't create directory: %v", err)
	}

	path := filepath.Join(dir, filename)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	return &Writer{
		file: file,
		csv:  csv.NewWriter(file),
	}
}

func (writer *Writer) Close() {
	writer.mu.Lock()
	defer writer.mu.Unlock()

	writer.csv.Flush()
	_ = writer.file.Close()
}

func (writer *Writer) Write(records ...[]string) {
	writer.mu.Lock()
	defer writer.mu.Unlock()

	for _, record := range records {
		err := writer.csv.Write(record)
		if err != nil {
			log.Fatalf("couldn't write to csv: %v", err)
		}
	}

	writer.csv.Flush()
}

func (writer *Writer) RecordProtos(messages ...protoreflect.Message) {

	var header []string
	var record []string

	for _, message := range messages {
		head, rec := MarshallProto(message)
		header = append(header, head...)
		record = append(record, rec...)
	}

	//stat, err := server.logFile.Stat()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//if stat.Size() == 0 {
	//	if err = server.writer.Write(headers); err != nil {
	//		log.Fatalln(err)
	//	}
	//}

	if !writer.headersWritten {
		writer.headersWritten = true
		writer.Write(header)
	}

	writer.Write(record)
}
