package cloudwatch

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/bradgignac/logspout-cloudwatch/test"
	. "gopkg.in/check.v1"
)

func TestLogStream(t *testing.T) {
	TestingT(t)
}

type LogStreamSuite struct {
	mock    *test.CloudWatchLogsMock
	session *session.Session
}

var _ = Suite(&LogStreamSuite{})

func (s *LogStreamSuite) SetUpTest(c *C) {
	s.mock = test.NewCloudWatchLogsMock()

	creds := credentials.NewStaticCredentials("id", "secret", "token")
	config := aws.NewConfig().
		WithCredentials(creds).
		WithEndpoint(s.mock.URL).
		WithRegion("us-east-1").
		WithDisableSSL(true)
	s.session = session.New(config)
}

func (s *LogStreamSuite) TearDownTest(c *C) {
	s.mock.Close()
}

func (s *LogStreamSuite) TestNewStream(c *C) {
	stream := NewLogStream("group", "new", s.session)
	err := stream.Init()

	c.Assert(err, IsNil)
	c.Assert(stream.Token, IsNil)
	c.Assert(s.mock.Streams, HasLen, 1)
}

func (s *LogStreamSuite) TestExistingStream(c *C) {
	s.mock.Streams = append(s.mock.Streams, &cloudwatchlogs.LogStream{
		LogStreamName:       aws.String("existing"),
		UploadSequenceToken: aws.String("existing"),
	})

	stream := NewLogStream("group", "existing", s.session)
	err := stream.Init()

	c.Assert(err, IsNil)
	c.Assert(stream.Token, NotNil)
	c.Assert(s.mock.Streams, HasLen, 1)
}
