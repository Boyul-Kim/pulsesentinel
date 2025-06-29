package event_ingestor

import (
	"fmt"
	"io"
	"log"

	pb "github.com/Boyul-Kim/pulsesentinel/proto/sentinel"
)

type EventIngestorServer struct {
	pb.UnimplementedEventIngestorServer
}

func (s *EventIngestorServer) StreamEvents(stream pb.EventIngestor_StreamEventsServer) error {
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
