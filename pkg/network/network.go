package network

import (
	"blockchain/pkg/blockchain"
	"blockchain/pkg/wallet"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache = make(map[string]*blockchain.Blockchain)
var walletCache = make(map[string]*wallet.Wallet)

const (
	BLOCKCHAIN = "blockchain"
	WALLET     = "wallet"
)

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
		walletCache[WALLET] = minersWallet
		chain = blockchain.InitBlockchain(minersWallet.BlockchainAddress(), s.Port())
		cache[BLOCKCHAIN] = chain
	}

	return chain
}

func (s *Server) GetWallet() *wallet.Wallet {
	w, ok := walletCache[WALLET]
	if !ok {
		w = wallet.NewWallet()
		walletCache[WALLET] = w
	}
	return w
}

func (s *Server) transactionHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Println("ERROR: Invalid Method")
	}
	w.Header().Add("Content-Type", "application/json")

	type ReqBody struct {
		Recipient string  `json:"recipient"`
		Amount    float64 `json:"amount"`
	}

	var body ReqBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chain := s.GetBlockchain()
	wallet := s.GetWallet()
	tx := blockchain.NewTransaction(wallet.BlockchainAddress(), body.Recipient, 0.1, chain)
	chain.AddTransaction(wallet.BlockchainAddress(), body.Recipient, 0.1, tx.GenerateSignature(wallet.PrivateKey()), wallet.PublicKey())

	fmt.Println("⛏️  Mining block...")
	chain.MineBlock()
	fmt.Println("️✅  Block successfully mined")

	_, err = io.WriteString(w, "Transaction received and will be added to the blockchain")
	if err != nil {
		log.Fatal("ERROR: Failed to send JSON")
	}

}

func (s *Server) walletHandler() {

}

func (s *Server) chainHandler(w http.ResponseWriter, req *http.Request) {
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
	http.HandleFunc("/", s.chainHandler)
	http.HandleFunc("/transaction", s.transactionHandler)
	fmt.Printf("Blockchain node is listening on: http://localhost:%d\n", s.Port())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(s.Port())), nil))
}
