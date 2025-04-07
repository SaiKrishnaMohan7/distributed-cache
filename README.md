# Distributed Cache

Is an caching implementation for a distributed system

## Goal

- Learn low level distributed systems
- Learn Go

## Features to be implemented (Project plan)

- ~~[ ] Leader election~~
  - ~~[ ] Raft or Paxos~~
- [ ] Custom Protocol
  - TCP: for Gossip node discovery
  - UDP (QUIC): for node to node communication
- [ ] Architecture
  - ~~Client - Server ?~~
  - [ ] Peer to Peer
    - [ ] Each node with a XOR distance based routing table
- ~~Consistent Hashing~~
- ~~TCP based comms between nodes~~
  - [ ] gRPC trial
- [ ] Persistence (maybe)
  - [ ] A hybrid between AOL (append only log) and snapshot so that log file compaction we don't need to worry about
- [x] Logging
- ~~Centralized Error Handling~~
- [ ] Telemetry (maybe)
- [ ] Metrics (maybe)
- [ ] Tests
  - [ ] Scale tests
