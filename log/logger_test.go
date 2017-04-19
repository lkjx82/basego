package log_test

import (
	log "basego/log"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {

	log.Init(log.Inf, "test", true)

	for i := 0; i < 10; i++ {
		go func(i int) {
			sr := ""
			for j := 0; j < i; j++ {
				sr = sr + "asdf asd "
			}
			log.I(sr, i)
			log.D(sr, i)
			log.E(sr, i)
			log.W(sr, i)
			//<-time.After(time.Second)
		}(i)
	}

	log.F("fatal", "ok")

	<-time.After(time.Second)
	log.Fini()
}
