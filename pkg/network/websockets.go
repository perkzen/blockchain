package network

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
)

type Event string

const (
	BLOCK_MINED Event = "BLOCK_MINED"
	NEW_NODE    Event = "NEW_NODE"
)

type Message[T any] struct {
	Data  T     `json:"data"`
	Event Event `json:"event"`
}

func ReadLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("ERROR: Failed to read from websocket")
			continue
		}

		var msg Message[string]
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			fmt.Println("ERROR: Failed to unmarshal JSON")
			continue
		}

		switch msg.Event {
		case BLOCK_MINED:
			fmt.Println("Block mined")
		case NEW_NODE:
			fmt.Println("New node")
		}

		fmt.Println("Received message:", msg)
		//_, err = ws.Write([]byte("Received message: " + msg))
		//if err != nil {
		//	fmt.Println("ERROR: Failed to write to websocket")
		//	continue
		//}
	}
}

func EmitEvent[T any](ws *websocket.Conn, event Event, data T) {
	msg := Message[T]{
		Data:  data,
		Event: event,
	}
	err := websocket.JSON.Send(ws, msg)
	if err != nil {
		fmt.Println("ERROR: Failed to send JSON")
	}
}

func NewWebSocketClient() *websocket.Conn {
	conn, err := websocket.Dial("ws://localhost:3001/ws", "", "http://localhost")
	if err != nil {
		fmt.Println("ERROR: Failed to connect to websocket")
		return nil
	}
	return conn
}
