package event_ingestor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	pb "github.com/Boyul-Kim/pulsesentinel/proto/sentinel"
	"github.com/segmentio/kafka-go"
)

type EventIngestorServer struct {
	pb.UnimplementedEventIngestorServer
	KafkaWriter *kafka.Writer
}

func NewEventIngestorServer() *EventIngestorServer {
	fmt.Println("Starting kafka producer...")
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Topic:        "event-logs",
		Balancer:     &kafka.Hash{},
		BatchSize:    100,
		BatchTimeout: 10e6,
		Async:        true,
	})
	return &EventIngestorServer{KafkaWriter: writer}
}

func (s *EventIngestorServer) StreamEvents(stream pb.EventIngestor_StreamEventsServer) error {
	ctx := context.Background()
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Printf("Stream error: %v", err)
			return err
		}
		log.Printf("Received event: %v", event.Event.Outcome)

		//marshal to JSON for kafka
		//TODO: serialize into protobufs later
		value, err := json.Marshal(event)
		if err != nil {
			log.Printf("Error marshaling: %v", err)
			continue
		}

		msg := kafka.Message{
			Key:   []byte(event.EventId),
			Value: value,
		}

		if err := s.KafkaWriter.WriteMessages(ctx, msg); err != nil {
			log.Printf("Failed to produce event to Kafka: %v", err)
		}

		switch event.Event.Outcome {
		case "success":
			fmt.Println("Benign event logged")
		case "failure":
			fmt.Printf("Potential harmful event detected: %v", event)
		}

		resp := &pb.EventResponse{Message: "Received event " + event.EventId}
		if err := stream.Send(resp); err != nil {
			log.Printf("Send error: %v", err)
			return err
		}
	}
}
