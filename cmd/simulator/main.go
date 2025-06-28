package main

import (
	"log"
	"os"

	"github.com/Boyul-Kim/pulsesentinel/internal/simulator"
)

func main() {
	println("STARTING SIMULATOR")
	if err := os.Truncate("/tmp/simulation_agent.log", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	simulator.GenerateRawAuditEvents()
}
