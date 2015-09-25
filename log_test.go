package cloudwatch

import (
	"bytes"
	"testing"
	"time"

	"github.com/gliderlabs/logspout/router"
	. "gopkg.in/check.v1"
)

func TestLog(t *testing.T) {
	TestingT(t)
}

type LogSuite struct{}

var _ = Suite(&LogSuite{})

func (s *LogSuite) TestSize(c *C) {
	var buffer bytes.Buffer

	for i := 0; i < 2048; i++ {
		buffer.WriteString("0")
	}

	msg := &router.Message{Data: buffer.String()}
	log := LogMessage{msg}

	c.Assert(log.Size(), Equals, 2048)
}

func (s *LogSuite) TestTimestamp(c *C) {
	now := time.Now()
	msg := &router.Message{Time: now}
	log := LogMessage{msg}

	c.Assert(log.Timestamp(), Equals, now.UnixNano()/int64(time.Millisecond))
}
