package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/martinlindhe/rogue"
	"github.com/martinlindhe/rogue/views"
	"github.com/plimble/ace"
)

var island *rogue.Island

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("rogue started")

	r := getRouter()
	island = rogue.NewIsland()

	r.GET("/", func(c *ace.C) {
		c.String(200, views.Index())
	})

	// XXX run http server in separate process? we also need a websock server

	// listen and serve on 0.0.0.0:8080
	appPort := 3322
	listenAt := fmt.Sprintf(":%d", appPort)

	log.Infof("Starting http server on %s", listenAt)

	r.Run(listenAt)
}

func getRouter() *ace.Ace {

	// ace with Logger, Recovery
	r := ace.Default()

	//	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/ping", pingController)

	r.GET("/island/full", getFullIslandController)

	r.Static("/js", "./public/js")
	r.Static("/css", "./public/css")
	r.Static("/fonts", "./public/fonts")
	r.Static("/img", "./public/img")
	//r.Static("/flags", "./public/flags")
	//r.LoadHTMLFiles("./public/index.html")
	return r
}

// curl -v "http://localhost:8080/ping"
func pingController(c *ace.C) {
	c.JSON(200, map[string]string{"pong": "now"})
}

type phaserTileMap struct {
	Version     int               `json:"version"`
	Width       int               `json:"width"`
	Height      int               `json:"height"`
	TileWidth   int               `json:"tilewidth"`
	TileHeight  int               `json:"tileheight"`
	Orientation string            `json:"orientation"`
	Layers      []phaserTileLayer `json:"layers"`
	TileSets    []phaserTileSet   `json:"tilesets"`
	// Properties ....  we skipped this
}

type phaserTileLayer struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Data    []int  `json:"data"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Opacity int    `json:"opacity"`
	Visible bool   `json:"visible"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
}

type phaserTileSet struct {
	FirstGid    int    `json:"firstgid"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	ImageHeight int    `json:"imageheight"`
	ImageWidth  int    `json:"imagewidth"`
	Margin      int    `json:"margin"`
	Spacing     int    `json:"spacing"`
	TileHeight  int    `json:"tileheight"`
	TileWidth   int    `json:"tilewidth"`

	// Properties ....  we skipped this
}

func getFullIslandController(c *ace.C) {
	// NOTE: this is useful in early stage for world debugging.
	// later on, the game would only expose a small area around the player

	// return a map in Phaser.Tilemap.TILED_JSON format

	// return width, height, heigthmap only
	var tileMap phaserTileMap
	tileMap.Version = 1
	tileMap.Width = island.Width
	tileMap.Height = island.Height
	tileMap.TileWidth = 16
	tileMap.TileHeight = 16
	tileMap.Orientation = "orthogonal"

	var layer phaserTileLayer
	layer.Data = island.HeightsAsFlatTilemap()
	layer.Width = 170
	layer.Height = 44
	layer.Visible = true
	layer.Opacity = 1
	layer.Type = "tilelayer"
	layer.Name = "layer1"
	tileMap.Layers = append(tileMap.Layers, layer)

	var tileset phaserTileSet
	tileset.FirstGid = 1
	tileset.Image = "someimg.png"
	tileset.ImageHeight = 128
	tileset.ImageWidth = 176
	tileset.Name = "some tileset name"
	tileset.TileWidth = 16
	tileset.TileHeight = 16
	tileMap.TileSets = append(tileMap.TileSets, tileset)

	c.JSON(http.StatusOK, tileMap)
}
