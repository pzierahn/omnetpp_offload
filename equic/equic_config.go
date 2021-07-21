package equic

import (
	"github.com/lucas-clemente/quic-go"
	"time"
)

var Config = &quic.Config{
	KeepAlive:      true,
	MaxIdleTimeout: time.Millisecond * 2000,
}
