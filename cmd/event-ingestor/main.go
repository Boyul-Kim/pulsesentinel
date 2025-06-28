package main

import (
	"log"
	"net"

	event_ingestor "github.com/Boyul-Kim/pulsesentinel/internal/event-ingestor"

	pb "github.com/Boyul-Kim/pulsesentinel/proto/sentinel"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterEventIngestorServer(s, &event_ingestor.EventIngestorServer{})
	log.Println("EVENT INGESTOR SERVICE LISTENING ON :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
