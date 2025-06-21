package main

import (
	"github.com/Boyul-Kim/pulsesentinel/internal/agent"
)

func main() {
	paths := []string{

		//Unfortunately, I am developing on WSL2 Ubuntu and cannot use auditd to access the audit logs. Using a simulator to simulate audit logs for now
		// "/var/log/syslog",
		// "/var/log/auth.log",
		// "/var/log/kern.log",
		//"/var/log/audit/audit.log",
	}
	agent.Watch(paths)
}
