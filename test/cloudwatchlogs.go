package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// MockStream stores the state of a fake stream.
type MockStream struct {
	LogCount int
	Token    int
}

// CloudWatchLogsMock mocks the CloudFront Logs API.
type CloudWatchLogsMock struct {
	*httptest.Server

	Groups  map[string]map[string]*MockStream
	Streams []*cloudwatchlogs.LogStream
}

// NewCloudWatchLogsMock instantiates a mock CloudFront Logs server.
func NewCloudWatchLogsMock() *CloudWatchLogsMock {
	mock := &CloudWatchLogsMock{}
	mock.Server = httptest.NewServer(mock)
	mock.Streams = []*cloudwatchlogs.LogStream{}
	mock.Groups = map[string]map[string]*MockStream{}

	return mock
}

// AddStream registers a new stream.
func (m *CloudWatchLogsMock) AddStream(group, stream string) {
	g, ok := m.Groups[group]
	if !ok {
		g = map[string]*MockStream{}
		m.Groups[group] = g
	}

	s, ok := g[stream]
	if !ok {
		s = &MockStream{}
		g[stream] = s
	}
}

// GetStreams returns a list of streams in a group.
func (m *CloudWatchLogsMock) GetStreams(group string) map[string]*MockStream {
	return m.Groups[group]
}

// GetStream returns a stream in a group.
func (m *CloudWatchLogsMock) GetStream(group, stream string) *MockStream {
	if group, ok := m.Groups[group]; ok {
		return group[stream]
	}

	return nil
}

func (m *CloudWatchLogsMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.Header["X-Amz-Target"][0]

	switch action {
	case "Logs_20140328.DescribeLogStreams":
		m.describeLogStreams(w, r)
	case "Logs_20140328.CreateLogStream":
		m.createLogStream(w, r)
	case "Logs_20140328.PutLogEvents":
		m.putLogEvents(w, r)
	}
}

func (m *CloudWatchLogsMock) describeLogStreams(w http.ResponseWriter, r *http.Request) {
	data := &cloudwatchlogs.DescribeLogStreamsInput{}
	m.readJSON(r.Body, data)

	group := aws.StringValue(data.LogGroupName)
	streams := []interface{}{}

	for n, s := range m.Groups[group] {
		streams = append(streams, map[string]string{
			"logStreamName":       n,
			"uploadSequenceToken": strconv.Itoa(s.Token),
		})
	}

	m.writeJSON(w, &map[string]interface{}{
		"logStreams": streams,
	})
}

func (m *CloudWatchLogsMock) createLogStream(w http.ResponseWriter, r *http.Request) {
	data := &cloudwatchlogs.CreateLogStreamInput{}
	m.readJSON(r.Body, data)

	group := aws.StringValue(data.LogGroupName)
	stream := aws.StringValue(data.LogStreamName)

	m.AddStream(group, stream)
	m.writeJSON(w, &map[string]interface{}{})
}

func (m *CloudWatchLogsMock) putLogEvents(w http.ResponseWriter, r *http.Request) {
	data := &cloudwatchlogs.PutLogEventsInput{}
	m.readJSON(r.Body, data)

	group := aws.StringValue(data.LogGroupName)
	stream := aws.StringValue(data.LogStreamName)
	token := aws.StringValue(data.SequenceToken)

	s := m.GetStream(group, stream)
	if strconv.Itoa(s.Token) != token && s.Token != 0 {
		w.WriteHeader(http.StatusBadRequest)
		m.writeJSON(w, &map[string]interface{}{
			"__type":  "InvalidSequenceTokenException",
			"message": "The given sequenceToken is invalid.",
		})
		return
	}

	s.LogCount += len(data.LogEvents)
	s.Token++

	m.writeJSON(w, &map[string]interface{}{
		"nextSequenceToken": strconv.Itoa(s.Token),
	})
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
