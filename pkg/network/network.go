package network

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	port       uint16
	cache      map[string]interface{}
	knownNodes []string
	conns      map[*websocket.Conn]bool
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

	if s.port != 3001 {
		client := NewWebSocketClient()
		EmitEvent(client, NEW_NODE, fmt.Sprintf("localhost:%d", s.Port()))
	}

	//s.AddNode("localhost:" + strconv.Itoa(int(s.Port())))
	//s.SearchNodes()
	//s.SyncChains()
	//s.MineBlock()
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(s.Port())), nil))
}
