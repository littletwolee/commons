package commons

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	ConsLogger *logger
	log        *logrus.Logger
)

type logger struct{}

func init() {
	ConsLogger = &logger{}
	log = logrus.New()
	logpath := ConsConfig.GetValue("logs", "path")
	loglevel := ConsConfig.GetValue("logs", "loglevel")
	log.Formatter = new(logrus.TextFormatter)
	file, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	level, err := logrus.ParseLevel(loglevel)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")
	}
	log.Level = level
}

// @Title checkErr
// @Description check error
// @Parameters
//            err            error          error
func (e *logger) LogErr(err error) {
	if err != nil {
		log.Error(err)
	}
}

// @Title LogInfo
// @Description log info
// @Parameters
//            msg            string          msg
func (e *logger) LogInfo(msg string) {
	logrus.Info(msg)
}
