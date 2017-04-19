package log

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

/* ----------------------------------------------------------------------------
	async log: write to file, and roll file by dialy, file name: log file + year + month + day + .log
	init: LogInit (LogLvl, Log File Name, Is Show In Console)
	close: LogFini ()
	log info: Log (...)
	log debug: LogD (...)
	log warning: LogE (...)
	log error: LogE (...)
	log fatal: LogF (...)
----------------------------------------------------------------------------*/

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
// Log Info
func Log(v ...interface{}) {
	log(Inf, v...)
}

// ----------------------------------------------------------------------------
// Log Debug
func LogD(v ...interface{}) {
	log(Dbg, v...)
}

// ----------------------------------------------------------------------------
// Log Warning
func LogW(v ...interface{}) {
	log(War, v...)
}

// ----------------------------------------------------------------------------
// Log Error
func LogE(v ...interface{}) {
	log(Err, v...)
}

// ----------------------------------------------------------------------------
// Log Fatal
func LogF(v ...interface{}) {
	log(Fat, v...)
}

// ----------------------------------------------------------------------------

type logEntity struct {
	lvl  LogLvl
	time time.Time
	msg  string
	file string
	line int
}

// ----------------------------------------------------------------------------

type asynLogger struct {
	lvl     LogLvl
	c       chan *logEntity
	console bool
}

// ----------------------------------------------------------------------------

var _logger *asynLogger = nil

// ----------------------------------------------------------------------------
// switch log msg show in console
func SetLogConsole(show bool) {
	_logger.console = show
}

// ----------------------------------------------------------------------------
// log init
func Init(lvl LogLvl, logfile string, console bool) error {
	_logger = &asynLogger{}
	_logger.c = make(chan *logEntity, 100)
	_logger.lvl = lvl
	_logger.console = console
	lastTime := time.Now()
	file, err := logRollFile(nil, logfile, lastTime, lastTime)
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
				file, err = logRollFile(file, logfile, lastTime, e.time)
				if err != nil {
					return
				}

				logStr := fmt.Sprintf("%s-%s| %s  (%s:%d)\n", logLvlStr[e.lvl], e.time.Format("15:04:05.999"), e.msg, e.file, e.line)
				file.Write([]byte(logStr))

				if e.lvl == Fat {
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
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	msg := fmt.Sprint(v...)
	e := logEntity{}
	e.lvl = lvl
	e.file = file
	e.line = line
	e.msg = msg
	e.time = time.Now()

	if _logger.console {
		fmt.Println(logLvlStr[lvl], e.time.Format("15:04:05.999"), e.msg, e.lvl, e.file, e.line)
	}
	_logger.c <- &e
}

// ----------------------------------------------------------------------------
// roll file by dialy
func logRollFile(file *os.File, logfile string, fileTime time.Time, logTime time.Time) (*os.File, error) {

	if fileTime.Year() != logTime.Year() || fileTime.YearDay() != logTime.YearDay() || file == nil {
		if file != nil {
			file.Close()
		}
		timeStr := logTime.Format("2006-01-02")
		var err error
		file, err = os.OpenFile(logfile+timeStr+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		return file, err
	}

	return file, nil
}

// ----------------------------------------------------------------------------
// close log
func LogFini() {
	if _logger != nil {
		_logger.c <- nil
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
