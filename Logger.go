package net

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	logger = log.New(ioutil.Discard, ``, log.LstdFlags) // logger is for debugging in epsp
	logMu  sync.Mutex
)

// SetLogger set logger.
func SetLogger(l *log.Logger) {
	if l == nil {
		l = log.New(ioutil.Discard, ``, log.LstdFlags) // logger is for debugging in epsp
	}
	logMu.Lock()
	logger = l
	logMu.Unlock()
}

// SetLoggerDebug set to debug mode
func SetLoggerDebug() {
	logger = log.New(os.Stderr, ``, log.LstdFlags)
}

func logf(format string, v ...interface{}) {
	logMu.Lock()
	logger.Printf(format, v...)
	logMu.Unlock()
}

func logln(v ...interface{}) {
	logMu.Lock()
	logger.Print(v...)
	logMu.Unlock()
}
