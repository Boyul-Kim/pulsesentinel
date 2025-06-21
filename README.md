# pulsesentinel

PulseSentinel is a high-performance, real-time security event pipeline and monitoring platform.
Collect, process, and analyze security logs from distributed agents with scalable microservices built in Go, Kafka, PostgreSQL, and Elasticsearch.
Features (will) include real-time event ingestion, intelligent alerting, resilient batch processing, and flexible search APIs for security analytics and dashboards.

## agent
To use the agent on a linux host, it will need auditd

sudo apt update
sudo apt install auditd
sudo service auditd start

NOTE:
I am currently using WSL2 Ubuntu. There are limitations with running auditd so I can't ready off "/var/log/audit/audit.log".
Instead, I am using a simulator to mock events in the ECS format.

