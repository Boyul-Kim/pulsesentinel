package agent

import (
	"fmt"
	"time"

	"math/rand"
)

type ECSEvent struct {
	Timestamp   string                 `json:"@timestamp"`
	Event       map[string]interface{} `json:"event"`
	Host        map[string]interface{} `json:"host"`
	User        map[string]interface{} `json:"user,omitempty"`
	Source      map[string]interface{} `json:"source,omitempty"`
	Destination map[string]interface{} `json:"destination,omitempty"`
	Process     map[string]interface{} `json:"process,omitempty"`
	File        map[string]interface{} `json:"file,omitempty"`
	Related     map[string]interface{} `json:"related,omitempty"`
}

func generateECSLoginEvent(agentID, host, user, ip, outcome string) ECSEvent {
	return ECSEvent{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Event: map[string]interface{}{
			"id":       fmt.Sprintf("%d", rand.Int63()),
			"category": []string{"authentication"},
			"type":     []string{"start"},
			"outcome":  outcome, // "success" or "failure"
			"action":   "user_login",
		},
		Host: map[string]interface{}{
			"hostname": host,
			"id":       agentID,
		},
		User: map[string]interface{}{
			"name": user,
		},
		Source: map[string]interface{}{
			"ip": ip,
		},
	}
}
