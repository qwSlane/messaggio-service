package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"sync"
)

// Processor methods must implement kafka.Worker func method interface
type Processor interface {
	ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)
}

// Worker kafka consumer worker fetch and process messages from reader
type Worker func(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)
