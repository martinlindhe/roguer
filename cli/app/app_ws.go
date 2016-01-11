package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/martinlindhe/roguer"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsResponse struct {
	Type string
}

type playerSpawnResponse struct {
	moveResponse
	Token string
}

type moveResponse struct {
	wsResponse
	X           float64
	Y           float64
	LocalSpawns []rogue.LocalSpawns
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
			pos, token := island.NewPlayer(parts[1])

			var res playerSpawnResponse
			res.Type = "xy"
			res.X = pos.X
			res.Y = pos.Y
			res.Token = token
			res.LocalSpawns = island.DescribeLocalArea(pos)

			b, _ = json.Marshal(res)
			log.Printf("new player %s spawned at %s", parts[1], pos)

			// XXX broadcast a "new player" event to all

		case "continue":
			pos, token, err := island.ContinuePlayer(parts[1])
			if err != nil {
				errMsg := fmt.Sprintf("%v", err)
				b = []byte(`{"Type": "error", "Message": "` + errMsg + `"}`)
				break
			}

			var res playerSpawnResponse
			res.Type = "xy"
			res.X = pos.X
			res.Y = pos.Y
			res.Token = token
			res.LocalSpawns = island.DescribeLocalArea(*pos)

			b, _ = json.Marshal(res)
			log.Printf("continuing player %s spawned at %s", parts[1], pos)

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
				b = []byte(`{"Type": "error", "Message": "invalid token"}`)
				break
			}

			oldPos := player.Spawn.Position
			player.Spawn.Position.X = float64(x)
			player.Spawn.Position.Y = float64(y)

			log.Printf("Player %s moved from %s to %s", player.Name, oldPos, player.Spawn.Position)

			var res moveResponse
			res.Type = "move_res"
			res.X = player.Spawn.Position.X
			res.Y = player.Spawn.Position.X
			res.LocalSpawns = island.DescribeLocalArea(player.Spawn.Position)
			b, _ = json.Marshal(res)

		default:
			errMsg := fmt.Sprintf("unknown command %s", parts[0])
			b = []byte(`{"Type": "error", "Message": "` + errMsg + `"}`)
			log.Errorf("unknown command %s", parts[0])
		}

		conn.WriteMessage(t, b)
	}
}
