package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// CloudWatchLogsMock mocks the CloudFront Logs API.
type CloudWatchLogsMock struct {
	*httptest.Server

	Streams []*cloudwatchlogs.LogStream
}

// NewCloudWatchLogsMock instantiates a mock CloudFront Logs server.
func NewCloudWatchLogsMock() *CloudWatchLogsMock {
	mock := &CloudWatchLogsMock{}
	mock.Server = httptest.NewServer(mock)
	mock.Streams = []*cloudwatchlogs.LogStream{}

	return mock
}

// AddStream registers a new stream.
func (m *CloudWatchLogsMock) AddStream(group, stream string) {
	fmt.Println("Adding stream...")
}

func (m *CloudWatchLogsMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Header["X-Amz-Target"][0]

	switch action {
	case "Logs_20140328.DescribeLogStreams":
		m.describeLogStreams(w, r)
	case "Logs_20140328.CreateLogStream":
		m.createLogStream(w, r)
	}
}

func (m *CloudWatchLogsMock) describeLogStreams(w http.ResponseWriter, r *http.Request) {
	var streams []interface{}

	for _, s := range m.Streams {
		streams = append(streams, map[string]string{
			"logStreamName":       *s.LogStreamName,
			"uploadSequenceToken": *s.UploadSequenceToken,
		})
	}

	m.writeJSON(w, &map[string]interface{}{
		"logStreams": streams,
	})
}

func (m *CloudWatchLogsMock) createLogStream(w http.ResponseWriter, r *http.Request) {
	data := &cloudwatchlogs.CreateLogStreamInput{}

	m.readJSON(r.Body, data)
	m.Streams = append(m.Streams, &cloudwatchlogs.LogStream{
		LogStreamName:       data.LogStreamName,
		UploadSequenceToken: aws.String("new stream"),
	})

	m.writeJSON(w, &map[string]interface{}{})
}

func (m *CloudWatchLogsMock) readJSON(body io.ReadCloser, data interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(data)
}

func (m *CloudWatchLogsMock) writeJSON(w http.ResponseWriter, data interface{}) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}
