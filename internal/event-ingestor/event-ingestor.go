package event_ingestor

import (
	"context"
	"log"

	pb "github.com/Boyul-Kim/pulsesentinel/proto/sentinel"
)

type EventIngestorServer struct {
	pb.UnimplementedEventIngestorServer
}

func (s *EventIngestorServer) SendEvent(ctx context.Context, event *pb.SecurityEvent) (*pb.EventResponse, error) {
	log.Printf("Received event: %v", event)
	return &pb.EventResponse{Message: "Event received"}, nil
}
