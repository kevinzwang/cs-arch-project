package games

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Game is an interface for games
type Game interface {
	Execute(*websocket.Conn)
}

var queuing map[string]*websocket.Conn

func init() {
	queuing = make(map[string]*websocket.Conn)
}

func chanReadJSON(conn *websocket.Conn, resp chan<- map[string]interface{}) {
	var foo map[string]interface{}
	err := conn.ReadJSON(&foo)
	if err != nil {
		fmt.Printf("error reading JSON: %v\n", err)
		resp <- nil
	}
	resp <- foo
}
