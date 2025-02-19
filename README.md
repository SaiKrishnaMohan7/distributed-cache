# Distributed Cache

Is an caching implementation for a distributed system

## Goal

- Learn low level distributed systems
- Learn Go

## Features to be implemented (Project plan)

- [ ] Leader election
  - [ ] Raft or Paxos
- [ ] Custom Protocol
  - TCP: for Gossip node discovery
  - UDP (QUIC): for node to node communication
- [ ] Architecture
  - ~Client - Server ?~
  - [ ] Peer to Peer
    - [ ] Each node with a XOR distance based routing table
- ~Consistent Hashing~
- ~TCP based comms between nodes~
  - [ ] gRPC trial
- [ ] Persistence (maybe)
  - [ ] A hybrid between AOL (append only log) and snapshot so that log file compaction we don't need to worry about

- [ ] Logging
- ~Centralized Error Handling~
- [ ] Telemetry (maybe)
- [ ] Metrics (maybe)
- [ ] Tests
  - [ ] Scale tests

## Motivations

First personal project repo in a long time. I let the pandemic destroy me. This led to a massive burnout, depression, weight gain, drinking problem... Past 5 years have been a war with self that I have persevered through. Thanks to my wife and my family. I had lost all discipline to do anything. I did not want to code at all, escapism and just so so mentally and physically exhuasted and not being able to get out of bed. This note is a thank you to life and to signal, I AM BACK IN THE GAME.
