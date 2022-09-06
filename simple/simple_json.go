package simple

import (
	"compress/gzip"
	"encoding/json"
	"log"
	"os"
)

func WritePrettyJson(filename string, bytes []byte) {

	var tmp interface{}
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		log.Panic(err)
	}

	jBytes, err := json.MarshalIndent(tmp, "", "    ")
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(filename, jBytes, 0755)
	if err != nil {
		log.Panic(err)
	}
}

func Prettify(data interface{}) ([]byte, string) {

	jBytes, err := PrettyBytesErr(data)
	if err != nil {
		log.Panic(err)
	}

	return jBytes, string(jBytes)
}

func PrettyString(data interface{}) string {

	_, jString := Prettify(data)
	return jString
}

func PrettyBytes(data interface{}) []byte {

	jBytes, _ := Prettify(data)
	return jBytes
}

func PrettyBytesErr(data interface{}) ([]byte, error) {

	jBytes, err := json.MarshalIndent(data, "", "    ")
	return jBytes, err
}

func WritePretty(filename string, data interface{}) {

	bytes, err := json.MarshalIndent(data, "", "    ")
	err = os.WriteFile(filename, bytes, 0755)
	if err != nil {
		log.Panic(err)
	}
}

func WritePrettyGz(filename string, data interface{}) {

	file, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644)

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		_ = file.Close()
	}()

	wr := gzip.NewWriter(file)
	defer func() {
		_ = wr.Close()
	}()

	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Panic(err)
	}

	_, err = wr.Write(bytes)
	if err != nil {
		log.Panic(err)
	}
}

func WritePrettyBytes(filename string, data []byte) {

	var tmp interface{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		log.Panic(err)
	}

	bytes, err := json.MarshalIndent(tmp, "", "    ")
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(filename, bytes, 0755)
	if err != nil {
		log.Panic(err)
	}
}

func RWPrettify(filename string) {

	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	WritePrettyJson(filename, bytes)
}

func UnmarshallFile(filepath string, obj interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	err = json.NewDecoder(file).Decode(&obj)

	return err
}
