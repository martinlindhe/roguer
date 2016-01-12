package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gobuild/log"
	"github.com/gorilla/websocket"
	"github.com/martinlindhe/roguer"
)

// client represents a single user
type client struct {
	// socket is the web socket for this client.
	socket *websocket.Conn

	// send is a channel on which messages are sent.
	send chan []byte
}

func (c *client) read() {
	for {
		if t, msg, err := c.socket.ReadMessage(); err == nil {
			//c.room.forward <- msg

			b := []byte{}

			parts := strings.SplitN(string(msg), " ", 2)

			switch parts[0] {
			case "new_player":
				pos, token := island.NewPlayer(parts[1], c.socket)

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

			c.socket.WriteMessage(t, b)

		} else {
			break
		}
	}

	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
