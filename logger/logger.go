package logger

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/littletwolee/commons/file"
	"github.com/sirupsen/logrus"
)

var (
	m sync.Mutex
	c *console
	l *log
)

const (
	timeFormat  string = "2006-01-02T15:04:05Z"
	fieldTime   string = "time"
	fieldFile   string = "file"
	fieldFunc   string = "func"
	fieldLine   string = "line"
	errLogger   string = "error"
	infoLogger  string = "info"
	panicLogger string = "panic"

	ERROR_PARSE = "Failed parse log level"
)

type logger struct {
	m map[string]*logrus.Logger
}

type ilogger interface {
	Error(i interface{})
	Info(i interface{})
	Panic(i interface{})
	ErrorF(format string, a ...interface{})
	InfoF(format string, a ...interface{})
	PanicF(format string, a ...interface{})
	init()
}

func field(l *logrus.Logger) *logrus.Entry {
	pc, fileName, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return l.WithFields(logrus.Fields{
		fieldFile: fileName,
		fieldFunc: funcName[:len(funcName)-2],
		fieldLine: line})
}

func (l *logger) Error(i interface{}) {
	field(l.m[errLogger]).Error(i)
}

func (l *logger) Info(i interface{}) {
	field(l.m[infoLogger]).Info(i)
}

func (l *logger) Panic(i interface{}) {
	field(l.m[panicLogger]).Panic(i)
}

func (l *logger) ErrorF(format string, a ...interface{}) {
	field(l.m[errLogger]).Errorf(format, a...)
}

func (l *logger) InfoF(format string, a ...interface{}) {
	field(l.m[infoLogger]).Infof(format, a...)
}

func (l *logger) PanicF(format string, a ...interface{}) {
	field(l.m[panicLogger]).Panicf(format, a...)
}

func (l *logger) getLogger() {
	m := make(map[string]*logrus.Logger)
	m[errLogger] = l.getNew(errLogger)
	m[infoLogger] = l.getNew(infoLogger)
	m[panicLogger] = l.getNew(panicLogger)
	l.m = m
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
	for k, v := range l.m {
		file := l.getFile(k)
		v.Out = file
	}
}

func (l *log) getFile(fileName string) *os.File {
	consFile := file.GetFile()
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

func (l *log) Error(i interface{}) {
	field(l.m[errLogger]).Error(i)
}

func (l *log) Info(i interface{}) {
	field(l.m[infoLogger]).Info(i)
}

func (l *log) Panic(i interface{}) {
	field(l.m[panicLogger]).Panic(i)
}

func (l *log) ErrorF(format string, a ...interface{}) {
	field(l.m[errLogger]).Errorf(format, a...)
}

func (l *log) InfoF(format string, a ...interface{}) {
	field(l.m[infoLogger]).Infof(format, a...)
}

func (l *log) PanicF(format string, a ...interface{}) {
	field(l.m[panicLogger]).Panicf(format, a...)
}
