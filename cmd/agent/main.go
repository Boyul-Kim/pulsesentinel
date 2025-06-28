package main

import (
	"github.com/Boyul-Kim/pulsesentinel/internal/agent"
)

func main() {
	println("STARTING AGENT")
	//Unfortunately, I am developing on WSL2 Ubuntu and cannot use auditd to access the audit logs. Using a simulator to simulate audit logs for now
	agentService := agent.InitAgent("localhost:50051")
	agentService.Watch("/tmp/simulation_agent.log")
}
