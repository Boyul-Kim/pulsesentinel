package simulator

import (
	"fmt"
	"os"
	"time"
)

func GenerateRawAuditEvents() {
	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				writeBenignAuditLogEvent()
				writeSuspiciousAuditLogEvent()
			}
		}
	}()

	time.Sleep(5 * time.Second)
	ticker.Stop()
	done <- true
	println("SIMULATION DONE")
}

// Writes a benign login event to /tmp/simulation_agent.log in audit.log style
func writeBenignAuditLogEvent() error {
	now := time.Now()
	epoch := now.Unix()
	eventID := 1003 // increment or randomize as needed

	msgID := fmt.Sprintf("audit(%d.000:%d)", epoch, eventID)
	lines := []string{
		fmt.Sprintf("type=USER_AUTH msg=%s: pid=24567 uid=1000 auid=1000 ses=3 msg='op=PAM:authentication grantors=pam_unix acct=\"alice\" exe=\"/usr/sbin/sshd\" hostname=? addr=192.168.10.15 terminal=ssh res=success'", msgID),
		fmt.Sprintf("type=USER_ACCT msg=%s: pid=24567 uid=1000 auid=1000 ses=3 msg='op=PAM:accounting grantors=pam_unix acct=\"alice\" exe=\"/usr/sbin/sshd\" hostname=? addr=192.168.10.15 terminal=ssh res=success'", msgID),
		fmt.Sprintf("type=USER_LOGIN msg=%s: pid=24567 uid=1000 auid=1000 ses=3 msg='op=login acct=\"alice\" exe=\"/usr/sbin/sshd\" hostname=? addr=192.168.10.15 terminal=ssh res=success'", msgID),
	}
	return appendLinesToFile("/tmp/simulation_agent.log", lines)
}

// Writes a suspicious failed login event to /tmp/simulation_agent.log in audit.log style
func writeSuspiciousAuditLogEvent() error {
	now := time.Now()
	epoch := now.Unix()
	eventID := 1002 // increment or randomize as you wish

	msgID := fmt.Sprintf("audit(%d.000:%d)", epoch, eventID)
	lines := []string{
		fmt.Sprintf("type=USER_AUTH msg=%s: pid=23212 uid=0 auid=4294967295 ses=4294967295 msg='op=PAM:authentication grantors=test acct=\"root\" exe=\"/usr/sbin/sshd\" hostname=test addr=203.0.113.77 terminal=ssh res=failed'", msgID),
		fmt.Sprintf("type=USER_LOGIN msg=%s: pid=23212 uid=0 auid=4294967295 ses=4294967295 msg='op=login acct=\"root\" exe=\"/usr/sbin/sshd\" hostname=test addr=203.0.113.77 terminal=ssh res=failed'", msgID),
	}
	return appendLinesToFile("/tmp/simulation_agent.log", lines)
}

func appendLinesToFile(path string, lines []string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}
