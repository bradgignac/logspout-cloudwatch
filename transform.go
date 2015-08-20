package cloudwatch

import "github.com/gliderlabs/logspout/router"

func transform(messages <-chan *router.Message) <-chan Log {
	logs := make(chan Log)

	go func() {
		defer close(logs)

		for msg := range messages {
			logs <- transformMessage(msg)
		}
	}()

	return logs
}

func transformMessage(msg *router.Message) Log {
	return &LogMessage{msg}
}
