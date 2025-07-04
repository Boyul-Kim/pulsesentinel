package main

import (
	"github.com/Boyul-Kim/pulsesentinel/internal/monitor"
	"github.com/segmentio/kafka-go"
)

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "alerts",
		GroupID:  "monitoring",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	monitor := &monitor.MonitorServer{KafkaConsumer: r}
	monitor.ConsumeEvents()
}
