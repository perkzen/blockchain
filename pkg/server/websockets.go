package server

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
)

type Event string

const (
	BLOCK_MINED Event = "BLOCK_MINED"
	CONNECT     Event = "CONNECT"
	DISCONNECT  Event = "DISCONNECT"
	NEW_NODE    Event = "NEW_NODE"
	REMOVE_NODE Event = "REMOVE_NODE"
)

type Message[T any] struct {
	Data  T     `json:"data"`
	Event Event `json:"event"`
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

		var msg Message[any]
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			fmt.Println("ERROR: Failed to unmarshal JSON", err)
			continue
		}

		switch msg.Event {
		case BLOCK_MINED:
			fmt.Println("Block mined")
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

func emitEvent[T any](ws *websocket.Conn, event Event, data T) {
	msg := Message[T]{
		Data:  data,
		Event: event,
	}
	err := websocket.JSON.Send(ws, msg)
	if err != nil {
		fmt.Println("ERROR: Failed to send JSON")
	}
}

func newWebSocketClient(addr string) *websocket.Conn {
	conn, err := websocket.Dial("ws://"+addr+"/ws", "", "http://localhost")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
