# ros-iface-streaming

**ros-iface-streaming** is a microservice that connects to Mikrotik (RouterOS) devices and streams real-time traffic data from a specific network interface using sockets

# Initial tools
- go 1.24
- taskfile `go install github.com/air-verse/air@latest`

# Env
.env file example:
```ssh
APP_NAME="ros-iface-streaming"
APP_ENV="dev"

HTTP_URL="localhost"
HTTP_PORT="8080"
HTTP_ALLOWED_ORIGINS="localhost:8080,"

ROUTEROS_PORT="8728"
ROUTEROS_USER=""
ROUTEROS_PASS=""
```

# Usage
- task dev

# Response
- {"iface": "ether1", "rx": 123456, "tx": 123456}

# TODO
- [ ] add metrics
