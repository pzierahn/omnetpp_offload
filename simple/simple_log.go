package simple

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func WriteLogToFile(prefix, path string) {

	date := time.Now()
	logName := fmt.Sprintf("%s.%s.%s.log",
		prefix, date.Format("2006-02-01"), date.Format("15-04-05"))

	logDir := filepath.Join(path, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalln(err)
	}

	logPath := filepath.Join(logDir, logName)
	log.Printf("Write logs to %s", logPath)

	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stderr, file)
	log.SetOutput(mw)
}
