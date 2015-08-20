package cloudwatch

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/gliderlabs/logspout/router"
)

const NumMessages = 2000000

func generateBatch(size int) []*cloudwatchlogs.InputLogEvent {
	batch := make([]*cloudwatchlogs.InputLogEvent, size)

	for i := 0; i < size; i++ {
		msg := randomdata.Paragraph()
		json := fmt.Sprintf("{ \"message\": \"%s\", \"category\": \"logspout\" }", msg)
		now := time.Now().Unix() * 1000

		batch[i] = &cloudwatchlogs.InputLogEvent{
			Message:   aws.String(json),
			Timestamp: aws.Int64(now),
		}
	}

	return batch
}

func TestCloudWatchAdapter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	os.Setenv("AWS_LOG_GROUP", "logspout-cloudwatch")
	os.Setenv("AWS_LOG_STREAM", "integration")

	route := &router.Route{}
	messages := make(chan *router.Message)

	adapter, err := NewAdapter(route)
	if err != nil {
		t.Error(err)
		return
	}

	go adapter.Stream(messages)
	for i := 0; i < NumMessages; i++ {
		messages <- &router.Message{Data: randomdata.Paragraph(), Time: time.Now()}
	}

	close(messages)
}
