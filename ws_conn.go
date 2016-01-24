package rogue

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

var (
	upgrader          = websocket.Upgrader{}
	messageBufferSize = 8192
)

type messageResponse struct {
	Type    string
	Message string
	Time    int64 // ticks since server reset
}

type moveResponse struct {
	Type        string
	X           float64
	Y           float64
	LocalSpawns []LocalSpawn
}

type tickMessage struct {
	Type          string
	FormattedTime string // XXX: remove this when js have int-to-time code
	Time          int64
}

type playerSpawnResponse struct {
	moveResponse
	Token string
}

func serveWebsocket(g *Game, w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}

	client := &client{
		socket: socket,
		game:   g,
		send:   make(chan []byte, messageBufferSize),
	}

	go client.write()
	client.read()
}
