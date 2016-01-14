package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/martinlindhe/roguer"
	"github.com/plimble/ace"
)

var (
	island       rogue.Island
	islandMap    []byte
	appPort      = 3322
	tickDuration = 3 * time.Second // 1 game tick = 3 real world seconds
)

func main() {

	log.SetLevel(log.DebugLevel)

	mongo, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer mongo.Close()

	mongo.SetMode(mgo.Monotonic, true)

	db := mongo.DB("test").C("roguer")
	/*
		//XXX CRASH
		//island := make(*rogue.Island)
		err = db.Find(bson.M{"Seed": island.Seed}).All(&island)
		if err != nil {
			panic(err)
		}
	*/
	ticker := time.NewTicker(10*tickDuration + 1)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:

				log.Printf("-SNAPSHOTTING DB at %s\n", time.Now())

				mongo.Refresh()
				_, err = db.UpsertId(island.Seed, &island)
				if err != nil {
					log.Fatalf("ERROR saving db: %s", err)
					mongo.Refresh()

					_, err = db.UpsertId(island.Seed, &island)
					if err != nil {
						log.Fatalf("FATAL ERROR, failed twice saving db\n")
					}
				}

				log.Printf("-DONE at %s\n", time.Now())

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	r := getRouter()
	rogue.NewIsland()
	islandMap = rogue.PrecalcTilemap()

	listenAt := fmt.Sprintf(":%d", appPort)

	log.Infof("roguer server started, listening on %s", listenAt)

	go r.Run(listenAt)

	// initial tick
	island.Tick()

	c := time.Tick(tickDuration)
	for range c {
		// progress game world
		island.Tick()
	}
}

func getRouter() *ace.Ace {

	// ace with Logger, Recovery
	r := ace.Default()

	r.GET("/", func(c *ace.C) {
		body, _ := ioutil.ReadFile("views/index.html")
		c.String(200, string(body))
	})

	r.GET("/island/full", getFullIslandController)

	r.GET("/sprite/character", getTexturePackCharacterController)
	r.GET("/sprite/item", getTexturePackItemController)
	r.GET("/sprite/ground2", getTexturePackGround2Controller)

	r.GET("/ws", func(c *ace.C) {
		rogue.ServeWs(c.Writer, c.Request)
	})

	r.Static("/js", "./public/js")
	r.Static("/css", "./public/css")
	r.Static("/fonts", "./public/fonts")
	r.Static("/img", "./public/img")
	r.Static("/audio", "./public/audio")

	return r
}

// returns a map in Tiled json format, recognized by phaser.io
func getFullIslandController(c *ace.C) {

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
