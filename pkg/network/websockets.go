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
	CONNECT     Event = "CONNECT"
	DISCONNECT  Event = "DISCONNECT"
)

type Message[T any] struct {
	Data  T     `json:"data"`
	Event Event `json:"event"`
}

func ReadLoop(ws *websocket.Conn, server *Server) {
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
		case CONNECT:
			fmt.Println("New node")
			server.knownNodes = append(server.knownNodes, msg.Data)
		case DISCONNECT:
			fmt.Println("Node disconnected")
			for i, node := range server.knownNodes {
				if node == msg.Data {
					server.knownNodes = append(server.knownNodes[:i], server.knownNodes[i+1:]...)
					break
				}
			}
		default:
			fmt.Println("Unknown event")
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

func NewWebSocketClient(url string) *websocket.Conn {
	conn, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		fmt.Println("ERROR: Failed to connect to websocket")
		return nil
	}
	return conn
}