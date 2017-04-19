package log_test

import (
	"github.com/lkjx82/basego/log"
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
			log.Log(sr, i)
			<-time.After(time.Second)
		}(i)
	}

	<-time.After(time.Second)
	log.Fini()
}
