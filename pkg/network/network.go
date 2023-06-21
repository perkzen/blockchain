package network

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	port       uint16
	cache      map[string]interface{}
	knownNodes []string
}

func NewBlockchainServer(port uint16) *Server {
	return &Server{
		port:       port,
		cache:      make(map[string]interface{}),
		knownNodes: []string{"localhost:3001"},
	}
}

func (s *Server) Port() uint16 {
	return s.port
}

func (s *Server) Run() {
	s.initHandlers()
	fmt.Println("Listening on port", s.Port())
	//http.HandleFunc("/wallet", s.walletHandler)
	//s.AddNode("localhost:" + strconv.Itoa(int(s.Port())))
	//s.SearchNodes()
	//s.SyncChains()
	//s.MineBlock()
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(s.Port())), nil))
}
