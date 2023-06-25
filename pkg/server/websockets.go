package server

import (
	"blockchain/pkg/blockchain"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
)

type Event string

type Message struct {
	Data  interface{} `json:"data"`
	Event Event       `json:"event"`
}

func (m *Message) MarshallJSON() ([]byte, error) {
	return json.Marshal(m)
}

func readLoop(ws *websocket.Conn, s *Server) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		var msg Message
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			fmt.Println("ERROR: Failed to unmarshal JSON", err)
			continue
		}

		switch msg.Event {
		case NEW_BLOCK:
			fmt.Println("Block mined")
			var block blockchain.Block

			jsonString, _ := json.Marshal(msg.Data)

			err := json.Unmarshal(jsonString, &block)
			if err != nil {
				fmt.Println(err)
			}

			pow := blockchain.NewProofOfWork(&block)
			valid := pow.IsValid(block.Nonce)
			if valid {
				chain := s.getBlockchain()
				chain.Blocks = append(chain.Blocks, &block)
				fmt.Println("✅ New block added")
			}

		case CONNECT:
			fmt.Println("New connection")
			address := msg.Data.(string)
			s.addNode(address)
			broadcastEvent(s, NEW_NODE, address)
		case DISCONNECT:
			fmt.Println("Node disconnected")
			address := msg.Data.(string)
			s.removeNode(address)
			broadcastEvent(s, REMOVE_NODE, address)
		case NEW_NODE:
			fmt.Println("New node")
			address := msg.Data.(string)
			s.addNode(address)
		case REMOVE_NODE:
			fmt.Println("Node disconnected")
			address := msg.Data.(string)
			s.removeNode(address)
		default:
			fmt.Println("Unknown event")
		}
	}
}

func emitEvent[T interface{}](ws *websocket.Conn, event Event, data T) {
	msg := Message{
		Data:  data,
		Event: event,
	}

	m, _ := msg.MarshallJSON()

	_, err := ws.Write(m)
	if err != nil {
		fmt.Println(err)
	}
}

func broadcastEvent(s *Server, event Event, data interface{}) {
	for _, ws := range s.nodes {
		emitEvent(ws, event, data)
	}
}

func newWebSocketClient(addr string) *websocket.Conn {
	conn, err := websocket.Dial("ws://"+addr+"/ws", "", "http://localhost")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
