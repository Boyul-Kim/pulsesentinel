package agent

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hpcloud/tail"
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
}

func Watch(path string) {
	eventGroups := make(map[string][]string)
	t, err := tail.TailFile(path, tail.Config{Follow: true})
	if err != nil {
		log.Fatalf("Error with tail file: %s", err)
	}

	for line := range t.Lines {
		msgID := extractMsgID(line.Text)
		eventGroups[msgID] = append(eventGroups[msgID], line.Text)
		auditEvent := parseAuditLinesToMap(eventGroups[msgID])
		ecsEvent := mapAuditToECS(auditEvent)

		fmt.Println("ECS EVENT", ecsEvent)
		// outputECSJSON(ecsEvent)
		// delete(eventGroups, msgID)
	}
}

func extractMsgID(line string) string {
	re := regexp.MustCompile(`msg=audit\(([^)]+)\)`)
	match := re.FindStringSubmatch(line)

	if len(match) > 1 {
		return match[1]
	}

	return ""
}

func parseAuditLinesToMap(lines []string) map[string]string {
	fields := make(map[string]string)
	kvRe := regexp.MustCompile(`(\w+)=("[^"]*"|'[^']*'|[^\s]+)`)
	for _, line := range lines {
		if idx := strings.Index(line, "type="); idx != -1 {
			parts := strings.SplitN(line[idx:], " ", 2)
			typeKV := strings.SplitN(parts[0], "=", 2)
			if len(typeKV) == 2 {
				fields["event_type"] = typeKV[1]
			}
		}
		for _, match := range kvRe.FindAllStringSubmatch(line, -1) {
			key := match[1]
			val := match[2]
			val = strings.Trim(val, `"'`)
			fields[key] = val
		}
	}
	return fields
}

func mapAuditToECS(audit map[string]string) ECSEvent {
	ecs := ECSEvent{
		Event:       make(map[string]interface{}),
		Host:        make(map[string]interface{}),
		User:        make(map[string]interface{}),
		Source:      make(map[string]interface{}),
		Destination: make(map[string]interface{}),
		Process:     make(map[string]interface{}),
		File:        make(map[string]interface{}),
	}

	if ts, ok := audit["timestamp"]; ok {
		ecs.Timestamp = ts // Should be ISO8601 or RFC3339 format
	}
	if id, ok := audit["msgid"]; ok {
		ecs.Event["id"] = id
	}
	if eventType, ok := audit["event_type"]; ok {
		ecs.Event["type"] = []string{eventType}
	}
	if category, ok := audit["event_category"]; ok {
		ecs.Event["category"] = []string{category}
	}
	if action, ok := audit["op"]; ok {
		ecs.Event["action"] = action
	}
	if outcome, ok := audit["res"]; ok {
		if outcome == "success" {
			ecs.Event["outcome"] = "success"
		} else if outcome == "failed" {
			ecs.Event["outcome"] = "failure"
		}
	}

	if host, ok := audit["hostname"]; ok {
		ecs.Host["hostname"] = host
	}
	if agentID, ok := audit["agent_id"]; ok {
		ecs.Host["id"] = agentID
	}

	if user, ok := audit["acct"]; ok {
		ecs.User["name"] = user
	}
	if uid, ok := audit["uid"]; ok {
		ecs.User["id"] = uid
	}
	if auid, ok := audit["auid"]; ok {
		ecs.User["effective_id"] = auid
	}

	if ip, ok := audit["addr"]; ok {
		ecs.Source["ip"] = ip
	}
	if srcPort, ok := audit["port"]; ok {
		ecs.Source["port"] = srcPort
	}

	if dstIP, ok := audit["dst_ip"]; ok {
		ecs.Destination["ip"] = dstIP
	}
	if dstPort, ok := audit["dst_port"]; ok {
		ecs.Destination["port"] = dstPort
	}

	if exe, ok := audit["exe"]; ok {
		ecs.Process["executable"] = exe
	}
	if comm, ok := audit["comm"]; ok {
		ecs.Process["name"] = comm
	}
	if pid, ok := audit["pid"]; ok {
		ecs.Process["pid"] = pid
	}
	if ppid, ok := audit["ppid"]; ok {
		ecs.Process["parent"] = map[string]interface{}{"pid": ppid}
	}
	if a0, ok := audit["a0"]; ok {
		ecs.Process["args"] = []string{a0}
	}

	if file, ok := audit["name"]; ok {
		ecs.File["path"] = file
	}
	if op, ok := audit["operation"]; ok {
		ecs.File["access"] = op
	}

	return ecs
}
