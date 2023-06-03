package network

import (
	"blockchain/pkg/blockchain"
	"blockchain/pkg/wallet"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache = make(map[string]*blockchain.Blockchain)

const BLOCKCHAIN = "blockchain"

type Server struct {
	port uint16
}

func NewBlockchainServer(port uint16) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) GetBlockchain() *blockchain.Blockchain {
	chain, ok := cache[BLOCKCHAIN]
	if !ok {
		minersWallet := wallet.NewWallet()
		chain = blockchain.InitBlockchain(minersWallet.BlockchainAddress(), s.Port())
		cache[BLOCKCHAIN] = chain
	}

	return chain
}

func (s *Server) GetChain(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Println("ERROR: Invalid Method")
	}
	w.Header().Add("Content-Type", "application/json")

	chain := s.GetBlockchain()
	m, _ := chain.MarshalJSON()
	_, err := io.WriteString(w, string(m[:]))
	if err != nil {
		log.Fatal("ERROR: Failed to send JSON")
	}

}

func (s *Server) Port() uint16 {
	return s.port
}

func (s *Server) Run() {
	http.HandleFunc("/", s.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(s.Port())), nil))
}
