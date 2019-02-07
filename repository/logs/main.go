package logs

import "github.com/sirupsen/logrus" // nolint

var Logger *logrus.Logger

type LoggerStruct struct {
	Logger *logrus.Logger
}

func (ls *LoggerStruct) InitLogger() {
	if ls.Logger == nil {
		ls.Logger = logrus.New()
	}
}

func (ls *LoggerStruct) SetLevel(level logrus.Level) {
	ls.InitLogger()
	ls.Logger.SetLevel(level)
}

func (ls *LoggerStruct) GetLogger() *logrus.Logger {
	return ls.Logger
}

func SetDebug() {
	Logger.SetLevel(logrus.DebugLevel)
}

func init() { // nolint
	lg := &LoggerStruct{}
	Logger = lg.GetLogger()
}
