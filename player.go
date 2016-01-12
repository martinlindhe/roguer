package rogue

import "github.com/gorilla/websocket"

// Player ...
type Player struct {
	Name   string
	Token  string
	Spawn  *Obj // points to a npc in game
	Socket *websocket.Conn
}
