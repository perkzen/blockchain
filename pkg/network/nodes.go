package network

import "fmt"

func (s *Server) addNode() {

}

func (s *Server) removeNode() {

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
	var addresses []string
	for addr := range s.nodes {
		addresses = append(addresses, addr)
	}

	return addresses
}
