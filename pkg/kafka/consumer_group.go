package kafka

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

type consumerGroup struct {
	Brokers []string
	GroupID string
	log     *logrus.Logger
}

// NewConsumerGroup kafka consumer group constructor
func NewConsumerGroup(brokers []string, groupID string, log *logrus.Logger) *consumerGroup {
	return &consumerGroup{Brokers: brokers, GroupID: groupID, log: log}
}

// ConsumeTopic start consumer group with given worker and pool size
func (c *consumerGroup) ConsumeTopic(ctx context.Context, groupTopics []string, poolSize int, worker Worker) {
	r := NewReader(c.Brokers, groupTopics, c.GroupID)

	defer func() {
		if err := r.Close(); err != nil {
			c.log.Warnf("consumerGroup.r.Close: %v", err)
		}
	}()

	c.log.Infof("Starting consumer groupID: %s, topic: %+v, pool size: %v", c.GroupID, groupTopics, poolSize)

	wg := &sync.WaitGroup{}
	for i := 0; i <= poolSize; i++ {
		wg.Add(1)
		go worker(ctx, r, wg, i)
	}
	wg.Wait()
}
