package utils

import (
	"bytes"
	"os"

	logrus "github.com/Sirupsen/logrus"
)

// Logger is a global logging instance has to be used everywhere
var Logger = logrus.New()

type disabledFormatter struct{}

func init() {
	EnableLogging()
}

func (df disabledFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := bytes.NewBufferString(entry.Message)
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}

// EnableLogging enables verbose logging mode.
func EnableLogging() {
	Logger.Out = os.Stderr
	Logger.Level = logrus.DebugLevel
}

// DisableLogging disables logging.
func DisableLogging() {
	Logger.Out = os.Stderr
	Logger.Level = logrus.PanicLevel
	Logger.Formatter = disabledFormatter{}
}
