# UDP-chat Server

UDP server that receives messages from clients and broadcasts it to all connected clients.

This project follows [this standards](https://github.com/golang-standards/project-layout) for internal structure better understanding.

## What do you need to run it
- To configure a local Redis cluster to save clients and messages

It currently only supports Redis. If do you want to support other cache providers, create a new implementation
for every new provider to `cache.Client` interface in `server/internal/infrastructure/cache`. No need of changes 
in service layer, make sure only to change the line ` cache.NewRedisConn()` to the new implementation constructor in main.

## Environment variables

```shell
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

You can send messages to this server manually, for example:

```shell
# create new client
echo '{"type":"NEW_CLIENT", "from": {"name": "Carol"}}' | nc -w1 -u 127.0.0.1 1337

# broadcast new message to clients
echo '{"body":"olá","type":"NEW_MESSAGE", "from": {"name": "Carol"}}' | nc -w1 -u 127.0.0.1 1337
```