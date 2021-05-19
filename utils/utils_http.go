package utils

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
)

func Marshal(message proto.Message) (bytes []byte, err error) {
	marshal := protojson.MarshalOptions{
		Indent:          "  ",
		EmitUnpopulated: true,
	}

	bytes, err = marshal.Marshal(message)

	return
}

func Response(writer http.ResponseWriter, message proto.Message, proto bool) {
	if proto {
		ResponseProto(writer, message)
	} else {
		ResponseJson(writer, message)
	}
}

func ResponseJson(writer http.ResponseWriter, message proto.Message) {
	if bytes, err := Marshal(message); err != nil {
		writer.WriteHeader(400)
		_, _ = fmt.Fprint(writer, err.Error())
		log.Println(err)
	} else {
		_, _ = writer.Write(bytes)
	}
}

func ResponseProto(writer http.ResponseWriter, message proto.Message) {
	out, err := proto.Marshal(message)
	if err != nil {
		log.Println("Failed to encode message :", err)
		return
	}

	_, _ = writer.Write(out)
}
