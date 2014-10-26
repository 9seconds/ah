package utils

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	logrus "github.com/Sirupsen/logrus"
)

// Logger is a global logging instance has to be used everywhere
var Logger = logrus.New()

type (
	disabledFormatter struct{}
	enabledFormatter  struct{}
)

func init() {
	EnableLogging()
}

func (df disabledFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := bytes.NewBufferString(entry.Message)
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}

func (ef enabledFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := bytes.NewBufferString(strings.ToUpper(entry.Level.String()))

	buffer.WriteString(" | ")
	buffer.WriteString(entry.Message)
	if len(entry.Data) > 0 {
		buffer.WriteString("\n     [ ")
		for key, value := range entry.Data {
			buffer.WriteString(key)
			buffer.WriteString("=<")
			fmt.Fprint(buffer, value)
			buffer.WriteString("> ")
		}
		buffer.WriteString("]")
	}
	buffer.WriteByte('\n')

	return buffer.Bytes(), nil
}

// EnableLogging enables verbose logging mode.
func EnableLogging() {
	Logger.Out = os.Stderr
	Logger.Level = logrus.DebugLevel
	Logger.Formatter = enabledFormatter{}
}

// DisableLogging disables logging.
func DisableLogging() {
	Logger.Out = os.Stderr
	Logger.Level = logrus.PanicLevel
	Logger.Formatter = disabledFormatter{}
}
