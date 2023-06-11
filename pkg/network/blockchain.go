package network

import (
	"blockchain/pkg/blockchain"
	"blockchain/pkg/wallet"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) GetBlockchain() *blockchain.Blockchain {
	chain, ok := chainCache[BLOCKCHAIN]
	if !ok {
		chains := make(map[string]int)
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
			chain = blockchain.InitBlockchain(minersWallet.BlockchainAddress(), s.Port())
			chainCache[BLOCKCHAIN] = chain
			return chain
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
		chain = c
	}
	return chain
}

func (s *Server) SyncChains() {

}
