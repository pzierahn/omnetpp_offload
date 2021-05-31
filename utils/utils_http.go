package utils

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
)

func Marshal(message proto.Message) (bytes []byte, err error) {
	//marshal := protojson.MarshalOptions{
	//	Indent:          "  ",
	//	EmitUnpopulated: true,
	//}
	//
	//bytes, err = marshal.Marshal(message)

	bytes, err = json.MarshalIndent(message, "", "  ")

	return
}

func Response(writer http.ResponseWriter, message proto.Message, proto bool) {
	if proto {
		writer.Header().Set("Content-Type", "application/protobuf")
		ResponseProto(writer, message)
	} else {
		writer.Header().Set("Content-Type", "application/json")
		ResponseJson(writer, message)
	}
}

func ResponseJson(writer http.ResponseWriter, message proto.Message) {
	if bytes, err := Marshal(message); err != nil {
		writer.WriteHeader(503)
		_, _ = fmt.Fprint(writer, err.Error())
		log.Println(err)
	} else {
		_, _ = writer.Write(bytes)
	}
}

func ResponseProto(writer http.ResponseWriter, message proto.Message) {
	out, err := proto.Marshal(message)
	if err != nil {
		writer.WriteHeader(503)
		log.Println("Failed to encode message :", err)
		return
	}

	_, _ = writer.Write(out)
}
