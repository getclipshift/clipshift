package logger

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func init() {
	Log = logrus.New()
}

func LoggerInit(loglevel string) {
	logPath, _ := homedir.Expand("~/.clipshift/log.txt")
	logfile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Out = logfile
		stats, _ := logfile.Stat()
		if stats.Size() > 5000000 {
			os.WriteFile(logPath, make([]byte, 0), 0666)
		}
	} else {
		Log.Error("Failed to log to file, using default stderr")
	}

	switch loglevel {
	case logrus.TraceLevel.String():
		Log.SetLevel(logrus.TraceLevel)
		Log.Debug("Trace logging enabled")
	case logrus.DebugLevel.String():
		Log.SetLevel(logrus.DebugLevel)
		Log.Debug("Debug logging enabled")
	case logrus.InfoLevel.String():
		Log.SetLevel(logrus.InfoLevel)
	case logrus.WarnLevel.String():
		Log.SetLevel(logrus.WarnLevel)
	default:
		Log.SetLevel(logrus.ErrorLevel)
	}
}
