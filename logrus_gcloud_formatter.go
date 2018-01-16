package logrus_gcloud_formatter

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

type LogrusGoogleCloudFormatter struct {
	Type string // if not empty use for logstash type field.

	// TimestampFormat sets the format used for timestamps.
	TimestampFormat string
}

func levelToString(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "debug"
	case logrus.InfoLevel:
		return "info"
	case logrus.WarnLevel:
		return "warning"
	case logrus.ErrorLevel:
		return "error"
	case logrus.FatalLevel:
		return "critical"
	case logrus.PanicLevel:
		return "critical"
	}

	return "info"
}

func (f *LogrusGoogleCloudFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := make(logrus.Fields)

	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			fields[k] = v.Error()
		default:
			fields[k] = v
		}
	}

	epoch := entry.Time.Unix()
	epoch_micros := epoch*int64(1000000) + int64(entry.Time.Nanosecond()/1000)

	fields["timestamp_micros"] = epoch_micros
	fields["timestamp"] = epoch
	fields["message"] = entry.Message
	fields["severity"] = levelToString(entry.Level)

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
