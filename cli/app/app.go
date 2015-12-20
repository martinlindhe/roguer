package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/martinlindhe/rogue"
	"github.com/martinlindhe/rogue/views"
	"github.com/plimble/ace"
)

func main() {

	/*
		for i := 0; i < 10; i++ {
			island.Tick()
		}
	*/

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

	log.Printf("Starting http server on %s\n", listenAt)

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

	fmt.Printf("Creating new island %s with seed %d\n", newIsland.Name, newIsland.Seed)

	//seed := time.Now().Unix()

	log.Printf("Generating island with seed %d ...\n", newIsland.Seed)
	island := rogue.GenerateIsland(newIsland.Seed, 220, 140)
	island.FillWithCritters()
	log.Println("Done generating island")

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
