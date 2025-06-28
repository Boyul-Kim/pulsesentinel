package agent

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	pb "github.com/Boyul-Kim/pulsesentinel/proto/sentinel"
	"github.com/hpcloud/tail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Agent struct {
	ECSEvents  map[string][]ECSEvent
	GrpcClient pb.EventIngestorClient
}

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

func InitAgent(addr string) *Agent {
	fmt.Println("STARTING CLIENT SERVER")

	// TODO:
	// -add batching/retry logicaddr
	// -context with timeouts

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	client := pb.NewEventIngestorClient(conn)
	return &Agent{
		ECSEvents:  make(map[string][]ECSEvent),
		GrpcClient: client,
	}
}

func (a *Agent) Watch(path string) {
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
		go a.sendECSEvents(ecsEvent)
		a.ECSEvents[msgID] = append(a.ECSEvents[msgID], ecsEvent)
		delete(eventGroups, msgID)
	}
}

func (a *Agent) sendECSEvents(event ECSEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	protoEvent := toProtoSecurityEvents(event)
	resp, err := a.GrpcClient.SendEvent(ctx, protoEvent)
	if err != nil {
		log.Printf("Error sending event: %v", err)
		return
	}

	log.Printf("Ingestor response: %s", resp.Message)

}

func toProtoSecurityEvents(ecs ECSEvent) *pb.SecurityEvent {
	var ts *timestamppb.Timestamp
	if ecs.Timestamp != "" {
		// RFC3339 format
		t, err := time.Parse(time.RFC3339, ecs.Timestamp)
		if err == nil {
			ts = timestamppb.New(t)
		} else {
			ts = timestamppb.Now()
		}
	} else {
		ts = timestamppb.Now()
	}

	return &pb.SecurityEvent{
		EventId:   getStr(ecs.Event, "id"),
		Timestamp: ts,
		Event: &pb.EventMeta{
			Category: getStrSlice(ecs.Event, "category"),
			Type:     getStrSlice(ecs.Event, "type"),
			Action:   getStr(ecs.Event, "action"),
			Outcome:  getStr(ecs.Event, "outcome"),
			Provider: getStr(ecs.Event, "provider"),
		},
		Host: &pb.HostMeta{
			Hostname: getStr(ecs.Host, "hostname"),
			Id:       getStr(ecs.Host, "id"),
		},
		User: &pb.UserMeta{
			Name:        getStr(ecs.User, "name"),
			Id:          getStr(ecs.User, "id"),
			EffectiveId: getStr(ecs.User, "effective_id"),
		},
		Source: &pb.SourceMeta{
			Ip:   getStr(ecs.Source, "ip"),
			Port: getInt32(ecs.Source, "port"),
		},
		Destination: &pb.DestinationMeta{
			Ip:   getStr(ecs.Destination, "ip"),
			Port: getInt32(ecs.Destination, "port"),
		},
		Process: &pb.ProcessMeta{
			Name:       getStr(ecs.Process, "name"),
			Executable: getStr(ecs.Process, "executable"),
			Pid:        getInt32(ecs.Process, "pid"),
			Ppid:       getInt32(ecs.Process, "ppid"),
			Args:       getStrSlice(ecs.Process, "args"),
		},
		File: &pb.FileMeta{
			Path:   getStr(ecs.File, "path"),
			Access: getStr(ecs.File, "access"),
		},
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
	innerKVRe := regexp.MustCompile(`(\w+)=("[^"]*"|'[^']*'|[^\s]+)`)
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
			if key == "msg" && strings.HasPrefix(val, "op=") {
				for _, inner := range innerKVRe.FindAllStringSubmatch(val, -1) {
					innerKey := inner[1]
					innerVal := inner[2]
					innerVal = strings.Trim(innerVal, `"'`)
					fields[innerKey] = innerVal
				}
			} else {
				fields[key] = val
			}
		}
	}
	return fields
}

func getStr(m map[string]interface{}, k string) string {
	if v, ok := m[k]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getStrSlice(m map[string]interface{}, k string) []string {
	if v, ok := m[k]; ok && v != nil {
		if arr, ok := v.([]string); ok {
			return arr
		} else if s, ok := v.(string); ok {
			return []string{s}
		}
	}
	return nil
}

func getInt32(m map[string]interface{}, k string) int32 {
	if v, ok := m[k]; ok && v != nil {
		switch vv := v.(type) {
		case int:
			return int32(vv)
		case int32:
			return vv
		case int64:
			return int32(vv)
		case float64:
			return int32(vv)
		case string:
			n, _ := strconv.Atoi(vv)
			return int32(n)
		}
	}
	return 0
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
		ecs.Timestamp = ts // RFC3339 format
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
