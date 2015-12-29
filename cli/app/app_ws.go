package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/qiniu/log"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsResponse struct {
	Type string
}

type newPlayerResponse struct {
	wsResponse
	X     float64
	Y     float64
	Token string
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("ws handler")

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		b := []byte{}

		parts := strings.SplitN(string(msg), " ", 2)

		switch parts[0] {
		case "new_player":
			// XXX create new player etc
			pos := island.RandomPointAboveWater()

			var res newPlayerResponse
			res.Type = "xy"
			res.X = pos.X
			res.Y = pos.Y
			res.Token = newJwt()

			b, _ = json.Marshal(res)
			log.Printf("new player %s spawned at %s", parts[1], pos)

			// XXX broadcast a "new player" event to all

		case "ping":
			b = []byte(`{"Type": "pong"}`)

		default:
			b = []byte(fmt.Sprintf("unknown command %s", parts[0]))
			log.Errorf("unknown command %s", parts[0])
		}

		conn.WriteMessage(t, b)
	}
}
