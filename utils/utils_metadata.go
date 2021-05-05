package utils

import (
	"fmt"
	"google.golang.org/grpc/metadata"
	"strconv"
)

func MetaStringFallback(md metadata.MD, key, fallback string) (value string) {

	value = fallback

	values := md.Get(key)

	if len(values) > 0 {
		value = values[0]
	}

	return
}

func MetaIntFallback(md metadata.MD, key string, fallback int) (value int) {

	value = fallback

	values := md.Get(key)

	if len(values) > 0 {
		var err error
		value, err = strconv.Atoi(values[0])
		if err != nil {
			value = fallback
		}
	}

	return
}

func MetaString(md metadata.MD, key string) (value string, err error) {

	values := md.Get(key)

	if len(values) > 0 {
		value = values[0]
	} else {
		err = fmt.Errorf("missing '%s' in md ", key)
	}

	return
}
