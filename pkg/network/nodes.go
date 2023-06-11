package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func (s *Server) AddNodeIfNotKnown(node string) {
	if node == fmt.Sprintf("localhost:%d", s.Port()) {
		return
	}

	var found bool
	for _, knownNode := range knownNodes {
		if knownNode == node {
			found = true
			break
		}
	}
	if !found {
		knownNodes = append(knownNodes, node)
	}
}

func (s *Server) AddNode(node string) {
	s.AddNodeIfNotKnown(node)

	for _, knownNode := range knownNodes {
		if knownNode == node {
			continue
		}

		body := map[string]string{"node": node}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			log.Fatal(err)
		}

		res, err := http.Post(fmt.Sprintf("http://%s/nodes", knownNode), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			log.Fatal("ERROR: Failed to add knownNode")
		}

		respBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(respBytes))
	}
}

func getNodes() {
	for _, node := range knownNodes {
		res, err := http.Get(fmt.Sprintf("http://%s/nodes", node))
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			log.Fatal("ERROR: Failed to get nodes")
		}

		var nodes []string
		err = json.NewDecoder(res.Body).Decode(&nodes)
		if err != nil {
			log.Fatal(err)
		}

		for _, node := range nodes {
			// check if node is already in knownNodes
			var found bool
			for _, knownNode := range knownNodes {
				if knownNode == node {
					found = true
					break
				}
			}
			if !found {
				knownNodes = append(knownNodes, node)
			}
		}
	}
}

func (s *Server) SearchNodes() {
	ticker := time.NewTicker(NODES_SYNC_TIMEOUT)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("👷‍ Searching for new nodes...")
				getNodes()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
