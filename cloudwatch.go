package cloudwatch

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gliderlabs/logspout/router"
)

// Set batch sizes based CloudWatch Logs limits found in the developer guide.
// While the actual size limit is 1 MB, we use 900 KB due to small differences
// in how batch size is calculated.
const batchSize = 900000
const batchLength = 10000
const batchDuration = 250 * time.Millisecond

func init() {
	router.AdapterFactories.Register(NewAdapter, "cloudwatch")
}

// Adapter ships logs to AWS CloudWatch.
type Adapter struct {
	route     *router.Route
	logstream *LogStream
	capacity  Capacity
}

// NewAdapter instances a new AWS CloudWatch adapter.
func NewAdapter(route *router.Route) (router.LogAdapter, error) {
	group := os.Getenv("AWS_LOG_GROUP")
	stream := os.Getenv("AWS_LOG_STREAM")
	logstream, err := NewLogStream(group, stream)
	if err != nil {
		return nil, err
	}

	capacity := Capacity{
		Size:     batchSize,
		Length:   batchLength,
		Duration: batchDuration,
	}

	log.Printf("Created CloudWatch adapter - group: %s, stream: %s", group, stream)

	return &Adapter{route: route, logstream: logstream, capacity: capacity}, nil
}

// Stream passes messages from a logspout message channel to AWS CloudWatch.
func (a *Adapter) Stream(logstream chan *router.Message) {
	log.Printf("CloudWatch adapter is streaming Docker logs")

	logs := transform(logstream)
	batches := batch(logs, a.capacity)

	for batch := range batches {
		events := make([]*cloudwatchlogs.InputLogEvent, len(batch))

		for i, log := range batch {
			events[i] = &cloudwatchlogs.InputLogEvent{
				Message:   aws.String(log.Body()),
				Timestamp: aws.Int64(log.Timestamp()),
			}
		}

		a.logstream.Log(events)
	}
}