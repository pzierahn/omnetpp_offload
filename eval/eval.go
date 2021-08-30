package eval

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
)

var (
	ScenarioId   = ""
	SimulationId = ""
	TrailId      = ""
)

const (
	_ = iota
	StepStart
	StepSuccess
	StepError
)

func WriteRuns(filename string) {

	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	for inx, record := range rrecords {

		headers, values := MarshallCSV(record)

		if inx == 0 {
			if err := w.Write(headers); err != nil {
				log.Fatalln(err)
			}
		}

		if err := w.Write(values); err != nil {
			log.Fatalln(err)
		}
	}

	w.Flush()

	err := ioutil.WriteFile(filename, buf.Bytes(), 0600)
	if err != nil {
		log.Fatalln(err)
	}
}

func WriteTransfers(filename string) {

	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	for inx, record := range trecords {

		headers, values := MarshallCSV(record)

		if inx == 0 {
			if err := w.Write(headers); err != nil {
				log.Fatalln(err)
			}
		}

		if err := w.Write(values); err != nil {
			log.Fatalln(err)
		}
	}

	w.Flush()

	err := ioutil.WriteFile(filename, buf.Bytes(), 0600)
	if err != nil {
		log.Fatalln(err)
	}
}

func WriteSetup(filename string) {

	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)

	for inx, record := range setup {

		headers, values := MarshallCSV(record)

		if inx == 0 {
			if err := w.Write(headers); err != nil {
				log.Fatalln(err)
			}
		}

		if err := w.Write(values); err != nil {
			log.Fatalln(err)
		}
	}

	w.Flush()

	err := ioutil.WriteFile(filename, buf.Bytes(), 0600)
	if err != nil {
		log.Fatalln(err)
	}
}

//func WriteJSON(filename string) {
//
//	byt, err := json.MarshalIndent(rrecords, "", "  ")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	err = ioutil.WriteFile(filename, byt, 0600)
//	if err != nil {
//		log.Fatalln(err)
//	}
//}
