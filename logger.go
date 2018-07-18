package commons

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	m sync.Mutex
	c *console
	l *log
)

const (
	timeFormat string = "2006-01-02T15:04:05Z"
	fieldTime  string = "time"
	fieldFile  string = "file"
	fieldFunc  string = "func"
	fieldLine  string = "line"
)

type logger struct {
	e, i, p *logrus.Logger
}

type ilogger interface {
	// &logger{}
	Error(format string, a ...interface{})
	Info(format string, a ...interface{})
	Panic(format string, a ...interface{})
	init()
}

func (l *logger) Error(format string, a ...interface{}) {
	pc, fileName, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	l.e.WithFields(logrus.Fields{
		//fieldTime: time.Now().Format(timeFormat),
		fieldFile: fileName,
		fieldFunc: funcName[:len(funcName)-2],
		fieldLine: line}).Errorf(format, a...)
}

func (l *logger) Info(format string, a ...interface{}) {
	l.i.WithField(fieldTime,
		time.Now().Format(timeFormat)).Infof(format, a...)
}

func (l *logger) Panic(format string, a ...interface{}) {
	pc, fileName, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	l.p.WithFields(logrus.Fields{
		fieldTime: time.Now().Format(timeFormat),
		fieldFile: fileName,
		fieldFunc: funcName[:len(funcName)-2],
		fieldLine: line}).Panicf(format, a...)
}

func (l *logger) getLogger() {
	l.e = l.getNew("error")
	l.i = l.getNew("info")
	l.p = l.getNew("panic")
}

func (l *logger) getNew(logLevel string) *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		ForceColors:            false,
		DisableColors:          false,
		FullTimestamp:          true,
		DisableSorting:         false,
		DisableLevelTruncation: true,
		QuoteEmptyFields:       false}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Error(ERROR_PARSE)
	}
	log.Level = level
	return log
}

func Console() ilogger {
	return _console()
}

type console struct {
	logger
}

func _console() *console {
	if c == nil {
		m.Lock()
		defer m.Unlock()
		if c == nil {
			c = &console{}
			c.init()
		}
	}
	return c
}

func (c *console) init() {
	c.getLogger()
}

func Log(logPath string) ilogger {
	return _log(logPath)
}

type log struct {
	logger
	m       sync.Mutex
	logPath string
}

func _log(logPath string) *log {
	if l == nil {
		m.Lock()
		defer m.Unlock()
		if l == nil {
			l = &log{logPath: logPath}
			l.init()
		}
	}
	return l
}

func (l *log) init() {
	l.getLogger()
}

func (l *log) getFile(fileName string) *os.File {
	consFile := GetFile()
	logPath, err := consFile.FormatPath(l.logPath)
	if err != nil {
		logrus.Error(err)
	}
	file, err := consFile.OpenFile(fmt.Sprintf("%s%s.log", logPath, fileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY)
	if err != nil {
		logrus.Error(err)
	}
	return file
}

func (l *log) Error(format string, a ...interface{}) {
	file := l.getFile(l.logger.e.Level.String())
	l.e.Out = file
	l.logger.Error(format, a)
}

func (l *log) Info(format string, a ...interface{}) {
	file := l.getFile(l.logger.i.Level.String())
	l.e.Out = file
	l.logger.Info(format, a...)
}

func (l *log) Panic(format string, a ...interface{}) {
	file := l.getFile(l.logger.p.Level.String())
	l.e.Out = file
	l.logger.Panic(format, a...)
}
