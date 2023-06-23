package server

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
)

type Event string

const (
	BLOCK_MINED Event = "BLOCK_MINED"
	CONNECT     Event = "CONNECT"
	DISCONNECT  Event = "DISCONNECT"
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
			continue
		}

		var msg Message[any]
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			fmt.Println("ERROR: Failed to unmarshal JSON")
			continue
		}

		switch msg.Event {
		case BLOCK_MINED:
			fmt.Println("Block mined")
		case CONNECT:
			fmt.Println("New node")
			address := msg.Data.(string)
			s.addNode(address, ws)
		case DISCONNECT:
			fmt.Println("Node disconnected")
			address := msg.Data.(string)
			s.removeNode(address)
		default:
			fmt.Println("Unknown event")
		}

		fmt.Println("Received message:", msg)
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

func newWebSocketClient(url string) *websocket.Conn {
	conn, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		fmt.Println("ERROR: Failed to connect to websocket")
		return nil
	}
	return conn
}
