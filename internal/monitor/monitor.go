package monitor

import (
	"context"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type MonitorServer struct {
	KafkaConsumer *kafka.Reader
}

func (s *MonitorServer) ConsumeEvents() {
	log.Printf("Starting monitor service...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer s.KafkaConsumer.Close()

	const (
		workers      = 8
		channelDepth = 32 // Managing backpressure
	)
	var wg sync.WaitGroup
	msgChan := make(chan kafka.Message, channelDepth)

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func(workerId int) {
			defer wg.Done()
			for m := range msgChan {
				s.processEvents(m, workerId)

				if err := s.KafkaConsumer.CommitMessages(ctx, m); err != nil {
					log.Printf("Error commting message: %v", err)
				}
			}
		}(i)
	}

	for {
		m, err := s.KafkaConsumer.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("Context error: %v", ctx.Err())
				break
			}

			log.Printf("Error fetching message: %v", err)
			continue
		}

		msgChan <- m
	}

	close(msgChan)
	wg.Wait()
	log.Printf("shutting down...")
}

func (s *MonitorServer) processEvents(m kafka.Message, workerId int) {
	log.Printf("processing events from worker %d: %v", workerId, m)
}
