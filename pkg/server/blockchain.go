package server

import (
	"blockchain/pkg/blockchain"
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

	serverAddr := fmt.Sprintf("localhost:%d", s.Port())

	if ok {
		chains[serverAddr] = chain.Length()
	}

	// searching nearby nodes for the longest chain
	for node := range s.nodes {
		// skip self
		if node == serverAddr {
			continue
		}

		res, err := http.Get(fmt.Sprintf("http://%s/", node))
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode == http.StatusOK {

			var c *blockchain.Blockchain
			err := json.NewDecoder(res.Body).Decode(&c)
			if err != nil {
				log.Fatal(err)
			}
			chains[node] = c.Length()
			fmt.Printf("📦🔗📦 %s: %d\n", node, c.Length())
		}
	}

	// get the longest chain
	var longestChainLength int
	var longestChainAddr string
	for node, length := range chains {
		if length > longestChainLength {
			longestChainLength = length
			longestChainAddr = node
		}
	}
	fmt.Printf("📦🔗📦 Longest chain: %s: %d\n", longestChainAddr, longestChainLength)

	// if the longest chain is the current node, return the chain
	if longestChainAddr == serverAddr {
		return chain
	}

	// if no other nodes are running, create a new blockchain
	if longestChainLength == 0 {
		fmt.Printf("📦🔗📦 No other nodes running, creating new blockchain\n")
		longestBlockchain = s.createBlockchain()
		return longestBlockchain
	}

	// get chain from longest chain node
	res, err := http.Get(fmt.Sprintf("http://%s/", longestChainAddr))
	if err != nil {
		log.Fatal(err)
	}
	var c *blockchain.Blockchain
	err = json.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}

	longestBlockchain = c
	return longestBlockchain
}

func (s *Server) getBlockchain() *blockchain.Blockchain {
	chain, ok := s.cache[BLOCKCHAIN].(*blockchain.Blockchain)
	if !ok {
		chain = s.getLongestChain()
		s.cache[BLOCKCHAIN] = chain
	}
	return chain
}

func (s *Server) createBlockchain() *blockchain.Blockchain {
	w := s.GetWallet()
	chain := blockchain.InitBlockchain(w.BlockchainAddress(), s.Port())
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

	chain := s.getBlockchain()
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
