# UDP-chat

The chat server made for you :)

It consists in two parts, the [client](client/README.md) and the [server](server/README.md).

## Roadmap

1. General
   - [x] Define basic architecture
   - [x] Create Redis infrastructure ([1](https://aws.amazon.com/pt/elasticache/redis/]))
2. Server
   - [x] Create base project
     - [x] Configure logging
   - [x] Receive messages from client
   - [x] Broadcast messages to clients
   - [x] Save messages on Redis
     - [x] Save connected clients
   - [x] Read messages from Redis
   - [ ] Delete messages
   - [ ] Flush DB when no one is connected to chat server anymore
   - [ ] Unit testing
   - [ ] Integration testing
3. Client
   - [ ] Create base project
     - [ ] Create graphical interface
   - [ ] Create BFF
   - [ ] Connect to chat server
   - [ ] Read messages from Redis
   - [ ] Send messages to server
   - [ ] Delete messages
   - [ ] Unit testing
   - [ ] Integration testing
4. Documentation
   - [ ] Requirements
   - [ ] How to build and run
   - [ ] How to run tests
   - [ ] Architecture diagrams
