package network

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (s *Server) transactionHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		fmt.Println("ERROR: Invalid Method")
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")

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
	chain.AddTransaction(wallet.BlockchainAddress(), body.Recipient, 0.1, wallet.PrivateKey(), wallet.PublicKey())

	_, err = io.WriteString(w, "Transaction received and will be added to the blockchain")
	if err != nil {
		log.Fatal("ERROR: Failed to send JSON")
	}

}

func (s *Server) chainHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Println("ERROR: Invalid Method")
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	chain := s.GetBlockchain()
	m, _ := chain.MarshalJSON()
	_, err := io.WriteString(w, string(m[:]))
	if err != nil {
		log.Fatal("ERROR: Failed to send JSON")
	}
}

func (s *Server) nodeHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodGet {

		nodes, _ := json.Marshal(s.knownNodes)

		_, err := io.WriteString(w, string(nodes))
		if err != nil {
			log.Fatal("ERROR: Failed to send JSON")
		}
	}

	if req.Method == http.MethodPost {

		type ReqBody struct {
			Node string `json:"node"`
		}

		var body ReqBody
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.AddNodeIfNotKnown(body.Node)

		_, err = io.WriteString(w, "Node added successfully")
		if err != nil {
			log.Fatal("ERROR: Failed to send JSON")
		}
	}
}

func (s *Server) walletHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if req.Method != http.MethodGet {
		fmt.Println("ERROR: Invalid Method")
	}

	wallet := s.GetWallet()
	chain := s.GetBlockchain()
	balance := chain.UTXO.GetBalance(wallet.BlockchainAddress())

	type RespBody struct {
		Address string  `json:"address"`
		Balance float32 `json:"balance"`
	}

	respBody := &RespBody{
		Address: wallet.BlockchainAddress(),
		Balance: balance,
	}

	m, _ := json.Marshal(respBody)
	_, err := io.WriteString(w, string(m[:]))
	if err != nil {
		log.Fatal("ERROR: Failed to send JSON")
	}
}

func (s *Server) initHandlers() {
	http.HandleFunc("/", s.chainHandler)
	http.HandleFunc("/transaction", s.transactionHandler)
	http.HandleFunc("/nodes", s.nodeHandler)
	http.HandleFunc("/wallet", s.walletHandler)
}
