package network

import (
	"blockchain/pkg/blockchain"
	"blockchain/pkg/wallet"
	"log"
	"net/http"
	"strconv"
)

var chainCache = make(map[string]*blockchain.Blockchain)
var walletCache = make(map[string]*wallet.Wallet)

const (
	BLOCKCHAIN = "blockchain"
	WALLET     = "wallet"
)

var knownNodes = []string{"localhost:3001"}

type Server struct {
	port uint16
}

func NewBlockchainServer(port uint16) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Port() uint16 {
	return s.port
}

func (s *Server) GetWallet() *wallet.Wallet {
	w, ok := walletCache[WALLET]
	if !ok {
		w = wallet.NewWallet()
		walletCache[WALLET] = w
	}
	return w
}

func (s *Server) Run() {
	http.HandleFunc("/", s.chainHandler)
	http.HandleFunc("/transaction", s.transactionHandler)
	http.HandleFunc("/nodes", s.nodeHandler)
	http.HandleFunc("/wallet", s.walletHandler)
	s.AddNode("localhost:" + strconv.Itoa(int(s.Port())))
	s.SearchNodes()
	s.SyncChains()
	s.MineBlock()
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(s.Port())), nil))
}
