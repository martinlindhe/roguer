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

func getFullIslandController(c *ace.C) {
	// NOTE: this is useful in early stage for world debugging.
	// later on, the game would only expose a small area around the player

	// return width, height, heigthmap only
	islandMap := struct {
		Width     int
		Height    int
		HeightMap [][]uint
	}{island.Width, island.Height, island.HeightMap}

	c.JSON(http.StatusOK, islandMap)
}
