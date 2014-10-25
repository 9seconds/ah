package utils

import (
	"os"
	"bytes"

	logrus "github.com/Sirupsen/logrus"
)

var Logger = logrus.New()

type disabledFormatter struct {}

func init() {
	EnableLogging()
}

func (df disabledFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := bytes.NewBufferString(entry.Message)
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}

func EnableLogging() {
	Logger.Out = os.Stderr
	Logger.Level = logrus.DebugLevel
}

func DisableLogging() {
	Logger.Out = os.Stderr
	Logger.Level = logrus.PanicLevel
	Logger.Formatter = disabledFormatter{}
}
