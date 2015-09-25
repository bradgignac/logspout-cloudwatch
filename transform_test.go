package cloudwatch

import (
	"testing"

	"github.com/gliderlabs/logspout/router"
	. "gopkg.in/check.v1"
)

func TestTransform(t *testing.T) {
	TestingT(t)
}

type TransformSuite struct{}

var _ = Suite(&TransformSuite{})

func (s *TransformSuite) TestTransformsMessageToLog(c *C) {
	messages := make(chan *router.Message)
	logs := transform(messages)

	messages <- &router.Message{Data: "hello world"}
	log := <-logs

	c.Assert(log.Body(), Equals, "hello world")
}
