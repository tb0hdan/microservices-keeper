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
	ls.InitLogger()
	return ls.Logger
}

type SlackLogger struct {
	Logger *logrus.Logger
}

func (ls *SlackLogger) Output(calldepth int, s string) (err error) {
	ls.Logger.Println(s)
	return
}

func NewSlackLogger() *SlackLogger {
	return &SlackLogger{Logger: Logger}
}

func SetDebug() {
	Logger.SetLevel(logrus.DebugLevel)
}

func GetDebug() bool {
	return Logger.Level == logrus.DebugLevel
}

func init() { // nolint
	lg := &LoggerStruct{}
	Logger = lg.GetLogger()
}
