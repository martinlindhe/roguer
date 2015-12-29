package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/martinlindhe/rogue"
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
	Name  string
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

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
			// create new player
			pos := island.RandomPointAboveWater()
			token := newJwt()

			spawn := new(rogue.Obj)
			spawn.Name = parts[1]
			spawn.Position = pos
			island.Spawns = append(island.Spawns, spawn)

			var player rogue.Player
			player.Name = parts[1]
			player.Token = token
			player.Spawn = spawn
			island.Players = append(island.Players, player)

			var res newPlayerResponse
			res.Type = "xy"
			res.X = pos.X
			res.Y = pos.Y
			res.Name = parts[1]
			res.Token = token

			b, _ = json.Marshal(res)
			log.Printf("new player %s spawned at %s", parts[1], pos)

			// XXX broadcast a "new player" event to all

		case "move":
			subcommand := strings.SplitN(parts[1], " ", 3)
			x, _ := strconv.Atoi(subcommand[0])
			y, _ := strconv.Atoi(subcommand[1])
			token := subcommand[2]

			// find user by token
			var player *rogue.Player
			for _, u := range island.Players {
				if u.Token == token {
					player = &u
				}
			}
			if player == nil {
				log.Errorf("Invalid token recieved: %s", token)
				b = []byte(`{"Type": "error, invalid token"}`)
			} else {
				oldPos := player.Spawn.Position
				player.Spawn.Position.X = float64(x)
				player.Spawn.Position.Y = float64(y)

				log.Printf("Player %s moved from %s to %s", player.Name, oldPos, player.Spawn.Position)
				b = []byte(`{"Type": "ok"}`)
			}

		default:
			b = []byte(fmt.Sprintf("unknown command %s", parts[0]))
			log.Errorf("unknown command %s", parts[0])
		}

		conn.WriteMessage(t, b)
	}
}
