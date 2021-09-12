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
		field := vals.Field(inx)
		tag := typ.Field(inx).Tag.Get("json")
		parts := strings.Split(tag, ",")

		if len(parts) < 1 {
			continue
		}

		header := parts[0]

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
