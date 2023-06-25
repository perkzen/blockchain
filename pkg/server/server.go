package server

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
	port  uint16
	cache map[string]interface{}
	nodes map[string]*websocket.Conn
}

func NewBlockchainServer(port uint16) *Server {

	knownNodes := []string{"localhost:3000"}

	s := &Server{
		port:  port,
		cache: make(map[string]interface{}),
		nodes: make(map[string]*websocket.Conn),
	}

	for _, node := range knownNodes {
		if node == fmt.Sprintf("localhost:%d", s.Port()) {
			continue
		}
		s.nodes[node] = newWebSocketClient(node)
	}

	return s

}

func (s *Server) Port() uint16 {
	return s.port
}

func (s *Server) gracefulShutdown() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		fmt.Printf("\nShutting down server on: http//:localhost:%d\n", s.Port())
		for _, conn := range s.nodes {
			emitEvent(conn, DISCONNECT, fmt.Sprintf("localhost:%d", s.Port()))
		}
		os.Exit(0)
	}()
}

func (s *Server) Run() {
	s.gracefulShutdown()
	s.initHandlers()
	s.connectToKnownNodes()
	s.getBlockchain()
	s.startMining()
	s.SyncChains()
	fmt.Printf("Listening on: http//:localhost:%d\n", s.Port())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(s.Port())), nil))

}
