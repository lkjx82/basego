package log_test

import (
	"fmt"
	"github.com/lkjx82/basego/log"
	"testing"
	"time"
)

type Hook struct {
}

func (this *Hook) OnLog(lg *log.LogEntity) {
	fmt.Println(lg.Lvl, lg.Func, lg.Line, lg.Msg, lg.Time)
}

func TestLogger(t *testing.T) {

	hook := Hook{}

	log.Init(log.Inf, "test", false, &hook)

	for i := 0; i < 10; i++ {
		go func(i int) {
			sr := ""
			for j := 0; j < i; j++ {
				sr = sr + "asdf asd "
			}
			log.I(sr, i)
			<-time.After(time.Second)
		}(i)
	}

	<-time.After(time.Second)
	log.Fina()
}
