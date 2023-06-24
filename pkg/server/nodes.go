package server

import (
	"fmt"
	"strconv"
	"sync"
)

var mut sync.Mutex

func (s *Server) addNode(addr string) {
	mut.Lock()
	defer mut.Unlock()

	if addr == fmt.Sprintf("localhost:%d", s.Port()) {
		return
	}

	client := newWebSocketClient(addr)
	s.nodes[addr] = client
}

func (s *Server) removeNode(addr string) {
	mut.Lock()
	defer mut.Unlock()
	delete(s.nodes, addr)
}

func broadcastEvent[T any](s *Server, event Event, data T) {
	for _, ws := range s.nodes {
		emitEvent(ws, event, data)
	}
}

func (s *Server) connectToKnownNodes() {
	for addr := range s.nodes {
		connAddr := fmt.Sprintf("localhost:%d", s.Port())
		client := newWebSocketClient(addr)
		mut.Lock()
		s.nodes[addr] = client
		mut.Unlock()
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
