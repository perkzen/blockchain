package network

import (
	"fmt"
	"golang.org/x/net/websocket"
	"strconv"
)

func (s *Server) addNode(addr string, ws *websocket.Conn) {
	s.nodes[addr] = ws
}

func (s *Server) removeNode(addr string) {
	delete(s.nodes, addr)
}

func (s *Server) connectToKnownNodes() {
	for addr := range s.nodes {
		connAddr := fmt.Sprintf("localhost:%d", s.Port())
		if addr == connAddr {
			continue
		}

		client := NewWebSocketClient("ws://" + addr + "/ws")
		s.nodes[addr] = client
		EmitEvent(client, CONNECT, connAddr)

	}
}

func (s *Server) getNodeAddresses() []string {
	addresses := []string{"localhost:" + strconv.Itoa(int(s.Port()))}
	for addr := range s.nodes {
		addresses = append(addresses, addr)
	}

	return addresses
}
