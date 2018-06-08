package log

import (
	"fmt"
	"github.com/iris-contrib/color"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	UNSPECIFIED int = iota
	TRACE
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var LOG_LEVEL_MAP = map[int]string{
	UNSPECIFIED: "UNSPECIFIED",
	TRACE:       "TRACE",
	DEBUG:       "DEBUG",
	INFO:        "INFO",
	WARN:        "WARN",
	ERROR:       "ERROR",
	FATAL:       "FATAL",
}

const (
	CALL_PATH   = 2
	TIME_FORMAT = "2006-01-02 15:04:05.0000"
)

type Logger struct {
	Writer io.Writer
	Level  int
	mu     sync.Mutex
}

var Log = NewLogger(os.Stdout, DEBUG)

func getShortFileName(file string) string {
	index := strings.LastIndex(file, "/")
	return file[index+1:]
}

func attr(sgr int) color.Attribute {
	return color.Attribute(sgr)
}

func NewLogger(writer io.Writer, level int) *Logger {
	return &Logger{
		Writer: writer,
		Level:  level,
	}
}

type LoggerIface interface {
	Log(level int, v ...interface{})
	Logf(level int, formater string, v ...interface{})
}

func (l *Logger) SetLevel(level int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Level = level
}

func (l *Logger) SetLevelByName(level string) {
	switch level {
	case "ERROR":
		{
			l.SetLevel(ERROR)
		}
	case "WARN":
		{
			l.SetLevel(WARN)
		}
	case "INFO":
		{
			l.SetLevel(INFO)
		}
	case "DEBUG":
		{
			l.SetLevel(DEBUG)
		}
	case "TRACE":
		{
			l.SetLevel(TRACE)
		}
	default:
		l.SetLevel(WARN)
	}
}

func (l *Logger) Log(level int, v ...interface{}) {
	timestamp := time.Now().Format(TIME_FORMAT)
	loglevel := LOG_LEVEL_MAP[level]
	l.mu.Lock()
	defer l.mu.Unlock()
	context := fmt.Sprint(v...)
	pc, file, line, _ := runtime.Caller(CALL_PATH)
	funcname := runtime.FuncForPC(pc).Name()
	file = getShortFileName(file)
	log := fmt.Sprintf("%s [%s] %s [%s] [%s:%d]", timestamp, loglevel, context, funcname, file, line)
	fmt.Fprintln(l.Writer, log)
}

func (l *Logger) Logf(level int, format string, v ...interface{}) {
	timestamp := time.Now().Format(TIME_FORMAT)
	loglevel := LOG_LEVEL_MAP[level]
	l.mu.Lock()
	defer l.mu.Unlock()
	context := fmt.Sprintf(format, v...)
	pc, file, line, _ := runtime.Caller(CALL_PATH)
	funcname := runtime.FuncForPC(pc).Name()
	file = getShortFileName(file)

	log := fmt.Sprintf("%s [%s] %s [%s] [%s:%d]", timestamp, loglevel, context, funcname, file, line)

	if level == ERROR {
		c := color.New(l.Writer, attr(int(color.FgHiRed)))
		coloredFmtPrinter := c.SprintFunc()
		fmt.Fprintln(l.Writer, coloredFmtPrinter(log))
	} else {
		fmt.Fprintln(l.Writer, log)
	}

}

func (l *Logger) Traceln(v ...interface{}) {
	if TRACE >= l.Level {
		l.Log(TRACE, v...)
	}
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	if TRACE >= l.Level {
		l.Logf(TRACE, format, v...)
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	if DEBUG >= l.Level {
		l.Log(DEBUG, v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if DEBUG >= l.Level {
		l.Logf(DEBUG, format, v...)
	}
}

func (l *Logger) Infoln(v ...interface{}) {
	if INFO >= l.Level {
		l.Log(INFO, v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if INFO >= l.Level {
		l.Logf(INFO, format, v...)
	}
}

func (l *Logger) Warnln(v ...interface{}) {
	if WARN >= l.Level {
		l.Log(WARN, v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if WARN >= l.Level {
		l.Logf(WARN, format, v...)
	}
}

func (l *Logger) Errorln(v ...interface{}) {
	if ERROR >= l.Level {
		l.Log(ERROR, v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if ERROR >= l.Level {

		l.Logf(ERROR, format, v...)
	}
}

func (l *Logger) Fatalln(v ...interface{}) {
	if FATAL >= l.Level {
		l.Log(FATAL, v...)
		os.Exit(1)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if FATAL >= l.Level {
		l.Logf(FATAL, format, v...)
		os.Exit(1)
	}
}
