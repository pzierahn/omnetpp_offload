package eval

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/reflect/protoreflect"
	"log"
	"reflect"
	"strings"
	"time"
)

func MarshallCSV(obj interface{}) (headers, values []string) {
	vals := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	for inx := 0; inx < vals.NumField(); inx++ {
		valueField := vals.Field(inx)
		typField := typ.Field(inx)

		header := typField.Name

		jsonTag := typField.Tag.Get("json")
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")

			if len(parts) > 0 {
				header = parts[0]
			}
		}

		headers = append(headers, header)

		switch val := valueField.Interface().(type) {
		case time.Time:
			timebyt, _ := val.MarshalText()
			values = append(values, string(timebyt))
		case string:
			values = append(values, val)
		case nil:
			values = append(values, "")
		case int:
			values = append(values, fmt.Sprint(val))
		case uint32:
			values = append(values, fmt.Sprint(val))
		case uint64:
			values = append(values, fmt.Sprint(val))
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

func MarshallProto(message protoreflect.Message) (headers, values []string) {
	fields := message.Descriptor().Fields()

	for inx := 0; inx < fields.Len(); inx++ {
		field := fields.Get(inx)
		value := message.Get(field)

		headers = append(headers, string(field.Name()))
		values = append(values, value.String())
	}

	return
}
