package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/martinlindhe/rogue"
	"github.com/martinlindhe/rogue/views"
	"github.com/plimble/ace"
)

var island *rogue.Island
var islandMap string

func precalcTilemap() string {
	var tileMap rogue.PhaserTileMap
	tileMap.Version = 1
	tileMap.Width = island.Width
	tileMap.Height = island.Height
	tileMap.TileWidth = 32
	tileMap.TileHeight = 32
	tileMap.Orientation = "orthogonal"

	var layer rogue.PhaserTileLayer
	layer.Data = island.HeightsAsFlatTilemap()
	layer.Width = island.Width
	layer.Height = island.Height
	layer.Visible = true
	layer.Opacity = 1
	layer.Type = "tilelayer"
	layer.Name = "layer1"
	tileMap.Layers = append(tileMap.Layers, layer)

	var tileset rogue.PhaserTileSet
	tileset.FirstGid = 0
	// NOTE: need to specify a tile in phaser later, .Name and .Image must be the same value (phaser 2.4.4, dec 2015)
	tileset.Name = "island_tiles"
	tileset.Image = "island_tiles"
	tileset.ImageHeight = 256
	tileset.ImageWidth = 256
	tileset.TileWidth = 32
	tileset.TileHeight = 32
	tileMap.TileSets = append(tileMap.TileSets, tileset)

	b, _ := json.Marshal(tileMap)

	return string(b)
}

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("rogue started")

	r := getRouter()
	island = rogue.NewIsland()

	islandMap = precalcTilemap()

	r.GET("/", func(c *ace.C) {
		c.String(200, views.Index())
	})

	// listen and serve on 0.0.0.0:3322
	appPort := 3322
	listenAt := fmt.Sprintf(":%d", appPort)

	log.Infof("Starting http server on %s", listenAt)

	r.Run(listenAt)
}

func getRouter() *ace.Ace {

	// ace with Logger, Recovery
	r := ace.Default()

	//	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.GET("/island/full", getFullIslandController)

	r.GET("/ws", func(c *ace.C) {
		wsHandler(c.Writer, c.Request)
	})

	r.Static("/js", "./public/js")
	r.Static("/css", "./public/css")
	r.Static("/fonts", "./public/fonts")
	r.Static("/img", "./public/img")
	r.Static("/audio", "./public/audio")
	//r.LoadHTMLFiles("./public/index.html")
	return r
}

func newJwt() string {

	token := jwt.New(jwt.SigningMethodHS256)

	signingKey := []byte("top secret")

	//token.Claims["foo"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

// returns a map in Tiled json format, recognized by phaser.io
func getFullIslandController(c *ace.C) {
	// NOTE: this is useful in early stage for world debugging.
	// later on, the game would only expose a small area around the player

	c.String(http.StatusOK, islandMap)
}
