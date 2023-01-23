package log

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type myFormatter struct{}

func (m *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// NOTE (xiangjun.ma) usage example <https://cloud.tencent.com/developer/article/1830710>
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	fName := filepath.Base(entry.Caller.File)
	newLog = fmt.Sprintf(
		"[%s] [%s] [%s:%d] %s\n",
		timestamp, entry.Level, fName, entry.Caller.Line, entry.Message,
	)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func init() {
}

func GetLogger(ctx context.Context) *logrus.Entry {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&myFormatter{})
	return logrus.WithContext(ctx)
}
