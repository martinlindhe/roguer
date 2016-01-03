package main

import (
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
var islandMap []byte

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("rogue started")

	r := getRouter()
	island = rogue.NewIsland()
	islandMap = rogue.PrecalcTilemap()

	// listen and serve on 0.0.0.0:3322
	appPort := 3322
	listenAt := fmt.Sprintf(":%d", appPort)

	log.Infof("Starting http server on %s", listenAt)

	go r.Run(listenAt)

	c := time.Tick(3 * time.Second)
	for range c {
		// progress game world
		island.Tick()
	}
}

func getRouter() *ace.Ace {

	// ace with Logger, Recovery
	r := ace.Default()

	//	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/", func(c *ace.C) {
		c.String(200, views.Index())
	})

	r.GET("/island/full", getFullIslandController)

	r.GET("/sprite/character", getTexturePackCharacterController)
	r.GET("/sprite/item", getTexturePackItemController)
	r.GET("/sprite/ground2", getTexturePackGround2Controller)

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

	c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(islandMap)
}

func getTexturePackCharacterController(c *ace.C) {

	ss, err := rogue.ParseSpritesetDefinition("resources/assets/tilesets/oddball/characters.yml")
	if err != nil {
		panic(err)
	}

	tp := rogue.GenerateTexturePacker(ss)
	c.JSON(http.StatusOK, tp)
}

func getTexturePackItemController(c *ace.C) {

	ss, err := rogue.ParseSpritesetDefinition("resources/assets/tilesets/oddball/items.yml")
	if err != nil {
		panic(err)
	}

	tp := rogue.GenerateTexturePacker(ss)
	c.JSON(http.StatusOK, tp)
}

func getTexturePackGround2Controller(c *ace.C) {

	ss, err := rogue.ParseSpritesetDefinition("resources/assets/tilesets/oddball/ground2.yml")
	if err != nil {
		panic(err)
	}

	tp := rogue.GenerateTexturePacker(ss)
	c.JSON(http.StatusOK, tp)
}
