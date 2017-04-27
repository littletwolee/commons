package commons

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	ConsLogger *log
	m          *sync.RWMutex
)

type log struct {
	ErrLog *logger
	MsgLog *logger
}

type logger struct {
	RWMutex *sync.RWMutex
	Log     *logrus.Logger
	Path    string
}

func init() {
	ConsLogger = &log{}
	logPath := ConsConfig.GetValue("logs", "path")
	logPath, err := ConsFile.FormatPath(logPath)
	if err != nil {
		logrus.Error(err)
	}
	err = ConsFile.PathExists(logPath, true)
	if err != nil {
		logrus.Error(err)
	}
	ConsLogger.ErrLog = getNew("error", logPath)
	ConsLogger.MsgLog = getNew("info", logPath)
}

// @Title getNew
// @Description get new logger point
// @Parameters
//            loglevel         string           log level
//            logPath          string           log path
// @Returns logger point:*logrus.Logger
func getNew(logLevel, logPath string) *logger {
	var (
		log *logger
		err error
	)
	logPath, err = ConsFile.FormatPath(logPath)
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

// @Title checkErr
// @Description check error
// @Parameters
//             errin            error          error
func (l *log) LogErr(errin error) {
	if errin != nil {
		l.ErrLog.RWMutex.Lock()
		defer l.ErrLog.RWMutex.Unlock()
		file, err := ConsFile.OpenFile(l.ErrLog.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY)
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
func (l *log) LogMsg(msg string) {
	l.MsgLog.RWMutex.Lock()
	defer l.MsgLog.RWMutex.Unlock()
	file, err := ConsFile.OpenFile(l.MsgLog.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY)
	if err != nil {
		logrus.Error(err)
	}
	defer file.Close()
	l.MsgLog.Log.Out = file
	l.MsgLog.Log.Info(msg)
}
