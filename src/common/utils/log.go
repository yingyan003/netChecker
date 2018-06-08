package utils

import mylog "github.com/maxwell92/gokits/log"

var log = mylog.Log

func SetLogLevel(level int) {
	log.SetLevel(level)
}

func GetLog() *mylog.Logger {
	return log
}
