package utils

import (
	"fmt"
	"google.golang.org/grpc/metadata"
)

func MetaStringFallback(md metadata.MD, key, fallback string) (value string) {

	value = fallback

	values := md.Get(key)

	if len(values) > 0 {
		value = values[0]
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
