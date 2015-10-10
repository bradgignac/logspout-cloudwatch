package cloudwatch

import (
	"testing"

	"github.com/gliderlabs/logspout/router"
	. "gopkg.in/check.v1"
)

func TestFilter(t *testing.T) {
	TestingT(t)
}

type FilterSuite struct{}

var _ = Suite(&FilterSuite{})

func (s *FilterSuite) TestFiltersEmptyMessages(c *C) {
	input := make(chan Log, 1)
	output := filter(input)

	input <- &LogMessage{&router.Message{Data: ""}}
	input <- &LogMessage{&router.Message{Data: "valid"}}
	log := <-output

	c.Assert(output, HasLen, 0)
	c.Assert(log.Body(), Equals, "valid")
}
