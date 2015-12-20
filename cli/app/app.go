package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/martinlindhe/rogue"
	"github.com/martinlindhe/rogue/views"
)

func main() {

	/*
		for i := 0; i < 10; i++ {
			island.Tick()
		}
	*/

	r := getRouter()

	// r.GET("/", views.Index())
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, views.Index())
	})

	// XXX run http server in separate process?

	// listen and serve on 0.0.0.0:8080
	appPort := 3322
	listenAt := fmt.Sprintf(":%d", appPort)

	log.Printf("Starting http server on %s\n", listenAt)

	r.Run(listenAt)
}

func getRouter() *gin.Engine {

	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())

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
func pingController(c *gin.Context) {
	c.JSON(200, gin.H{"pong": "now"})
}

// Binding from JSON
type NewIsland struct {
	Name string `json:"name" binding:"required"`
	Seed string `json:"seed" binding:"required"`
}

func newIslandController(c *gin.Context) {

	var jsonFormat NewIsland
	if c.BindJSON(&jsonFormat) != nil {
		fmt.Println("error decoding json")
		return
	}

	seed, err := strconv.ParseInt(jsonFormat.Seed, 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Creating new island %s with seed %s\n", jsonFormat.Name, jsonFormat.Seed)

	//seed := time.Now().Unix()

	log.Printf("Generating island with seed %d ...\n", seed)
	island := rogue.GenerateIsland(seed, 220, 140)
	island.FillWithCritters()
	log.Println("Done generating island")

	islandColImage := island.ColoredHeightMapAsImage()

	islandColImageName := fmt.Sprintf("./public/img/islands/%d.png", seed)
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
