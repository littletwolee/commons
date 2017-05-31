package commons

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sync"
)

var (
	m          *sync.RWMutex
	consLogger *Log
)

type Log struct {
	ErrLog   *logger
	MsgLog   *logger
	PanicLog *logger
}

type logger struct {
	RWMutex *sync.RWMutex
	Log     *logrus.Logger
	Path    string
}

func GetLogger() *Log {
	if consLogger != nil {
		return consLogger
	}
	consLogger = &Log{}
	logPath := ConsConfig.GetValue("logs", "path")
	logPath, err := consFile.FormatPath(logPath)
	if err != nil {
		logrus.Panic(err)
	}
	err = consFile.PathExists(logPath, true)
	if err != nil {
		logrus.Panic(err)
	}
	consLogger.ErrLog = consLogger.getNew("error", logPath)
	consLogger.MsgLog = consLogger.getNew("info", logPath)
	consLogger.PanicLog = consLogger.getNew("panic", logPath)
	return consLogger
}

// @Title getNew
// @Description get new logger point
// @Parameters
//            loglevel         string           log level
//            logPath          string           log path
// @Returns logger point:*logrus.Logger
func (l *Log) getNew(logLevel, logPath string) *logger {
	var (
		log      *logger
		err      error
		consFile *File
	)
	consFile = GetFile()
	logPath, err = consFile.FormatPath(logPath)
	if err != nil {
		logrus.Error(err)
	}
	log = &logger{
		RWMutex: new(sync.RWMutex),
		Log:     logrus.New(),
		Path:    fmt.Sprintf("%s%s.log", logPath, logLevel),
	}
	log.Log.Formatter = new(logrus.TextFormatter)
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Error(ErrorParse)
	}
	log.Log.Level = level
	goto RETURN
RETURN:
	return log
}

// @Title LogErr
// @Description log error
// @Parameters
//             errin            error          error
func (l *Log) LogErr(errin error) {
	if errin != nil {
		l.ErrLog.RWMutex.Lock()
		defer l.ErrLog.RWMutex.Unlock()
		file, err := consFile.OpenFile(l.ErrLog.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY)
		if err != nil {
			logrus.Error(err)
		}
		defer file.Close()
		l.ErrLog.Log.Out = file
		l.ErrLog.Log.Error(errin)
	}
}

// @Title LogMsg
// @Description log msg
// @Parameters
//            msg            string          msg
func (l *Log) LogMsg(msg string) {
	l.MsgLog.RWMutex.Lock()
	defer l.MsgLog.RWMutex.Unlock()
	file, err := consFile.OpenFile(l.MsgLog.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY)
	if err != nil {
		logrus.Error(err)
	}
	defer file.Close()
	l.MsgLog.Log.Out = file
	l.MsgLog.Log.Info(msg)
}

// @Title LogPanic
// @Description log panic
// @Parameters
//            msg            string          msg
func (l *Log) LogPanic(errin error) {
	if errin != nil {
		l.PanicLog.RWMutex.Lock()
		defer l.PanicLog.RWMutex.Unlock()
		file, err := consFile.OpenFile(l.PanicLog.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY)
		if err != nil {
			logrus.Error(err)
		}
		defer file.Close()
		l.PanicLog.Log.Out = file
		pc, _, line, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()
		l.PanicLog.Log.WithFields(logrus.Fields{
			"func": funcName[:len(funcName)-2],
			"line": line,
		}).Panic(errin)
	}
}

// @Title OutErr
// @Description out put error
// @Parameters
//             errin            error          error
func (l *Log) OutErr(errin error) {
	if errin != nil {
		l.ErrLog.Log.Info(errin)
	}
}

// @Title OutMsg
// @Description out put msg
// @Parameters
//            msg            string          msg
func (l *Log) OutMsg(msg string) {
	l.MsgLog.Log.Info(msg)
}

// @Title OutPanic
// @Description out put panic
// @Parameters
//            errin            error          error
func (l *Log) OutPanic(errin error) {
	if errin != nil {
		pc, _, line, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()
		l.PanicLog.Log.WithFields(logrus.Fields{
			"func": funcName[:len(funcName)-2],
			"line": line,
		}).Panic(errin)
	}
}
