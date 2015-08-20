package cloudwatch

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func TestBatch(t *testing.T) {
	TestingT(t)
}

type BatchSuite struct {
	in  chan Log
	out chan []Log
}

var _ = Suite(&BatchSuite{})

func (s *BatchSuite) SetUpTest(c *C) {
	s.in = make(chan Log)
	s.out = make(chan []Log)
}

// TODO: Test closing channel.

func (s *BatchSuite) TestBatcherWhenNotFull(c *C) {
	batcher := NewBatcher(s.in, s.out, Capacity{Length: 2, Size: 2})
	go batcher.Start()

	s.in <- &FakeLog{size: 0}

	c.Assert(batcher.Length(), Equals, 1)
}

func (s *BatchSuite) TestFlushWhenCountReached(c *C) {
	batcher := NewBatcher(s.in, s.out, Capacity{Length: 2, Size: 4})
	go batcher.Start()

	s.in <- &FakeLog{size: 1}
	s.in <- &FakeLog{size: 1}

	messages := <-s.out

	c.Assert(messages, HasLen, 2)
	c.Assert(batcher.Length(), Equals, 0)
}

func (s *BatchSuite) TestFlushWhenSizeReached(c *C) {
	batcher := NewBatcher(s.in, s.out, Capacity{Length: 4, Size: 2})
	go batcher.Start()

	s.in <- &FakeLog{size: 1}
	s.in <- &FakeLog{size: 1}

	messages := <-s.out

	c.Assert(messages, HasLen, 2)
	c.Assert(batcher.Length(), Equals, 0)
}

func (s *BatchSuite) TestFlushWhenSizeExceeded(c *C) {
	batcher := NewBatcher(s.in, s.out, Capacity{Length: 4, Size: 3})
	go batcher.Start()

	s.in <- &FakeLog{size: 2}
	s.in <- &FakeLog{size: 2}

	messages := <-s.out

	c.Assert(messages, HasLen, 1)
	c.Assert(batcher.Length(), Equals, 1)
}

func (s *BatchSuite) TestFlushOversizedMessage(c *C) {
	batcher := NewBatcher(s.in, s.out, Capacity{Length: 4, Size: 1})
	go batcher.Start()

	s.in <- &FakeLog{size: 2}

	messages := <-s.out

	c.Assert(messages, HasLen, 1)
	c.Assert(batcher.Length(), Equals, 0)
}

func (s *BatchSuite) TestFlushWhenDurationExceeded(c *C) {
	duration := 1 * time.Millisecond
	batcher := NewBatcher(s.in, s.out, Capacity{Length: 2, Size: 2, Duration: duration})
	go batcher.Start()

	s.in <- &FakeLog{size: 1}
	time.Sleep(2 * time.Millisecond)

	messages := <-s.out

	c.Assert(messages, HasLen, 1)
	c.Assert(batcher.Length(), Equals, 0)
}

func (s *BatchSuite) TestFlushWhenClosed(c *C) {
	batcher := NewBatcher(s.in, s.out, Capacity{Length: 2, Size: 2})
	go batcher.Start()

	s.in <- &FakeLog{size: 0}
	close(s.in)

	messages := <-s.out

	c.Assert(messages, HasLen, 1)
	c.Assert(batcher.Length(), Equals, 0)
}
