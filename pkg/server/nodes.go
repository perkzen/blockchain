package server

import (
	"fmt"
	"strconv"
)

func (s *Server) addNode(addr string) {
	client := newWebSocketClient(addr)
	s.nodes[addr] = client
}

func (s *Server) removeNode(addr string) {
	delete(s.nodes, addr)
}

func broadcastEvent[T any](s *Server, event Event, data T) {
	for _, ws := range s.nodes {
		emitEvent(ws, NEW_NODE, data)
	}
}

func (s *Server) connectToKnownNodes() {
	for addr := range s.nodes {
		connAddr := fmt.Sprintf("localhost:%d", s.Port())
		client := newWebSocketClient(addr)
		s.nodes[addr] = client
		emitEvent(client, CONNECT, connAddr)
	}
}

func (s *Server) getNodeAddresses() []string {
	addresses := []string{"localhost:" + strconv.Itoa(int(s.Port()))}
	for addr := range s.nodes {
		addresses = append(addresses, addr)
	}

	return addresses
}
