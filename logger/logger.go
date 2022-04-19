package logger

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func LoggerInit(loglevel string) {
	Log = logrus.New()

	switch loglevel {
	case logrus.TraceLevel.String():
		Log.SetLevel(logrus.TraceLevel)
		Log.Debug("Trace logging enabled")
		break
	case logrus.DebugLevel.String():
		Log.SetLevel(logrus.DebugLevel)
		Log.Debug("Debug logging enabled")
		break
	case logrus.InfoLevel.String():
		Log.SetLevel(logrus.InfoLevel)
		break
	case logrus.WarnLevel.String():
		Log.SetLevel(logrus.WarnLevel)
		break
	default:
		Log.SetLevel(logrus.ErrorLevel)
	}

	logPath, _ := homedir.Expand("~/.clipshift/log.txt")
	logfile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Out = logfile
	} else {
		Log.Error("Failed to log to file, using default stderr")
	}
}
