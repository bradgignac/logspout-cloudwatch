package cloudwatch

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func batch(logs <-chan Log, capacity Capacity) <-chan []Log {
	batches := make(chan []Log)
	batcher := NewBatcher(logs, batches, capacity)

	go func() {
		defer close(batches)
		batcher.Start()
	}()

	return batches
}

// Batcher buffers messages on a channel until a flush is triggered.
type Batcher struct {
	in  <-chan Log
	out chan<- []Log

	messages []Log
	capacity Capacity
	timer    <-chan time.Time
	size     int
}

// NewBatcher creates a Batcher that buffers message from the input channel to the
// output channel.
func NewBatcher(in <-chan Log, out chan<- []Log, capacity Capacity) *Batcher {
	return &Batcher{in: in, out: out, capacity: capacity}
}

// Length returns the current length of the buffer.
func (b *Batcher) Length() int {
	return len(b.messages)
}

// Start begins buffering messages from the input channel.
func (b *Batcher) Start() {
loop:
	for {
		select {
		case l, ok := <-b.in:
			if !ok {
				break loop
			}

			if b.willOverflow(l) {
				log.Debugf("Batch flushed to prevent size overflow - size: %d, capacity: %v", b.size, b.capacity)
				b.flush()
			}

			b.messages = append(b.messages, l)
			b.size += l.Size()

			if b.isFullSize() {
				log.Debugf("Batch flushed due to batch size - size: %d, capacity: %v", b.size, b.capacity)
				b.flush()
			} else if b.isFullLength() {
				log.Debugf("Batch flushed due to batch length - length: %d, capacity: %v", len(b.messages), b.capacity)
				b.flush()
			} else {
				b.startFlushTimer()
			}
		case <-b.timer:
			log.Debugf("Batch flushed due to timer - capacity: %v", b.capacity)
			b.flush()
		}
	}

	b.flush()
}

func (b *Batcher) willOverflow(log Log) bool {
	return b.size+log.Size() > b.capacity.Size
}

func (b *Batcher) isFullSize() bool {
	return b.size >= b.capacity.Size
}

func (b *Batcher) isFullLength() bool {
	return len(b.messages) == b.capacity.Length
}

func (b *Batcher) flush() {
	messages := make([]Log, len(b.messages))
	copy(messages, b.messages)

	b.timer = nil
	b.messages = nil
	b.size = 0

	if len(messages) != 0 {
		b.out <- messages
	}
}

func (b *Batcher) startFlushTimer() {
	if b.timer == nil && b.capacity.Duration > 0 {
		b.timer = time.After(b.capacity.Duration)
	}
}

// Capacity returns conditions that trigger a Batcher to flush.
type Capacity struct {
	Size     int
	Length   int
	Duration time.Duration
}
