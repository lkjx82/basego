package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

/* ----------------------------------------------------------------------------
	async log: write to File, and roll File by dialy, File name: log File + year + month + day + .log
	init: LogInit (LogLvl, I File Name, Is Show In Console)
	close: LogFini ()
	log info: I (...)
	log debug: D (...)
	log warning: W (...)
	log error: E (...)
	log fatal: F (...)
----------------------------------------------------------------------------*/

type Hook interface {
	OnLog(le *LogEntity)
}

type LogLvl int

// ----------------------------------------------------------------------------

const (
	Dbg = LogLvl(0)
	Inf = LogLvl(1)
	War = LogLvl(2)
	Err = LogLvl(3)
	Fat = LogLvl(4)
)

// ----------------------------------------------------------------------------
// I Info
func I(v ...interface{}) {
	log(Inf, v...)
}

// ----------------------------------------------------------------------------
// I Debug
func D(v ...interface{}) {
	log(Dbg, v...)
}

// ----------------------------------------------------------------------------
// I Warning
func W(v ...interface{}) {
	log(War, v...)
}

// ----------------------------------------------------------------------------
// I Error
func E(v ...interface{}) {
	log(Err, v...)
}

// ----------------------------------------------------------------------------
// I Fatal
func F(v ...interface{}) {
	log(Fat, v...)
}

// ----------------------------------------------------------------------------

type LogEntity struct {
	Lvl  LogLvl
	Time time.Time
	Msg  string
	File string
	Func string
	Line int
}

// ----------------------------------------------------------------------------

type asynLogger struct {
	lvl      LogLvl
	c        chan *LogEntity
	console  bool
	hook     Hook
	keepDays int
}

// ----------------------------------------------------------------------------

var _logger *asynLogger = nil

// ----------------------------------------------------------------------------
// switch log Msg show in console
func SetLogConsole(show bool) {
	_logger.console = show
}

func SetKeepDays(days int) {
	_logger.keepDays = days
}

// ----------------------------------------------------------------------------
// log init
func Init(lvl LogLvl, logfile string, console bool, hook Hook) error {
	_logger = &asynLogger{}
	_logger.c = make(chan *LogEntity, 100)
	_logger.lvl = lvl
	_logger.console = console
	_logger.hook = hook
	_logger.keepDays = 7
	lastTime := time.Now()
	file, err := rollFile(nil, logfile, lastTime, lastTime)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if file != nil {
				file.Close()
			}
		}()

		for e := range _logger.c {
			if e != nil {
				var err error
				file, err = rollFile(file, logfile, lastTime, e.Time)
				if err != nil {
					return
				}

				logStr := fmt.Sprintf("%s-%s| %s  (%s:%d)\n", logLvlStr[e.Lvl],
					e.Time.Format("15:04:05.999"), e.Msg, e.Func, e.Line)
				file.Write([]byte(logStr))

				if hook != nil {
					hook.OnLog(e)
				}

				if e.Lvl == Fat {
					os.Exit(0)
				}

			} else {
				return
			}
		}
	}()

	return nil
}

// ----------------------------------------------------------------------------

func log(lvl LogLvl, v ...interface{}) {
	if lvl < _logger.lvl {
		return
	}

	fun := "Func???"
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "File???"
		line = 0
	} else {
		fun = runtime.FuncForPC(pc).Name()
	}

	msg := fmt.Sprint(v...)
	e := LogEntity{}
	e.Lvl = lvl
	e.File = file
	e.Line = line
	e.Msg = msg
	e.Func = fun
	e.Time = time.Now()

	if _logger.console {
		fmt.Println(logLvlStr[lvl], e.Time.Format("15:04:05.999"), e.Msg, e.Lvl, e.Func, e.Line)
	}
	_logger.c <- &e
}

// ----------------------------------------------------------------------------
// roll File by dialy
func rollFile(file *os.File, logfile string, fileTime time.Time, logTime time.Time) (*os.File, error) {

	if fileTime.Year() != logTime.Year() || fileTime.YearDay() != logTime.YearDay() || file == nil {
		if file != nil {
			file.Close()
		}

		// 先删除7天前的
		delLogTime := logTime.AddDate(0, 0, -_logger.keepDays+1)
		// 多往前检查 keepDays
		for i := 0; i < _logger.keepDays; i++ {
			delLogTime = delLogTime.AddDate(0, 0, -1)
			delPath := fmt.Sprintf("%s%s.log", logfile, delLogTime.Format("2006-01-02"))
			fmt.Println("删除 path: ", delPath)
			os.Remove(delPath)
		}

		idx := strings.LastIndex(logfile, "/")
		if idx > 0 {
			filePath := logfile[:idx]
			if err := os.MkdirAll(filePath, 0755); err != nil {
				fmt.Println("mkdir error", filePath)
				// os.Exit(1)
				return nil, err
			}
		}

		timeStr := logTime.Format("2006-01-02")
		var err error
		file, err = os.OpenFile(logfile+timeStr+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		return file, err
	}

	return file, nil
}

// ----------------------------------------------------------------------------
// close log
func Fina() {
	if _logger != nil {
		_logger.c <- nil
		_logger = nil
	}
}

// ----------------------------------------------------------------------------

var logLvlStr = [5]string{
	"DBG",
	"INF",
	"WAR",
	"ERR",
	"FAT",
}

// ----------------------------------------------------------------------------
