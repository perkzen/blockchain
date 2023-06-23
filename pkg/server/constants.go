package server

import "time"

const (
	MINING_TIMEOUT     = 2 * time.Second
	CHAIN_SYNC_TIMEOUT = 10 * time.Second
	NODES_SYNC_TIMEOUT = 5 * time.Second
)

// cache keys
const (
	BLOCKCHAIN = "blockchain"
	WALLET     = "wallet"
)
