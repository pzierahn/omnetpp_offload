package simple

import (
	"net/url"
	"strconv"
)

func QueryBool(values url.Values, key string, fallback bool) (value bool) {

	value = fallback

	queryQ, _ := values[key]
	if len(queryQ) == 1 {
		value, _ = strconv.ParseBool(queryQ[0])
	}

	return
}

func QueryString(values url.Values, key, fallback string) (value string) {

	value = fallback

	queryQ, _ := values[key]
	if len(queryQ) == 1 {
		value = queryQ[0]
	}

	return
}

func QueryInt(values url.Values, key string, fallback int) (value int) {

	value = fallback

	queryQ, _ := values[key]
	if len(queryQ) == 1 {
		if dist, err := strconv.Atoi(queryQ[0]); err == nil {
			value = dist
		}
	}

	return
}
