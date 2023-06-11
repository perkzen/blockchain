package network

import (
	"fmt"
)

func addNode(port uint16) {
	node := fmt.Sprintf("localhost:%d", port)
	knownNodes = append(knownNodes, node)

	for _, node := range knownNodes {
		if node == fmt.Sprintf("localhost:%d", port) {
			continue
		}

		//	res, err := http.Post(fmt.Sprintf("http://%s/nodes", node), "application/json", nil)
	}
}
