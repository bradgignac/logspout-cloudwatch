package cloudwatch

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// LogStream ships logs to AWS CloudWatch.
type LogStream struct {
	Group   *string
	Stream  *string
	Token   *string
	service *cloudwatchlogs.CloudWatchLogs
}

// NewLogStream instantiates a Logger.
func NewLogStream(group, stream string, config client.ConfigProvider) *LogStream {
	cloudwatch := cloudwatchlogs.New(config)
	logstream := &LogStream{
		Group:   aws.String(group),
		Stream:  aws.String(stream),
		service: cloudwatch,
	}

	return logstream
}

// Init fetches the sequence token for a stream so logs can be streamed.
func (s *LogStream) Init() error {
	stream, err := s.findStream()
	if err != nil {
		return err
	}

	if stream != nil {
		s.Token = stream.UploadSequenceToken
		return nil
	}

	return s.createStream()
}

func (s *LogStream) createStream() error {
	params := &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  s.Group,
		LogStreamName: s.Stream,
	}

	_, err := s.service.CreateLogStream(params)

	return err
}

func (s *LogStream) findStream() (*cloudwatchlogs.LogStream, error) {
	params := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        s.Group,
		LogStreamNamePrefix: s.Stream,
		Limit:               aws.Int64(1),
	}

	resp, err := s.service.DescribeLogStreams(params)
	if err != nil {
		return nil, err
	}

	if len(resp.LogStreams) == 0 {
		return nil, nil
	}

	return resp.LogStreams[0], nil
}

// Log submits a batch of logs to the LogStream.
func (s *LogStream) Log(logs []*cloudwatchlogs.InputLogEvent) error {
	params := &cloudwatchlogs.PutLogEventsInput{
		LogEvents:     logs,
		LogGroupName:  s.Group,
		LogStreamName: s.Stream,
		SequenceToken: s.Token,
	}

	resp, err := s.service.PutLogEvents(params)
	awserr, _ := err.(awserr.Error)

	if awserr != nil {
		switch awserr.Code() {
		case "InvalidSequenceTokenException":
			log.Infof("Retrying log upload with new token - length %d, error, %v", len(logs), err)
			return s.retryBatchWithNewToken(logs)
		default:
			log.Errorf("Log upload failed - length: %d, error: %v", len(logs), err)
			return awserr
		}
	}

	if resp.RejectedLogEventsInfo != nil {
		log.Warnf("Log upload succeeded with rejected events - length: %d", len(logs))
	} else {
		log.Debugf("Log upload succeeded - length: %d", len(logs))
	}

	s.Token = resp.NextSequenceToken

	return nil
}

func (s *LogStream) retryBatchWithNewToken(logs []*cloudwatchlogs.InputLogEvent) error {
	if err := s.Init(); err != nil {
		return err
	}

	return s.Log(logs)
}
