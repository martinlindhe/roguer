package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/martinlindhe/roguer"
)

var (
	upgrader = websocket.Upgrader{}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 1 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 2 * time.Second
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

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		t, msg, err := ws.ReadMessage()
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
				res := messageResponse{Type: "error", Message: fmt.Sprintf("%v", err)}
				b, _ = json.Marshal(res)
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
				res := messageResponse{Type: "error", Message: "invalid token"}
				b, _ = json.Marshal(res)
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
			log.Errorf("unknown command %s", parts[0])

			res := messageResponse{Type: "error", Message: fmt.Sprintf("unknown command %s", parts[0])}
			b, _ = json.Marshal(res)
		}

		ws.WriteMessage(t, b)
	}
}

func writer(ws *websocket.Conn) {

	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case <-pingTicker.C:
			log.Println("pingTicker")
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	go writer(ws)
	reader(ws)
}
