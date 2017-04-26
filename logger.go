package commons

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
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
	logLevel := ConsConfig.GetValue("logs", "loglevel")
	logPath, err := ConsFile.FormatPath(logPath)
	if err != nil {
		logrus.Error(err)
	}
	err = ConsFile.PathExists(logPath, true)
	if err != nil {
		logrus.Error(err)
	}
	ConsLogger.ErrLog = getNew("err.log", logLevel, logPath)
	ConsLogger.MsgLog = getNew("msg.log", logLevel, logPath)
}

// @Title getNew
// @Description get new logger point
// @Parameters
//            logname          string           log name
//            loglevel         string           log level
// @Returns logger point:*logrus.Logger
func getNew(logName, logLevel, logPath string) *logger {
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
		Log:     &logrus.Logger{},
		Path:    fmt.Sprintf("%s%s", logPath, logName),
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
		file, err := ConsFile.OpenFile(l.ErrLog.Path)
		if err != nil {
			logrus.Error(err)
		}
		l.ErrLog.Log.Out = file
		l.ErrLog.Log.WithFields(logrus.Fields{
			"Time":  time.Now(),
			"Level": l.ErrLog.Log.Level,
		}).Error(err)
	}
}

// @Title LogMsg
// @Description log msg
// @Parameters
//            msg            string          msg
func (l *log) LogMsg(msg string) {
	fmt.Println(msg)
	// l.MsgLog.RWMutex.Lock()
	// defer l.MsgLog.RWMutex.Unlock()
	file, err := ConsFile.OpenFile(l.MsgLog.Path)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(file)
	l.MsgLog.Log.Out = file
	l.MsgLog.Log.WithFields(logrus.Fields{
		"Time":  time.Now(),
		"Level": l.MsgLog.Log.Level,
	}).Info(msg)
}
