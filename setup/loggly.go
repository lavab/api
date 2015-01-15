package setup

import (
	"github.com/lavab/logrus"
	"github.com/segmentio/go-loggly"
)

type logglyHook struct {
	Loggly *loggly.Client
}

func (h *logglyHook) Fire(entry *logrus.Entry) error {
	entry.Data["message"] = entry.Message
	return h.Loggly.Send(map[string]interface{}(entry.Data))
}

func (h *logglyHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
