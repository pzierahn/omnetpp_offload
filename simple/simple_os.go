package simple

import (
	"os"
	"strings"
)

func GetHostnameShort() (host string) {
	host, _ = os.Hostname()

	host = strings.TrimSuffix(host, ".local")
	host = strings.TrimSuffix(host, ".fritz.box")

	return
}
