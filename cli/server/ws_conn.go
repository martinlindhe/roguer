package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/martinlindhe/roguer"
)

var (
	upgrader          = websocket.Upgrader{}
	messageBufferSize = 8192
)

type messageResponse struct {
	Type    string
	Message string
}

type moveResponse struct {
	Type        string
	X           float64
	Y           float64
	LocalSpawns []rogue.LocalSpawns
}

type playerSpawnResponse struct {
	moveResponse
	Token string
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
	}

	go client.write()
	client.read()
}
