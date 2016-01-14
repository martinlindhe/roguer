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
	island           rogue.Island
	islandMap        []byte
	appPort          = 3322
	tickDuration     = 3 * time.Second // 1 game tick = 3 real world seconds
	snapshotInterval = 10*tickDuration + 1
	mongoSession     *mgo.Session
	dbName           = "roguer"
	enableAutosave   = false
)

func newOrResumeIsland() {
	rogue.NewIsland()

	// XXX loading island half-works, disabled for now //jan 2016
	/*
		island.Seed = 666666 // XXX dont hard code
		fmt.Printf("Resuming island with seed %d\n", island.Seed)
		err = db.Find(bson.M{"_id": island.Seed}).One(&island)
		if err != nil {
			//panic(err)
			fmt.Printf("ERROR resuming, creating new world")
			rogue.NewIsland()
		} else {
			island.LoadSpecs()
		}
	*/
}

func main() {

	log.SetLevel(log.DebugLevel)

	mongoSession, err := initMongo()
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()

	newOrResumeIsland()

	registerAutosaver()

	r := getRouter()

	islandMap = rogue.PrecalcTilemap()

	listenAt := fmt.Sprintf(":%d", appPort)

	log.Infof("roguer server started, listening on %s", listenAt)

	go r.Run(listenAt)

	c := time.Tick(tickDuration)
	for range c {
		// progress game world
		island.Tick()
	}
}

func initMongo() (*mgo.Session, error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	return session, nil
}

func registerAutosaver() {
	ticker := time.NewTicker(snapshotInterval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if !enableAutosave {
					break
				}

				log.Printf("---SAVE at %s\n", time.Now())

				mongoSession.Refresh()

				coll := mongoSession.DB(dbName).C("world")

				_, err := coll.UpsertId(island.Seed, island)
				if err != nil {
					log.Fatalf("ERROR saving db: %s", err)
					mongoSession.Refresh()

					_, err = coll.UpsertId(island.Seed, &island)
					if err != nil {
						log.Fatalf("FATAL ERROR, failed twice saving db\n")
					}
				}

				log.Printf("---DONE at %s\n", time.Now())

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
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
