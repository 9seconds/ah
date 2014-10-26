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

func (df disabledFormatter) Format(entry *logrus.Entry) (buf []byte, err error) {
	buffer := bytes.NewBufferString(entry.Message)
	buffer.WriteByte('\n')
	buf = buffer.Bytes()

	return
}

func (ef enabledFormatter) Format(entry *logrus.Entry) (buf []byte, err error) {
	level := strings.ToUpper(entry.Level.String())
	buffer := bytes.NewBufferString(level)

	indent := make([]byte, len(level)+1)
	for idx := 0; idx < len(indent); idx++ {
		indent[idx] = ' '
	}

	buffer.WriteString(" | ")
	buffer.WriteString(entry.Message)
	if len(entry.Data) > 0 {
		for key, value := range entry.Data {
			buffer.WriteByte('\n')
			buffer.Write(indent)
			buffer.WriteString("* ")
			buffer.WriteString(key)
			buffer.WriteString(" -> ")
			fmt.Fprint(buffer, value)
		}
	}
	buffer.WriteByte('\n')
	buf = buffer.Bytes()

	return
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
