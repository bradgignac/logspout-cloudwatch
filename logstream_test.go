package cloudwatch

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/bradgignac/logspout-cloudwatch/test"
	. "gopkg.in/check.v1"
)

const REGION = "us-west-2"

func TestLogStream(t *testing.T) {
	TestingT(t)
}

type LogStreamSuite struct {
	mock    *test.CloudWatchLogsMock
	session *session.Session
	stream  *LogStream
}

var _ = Suite(&LogStreamSuite{})

func (s *LogStreamSuite) SetUpTest(c *C) {
	s.mock = test.NewCloudWatchLogsMock()

	config := aws.NewConfig().
		WithEndpoint(s.mock.URL).
		WithRegion(REGION).
		WithDisableSSL(true)
	session := session.New(config)

	s.stream = NewLogStream("group", "stream", session)
}

func (s *LogStreamSuite) TearDownTest(c *C) {
	s.mock.Close()
}

func (s *LogStreamSuite) TestNewStream(c *C) {
	err := s.stream.Init()
	streams := s.mock.GetStreams("group")

	c.Assert(err, IsNil)
	c.Assert(s.stream.Token, IsNil)
	c.Assert(streams, HasLen, 1)
}

func (s *LogStreamSuite) TestExistingStream(c *C) {
	s.mock.AddStream("group", "stream")

	err := s.stream.Init()
	streams := s.mock.GetStreams("group")

	c.Assert(err, IsNil)
	c.Assert(s.stream.Token, NotNil)
	c.Assert(streams, HasLen, 1)
}

func (s *LogStreamSuite) TestPutLogsToNewStream(c *C) {
	s.mock.AddStream("group", "stream")

	logs := []*cloudwatchlogs.InputLogEvent{
		&cloudwatchlogs.InputLogEvent{
			Message:   aws.String("body"),
			Timestamp: aws.Int64(0),
		},
	}

	err := s.stream.Log(logs)
	stream := s.mock.GetStream("group", "stream")
	token := aws.StringValue(s.stream.Token)

	c.Assert(err, IsNil)
	c.Assert(token, Equals, "1")
	c.Assert(stream.LogCount, Equals, 1)
}

func (s *LogStreamSuite) TestPutLogsToExistingStream(c *C) {
	s.mock.AddStream("group", "stream")

	logs := []*cloudwatchlogs.InputLogEvent{
		&cloudwatchlogs.InputLogEvent{
			Message:   aws.String("body"),
			Timestamp: aws.Int64(0),
		},
	}
	s.stream.Log(logs)

	err := s.stream.Log(logs)
	stream := s.mock.GetStream("group", "stream")
	token := aws.StringValue(s.stream.Token)

	c.Assert(err, IsNil)
	c.Assert(token, Equals, "2")
	c.Assert(stream.LogCount, Equals, 2)
}
