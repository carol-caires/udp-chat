# UDP-chat Server

UDP server that receives messages from clients and broadcasts it to all connected clients.

This project follows [this standards](https://github.com/golang-standards/project-layout) for internal structure better understanding.

## What do you need to run it
- To configure a local Redis cluster to save clients and messages

## Environment variables

```
# General service configuration
HOST=127.0.0.1 # mandatory
PORT=1337 # mandatory
BLOCKING_DEADLINE_SECONDS=15
MAX_BUFFER_SIZE_BYTES=1024

# Redis configuration
REDIS_ADDR=localhost:6379 # mandatory
```

## How to run it

## How to test it