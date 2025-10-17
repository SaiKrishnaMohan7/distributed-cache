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

## TTL and Cleanup Behavior

### How TTL Works

- When you `Set()` a key with a TTL, it expires after that duration
- `Get()` will return an error for expired keys immediately (strict guarantee)
- Background cleanup removes expired keys from memory periodically

### Configuring Cleanup

Set `CACHE_CLEANUP_TICK` environment variable to control cleanup frequency:

```bash
CACHE_CLEANUP_TICK=1000  # Run cleanup every 1 second (1000ms)
```

**Tuning Guidelines:**

- **High TTL precision needs**: 100-500ms (e.g., sub-second TTLs)
- **Balanced**: 1000-5000ms (default recommendation)
- **Low memory pressure**: 10000-60000ms (cleanup less often)

**Note:** Cleanup timing only affects when memory is freed, not when keys become inaccessible.
TTL is always strictly enforced by `Get()`.
