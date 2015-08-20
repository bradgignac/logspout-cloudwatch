package cloudwatch

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// LogStream ships logs to AWS CloudWatch.
type LogStream struct {
	group   *string
	stream  *string
	token   *string
	service *cloudwatchlogs.CloudWatchLogs
}

// NewLogStream instantiates a Logger.
func NewLogStream(group, stream string) (*LogStream, error) {
	cloudwatch := cloudwatchlogs.New(nil)
	logstream := &LogStream{
		group:   aws.String(group),
		stream:  aws.String(stream),
		service: cloudwatch,
	}

	err := logstream.Init()
	if err != nil {
		return nil, err
	}

	return logstream, nil
}

// Init fetches the sequence token for a stream so logs can be streamed.
func (s *LogStream) Init() error {
	stream, err := s.findStream()
	if err != nil {
		return err
	}

	if stream != nil {
		s.token = stream.UploadSequenceToken
		return nil
	}

	return s.createStream()
}

func (s *LogStream) createStream() error {
	params := &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  s.group,
		LogStreamName: s.stream,
	}

	_, err := s.service.CreateLogStream(params)

	return err
}

func (s *LogStream) findStream() (*cloudwatchlogs.LogStream, error) {
	params := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        s.group,
		LogStreamNamePrefix: s.stream,
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
func (s *LogStream) Log(logs []*cloudwatchlogs.InputLogEvent) {
	params := &cloudwatchlogs.PutLogEventsInput{
		LogEvents:     logs,
		LogGroupName:  s.group,
		LogStreamName: s.stream,
		SequenceToken: s.token,
	}

	resp, err := s.service.PutLogEvents(params)
	if err != nil {
		log.Printf("Log upload failed - length: %d, error: %v", len(logs), err)
		return
	}

	if resp.RejectedLogEventsInfo != nil {
		log.Printf("Log upload succeeded with rejected events - length: %d", len(logs))
	} else {
		log.Printf("Log upload succeeded - length: %d", len(logs))
	}

	s.token = resp.NextSequenceToken
}
