package cloudwatch

import (
	"time"

	"github.com/gliderlabs/logspout/router"
)

type Log interface {
	Body() string
	Size() int
	Timestamp() int64
}

// LogMessage represents a log message to be sent to CloudWatch.
type LogMessage struct {
	*router.Message
}

// Body returns a string representation of the log.
func (l *LogMessage) Body() string {
	return l.Data
}

// Size returns the size of the log message in bytes.
func (l *LogMessage) Size() int {
	return len(l.Data)
}

// Timestamp returns the number of milliseconds since the epoch.
func (l *LogMessage) Timestamp() int64 {
	return l.Time.UnixNano() / int64(time.Millisecond)
}

type FakeLog struct {
	size int
}

func (l *FakeLog) Body() string {
	return ""
}

func (l *FakeLog) Size() int {
	return l.size
}

func (l *FakeLog) Timestamp() int64 {
	return 0
}
