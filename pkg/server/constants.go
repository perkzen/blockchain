package server

import "time"

const (
	MINING_TIMEOUT     = 3 * time.Second
	CHAIN_SYNC_TIMEOUT = 20 * time.Second
)

// cache keys
const (
	BLOCKCHAIN = "blockchain"
	WALLET     = "wallet"
)

// events
const (
	NEW_BLOCK   Event = "NEW_BLOCK"
	CONNECT     Event = "CONNECT"
	DISCONNECT  Event = "DISCONNECT"
	NEW_NODE    Event = "NEW_NODE"
	REMOVE_NODE Event = "REMOVE_NODE"
)
