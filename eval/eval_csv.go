package eval

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"
)

type xuint32 uint32
type xint int

func MarshallCSV(obj interface{}) (headers, values []string) {
	vals := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	for inx := 0; inx < vals.NumField(); inx++ {
		field := vals.Field(inx)
		header := typ.Field(inx).Tag.Get("csv")

		if header == "" {
			continue
		}

		headers = append(headers, header)

		switch val := field.Interface().(type) {
		case time.Time:
			timebyt, _ := val.MarshalText()
			values = append(values, string(timebyt))
		case string:
			values = append(values, val)
		case nil:
			values = append(values, "")
		case int:
			values = append(values, fmt.Sprint(val))
		case uint64:
			values = append(values, fmt.Sprint(val))
		case xuint32:
			values = append(values, fmt.Sprintf("0x%08x", val))
		case xint:
			values = append(values, fmt.Sprintf("0x%08x", val))
		case error:
			values = append(values, val.Error())
		default:
			byt, err := json.Marshal(val)
			if err != nil {
				log.Fatalln(err)
			}

			values = append(values, string(byt))
		}
	}

	return
}
