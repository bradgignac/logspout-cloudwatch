package cloudwatch

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/gliderlabs/logspout/router"
)

const NumMessages = 250000

func TestCloudWatchAdapter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode.")
	}

	route := &router.Route{Address: "logspout-cloudwatch"}
	messages := make(chan *router.Message)

	adapter, err := NewAdapter(route)
	if err != nil {
		t.Error(err)
		return
	}

	go adapter.Stream(messages)
	for i := 0; i < NumMessages; i++ {
		messages <- createMessage()
	}

	close(messages)
}

func createMessage() *router.Message {
	data := ""
	timestamp := time.Now()
	random := rand.Intn(100)

	if random != 0 {
		data = randomdata.Paragraph()
	}

	return &router.Message{Data: data, Time: timestamp}
}
