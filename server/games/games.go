package games

import "github.com/gorilla/websocket"

// Game is an interface for games
type Game interface {
	Execute(*websocket.Conn)
}

var queuing map[string]*websocket.Conn
