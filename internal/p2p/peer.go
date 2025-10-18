package p2p

// Peer struct is a subset of the Node struct
// This way we need to only store some of Node
// information i.e. Node ID and Address in a
// Node's peer list and no need to store the Node.cache
// or Node.Listener
type Peer struct {
	ID      string // Unique ID (could be generated from public key or random)
	Address string // e.g., "127.0.0.1:4000"
}
