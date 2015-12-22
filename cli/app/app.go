package main

import (
	"fmt"
	"image/png"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/martinlindhe/rogue"
	"github.com/martinlindhe/rogue/views"
	"github.com/plimble/ace"
)

func main() {

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	log.Info("rogue started")

	log.Debug("debug msg")

	r := getRouter()

	// r.GET("/", views.Index())
	r.GET("/", func(c *ace.C) {
		//c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, views.Index())
	})

	// XXX run http server in separate process?

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

	r.POST("/island/new", newIslandController)

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

func newIslandController(c *ace.C) {

	newIsland := struct {
		Name string `json:"name"`
		Seed int64  `json:"seed"`
	}{}

	c.ParseJSON(&newIsland)

	log.Infof("Generating island %s with seed %d ...", newIsland.Name, newIsland.Seed)
	island := rogue.GenerateIsland(newIsland.Seed, 220, 140)
	island.FillWithCritters()
	log.Info("Done generating island")

	// store island to disk as png
	islandColImage := island.ColoredHeightMapAsImage()
	islandColImageName := fmt.Sprintf("./public/img/islands/%d.png", newIsland.Seed)
	islandColImgFile, _ := os.Create(islandColImageName)
	png.Encode(islandColImgFile, islandColImage)
	/*
		islandImage := island.HeightMapAsImage()
		islandImgFile, _ := os.Create("island.png")
		png.Encode(islandImgFile, islandImage)
	*/
	// XXX return as json
	c.JSON(http.StatusOK, island)
}
