package network

import (
	"blockchain/pkg/blockchain"
	"blockchain/pkg/wallet"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (s *Server) getLongestChain() *blockchain.Blockchain {
	var longestBlockchain *blockchain.Blockchain

	chains := make(map[string]int)
	chain, ok := chainCache[BLOCKCHAIN]

	if ok {
		chains[fmt.Sprintf("localhost:%d", s.Port())] = chain.Length()
	}

	for _, node := range knownNodes {

		// skip self
		if node == fmt.Sprintf("localhost:%d", s.Port()) {
			continue
		}

		res, err := http.Get(fmt.Sprintf("http://%s/", node))
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode == http.StatusOK {

			var c blockchain.Blockchain
			err := json.NewDecoder(res.Body).Decode(&c)
			if err != nil {
				log.Fatal(err)
			}
			chains[node] = c.Length()
		}
	}

	// get the longest chain
	var longestChainLength int
	var longestChainNode string
	for node, length := range chains {
		if length > longestChainLength {
			longestChainLength = length
			longestChainNode = node
		}
	}

	// if no other nodes are running, create a new blockchain
	if longestChainNode == "" {
		minersWallet := wallet.NewWallet()
		walletCache[WALLET] = minersWallet
		longestBlockchain = blockchain.InitBlockchain(minersWallet.BlockchainAddress(), s.Port())
		chainCache[BLOCKCHAIN] = longestBlockchain
		return longestBlockchain
	}

	// get chain from longest chain node
	res, err := http.Get(fmt.Sprintf("http://%s/", longestChainNode))
	if err != nil {
		log.Fatal(err)
	}
	var c *blockchain.Blockchain
	err = json.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
	chainCache[BLOCKCHAIN] = c
	longestBlockchain = c
	return longestBlockchain
}

func (s *Server) GetBlockchain() *blockchain.Blockchain {
	chain, ok := chainCache[BLOCKCHAIN]
	if !ok {
		chain = s.getLongestChain()
	}
	return chain
}

func (s *Server) SyncChains() {
	ticker := time.NewTicker(CHAIN_SYNC_TIMEOUT)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("📦 🔗 📦 Syncing chains...")
				s.getLongestChain()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Server) MineBlock() {
	ticker := time.NewTicker(MINING_TIMEOUT)
	quit := make(chan struct{})

	chain := s.GetBlockchain()
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("⛏️  Mining block...")
				chain.MineBlock()
				fmt.Println("️✅  Block successfully mined")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
