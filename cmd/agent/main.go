package main

import (
	"github.com/Boyul-Kim/pulsesentinel/internal/agent"
)

func main() {
	paths := []string{
		"/var/log/syslog",
		"/var/log/auth.log",
		"/var/log/kern.log",
	}
	agent.Watch(paths)
}
