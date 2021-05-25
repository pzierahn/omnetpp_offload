package quick

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func decode(byt []byte, obj interface{}) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(byt))
	err = dec.Decode(obj)

	if err != nil {
		err = fmt.Errorf("error decoding %T: %v", obj, err)
		return
	}

	return
}

func encode(obj interface{}) (byt []byte, err error) {
	var objBuf bytes.Buffer
	objEnc := gob.NewEncoder(&objBuf)
	err = objEnc.Encode(obj)

	if err != nil {
		err = fmt.Errorf("error encoding gob %T: %v", obj, err)
		return
	}

	byt = objBuf.Bytes()

	return
}
