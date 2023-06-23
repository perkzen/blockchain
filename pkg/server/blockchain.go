package server

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
	chain, ok := s.cache[BLOCKCHAIN].(*blockchain.Blockchain)

	if ok {
		chains[fmt.Sprintf("localhost:%d", s.Port())] = chain.Length()
	}

	for node := range s.nodes {
		// skip self
		if node == fmt.Sprintf("localhost:%d", s.Port()) {
			continue
		}

		res, err := http.Get(fmt.Sprintf("http://%s/", node))
		if err != nil {
			log.Fatal(err)
		}
		if res.StatusCode == http.StatusOK {

			var c *blockchain.Blockchain
			err := json.NewDecoder(res.Body).Decode(&c)
			if err != nil {
				log.Fatal(err)
			}
			chains[node] = c.Length()
			fmt.Printf("🔗 %s: %d\n", node, c.Length())
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
	fmt.Printf("🔗 Longest chain: %s: %d\n", longestChainNode, longestChainLength)

	if longestChainNode == fmt.Sprintf("localhost:%d", s.Port()) {
		return chain
	}

	// if no other nodes are running, create a new blockchain
	if longestChainLength == 0 {
		fmt.Printf("🔗 No other nodes running, creating new blockchain\n")
		minersWallet := wallet.NewWallet()
		s.cache[WALLET] = minersWallet
		longestBlockchain = blockchain.InitBlockchain(minersWallet.BlockchainAddress(), s.Port())
		s.cache[BLOCKCHAIN] = longestBlockchain
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
	s.cache[BLOCKCHAIN] = c
	longestBlockchain = c
	return longestBlockchain
}

func (s *Server) GetBlockchain() *blockchain.Blockchain {
	chain, ok := s.cache[BLOCKCHAIN].(*blockchain.Blockchain)
	if !ok {
		chain = s.getLongestChain()
		s.cache[BLOCKCHAIN] = chain
	}
	return chain
}

func (s *Server) ValidateBlock() {
	// new get block
	// validate it with proof of work
	// if 51% thinks that its valid, add it to the chain

}

func (s *Server) BroadcastNewBlock() {
	// broadcast new block to all nodes
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
