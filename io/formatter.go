// enhanced https://github.com/t-tomalak/logrus-easy-formatter with fields tag added
package io

import (
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// Default log format will output [INFO]: 2006-01-02T15:04:05Z07:00 - Log message
	defaultLogFormat       = "[%lvl%]: %time% - %msg% %fields%"
	defaultTimestampFormat = time.RFC3339
)

// Formatter implements logrus.Formatter interface.
type Formatter struct {
	// Timestamp format
	TimestampFormat string
	// Available standard keys: time, msg, lvl
	// Also can include custom fields but limited to strings.
	// All of fields need to be wrapped inside %% i.e %time% %msg%
	LogFormat string
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%msg%", entry.Message, 1)

	level := strings.ToUpper(entry.Level.String())
	output = strings.Replace(output, "%lvl%", level, 1)

	if strings.Contains(output, "%fields%") {

		data := "{ "
		for k, val := range entry.Data {
			s := "unknown"
			switch v := val.(type) {
			case string:
				s = v
			case int:
				s = strconv.Itoa(v)
			case bool:
				s = strconv.FormatBool(v)
			case error:
				s = v.Error()
			}

			data += k + "=" + s + ", "
		}
		if len(data) > 2 {
			data = data[:len(data)-2]
		}
		data += " }"
		if data != "{  }" {
			output = strings.Replace(output, "%fields%", data, 1)
		} else {
			output = strings.Replace(output, "%fields%", "", 1)
		}
	}

	return []byte(output), nil
}