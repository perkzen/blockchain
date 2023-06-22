package network

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

type Server struct {
	port       uint16
	cache      map[string]interface{}
	knownNodes []string
	conns      []*websocket.Conn
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
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		fmt.Println("Shutting down...")
		for _, conn := range s.conns {
			EmitEvent(conn, DISCONNECT, fmt.Sprintf("localhost:%d", s.Port()))
		}
		os.Exit(0)
	}()

	s.initHandlers()
	fmt.Println("Listening on port", s.Port())

	for _, node := range s.knownNodes {
		if node == fmt.Sprintf("localhost:%d", s.Port()) {
			continue
		}
		client := NewWebSocketClient("ws://" + node + "/ws")
		s.conns = append(s.conns, client)
		fmt.Println(client.RemoteAddr())
		EmitEvent(client, CONNECT, fmt.Sprintf("localhost:%d", s.Port()))
	}

	//s.AddNode("localhost:" + strconv.Itoa(int(s.Port())))
	//s.SearchNodes()
	//s.SyncChains()
	//s.MineBlock()
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(s.Port())), nil))

}
