package rogue

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/martinlindhe/imgcat/lib"
	"github.com/nsf/termbox-go"
	"github.com/plimble/ace"
)

var (
	islandMap        []byte
	appPort          = 3322
	mainloopInterval = 100 * time.Millisecond
	gameTickIRL      = 3 * time.Second // 1 game tick = 3 real world seconds
	snapshotInterval = 10*gameTickIRL + 1
	mongoSession     *mgo.Session
	dbName           = "roguer"
	enableAutosave   = false
)

// BootServer ...
func BootServer() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	//termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	mongoSession, err := initServer()
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()

	serverLoop()
}

func initServer() (*mgo.Session, error) {
	mongoSession, err := initMongo()

	registerAutosaver()

	newOrResumeIsland()

	r := getHTTPRouter()
	islandMap = PrecalcTilemap()
	listenAt := fmt.Sprintf(":%d", appPort)

	go r.Run(listenAt)

	return mongoSession, err
}

func serverLoop() {

	// main loop
	var cnt time.Duration
	c := time.Tick(mainloopInterval)
	for range c {

		if !handleEvents() {
			break
		}

		cnt += mainloopInterval
		if cnt >= gameTickIRL {
			cnt = 0
			// progress game world
			island.Tick()
		}
		generalLog.repaintMostRecent()
	}
}

func newOrResumeIsland() {
	NewIsland()

	imgcat.CatFile("public/img/islands/current.png", os.Stdout)

	// XXX loading island half-works, disabled for now //jan 2016
	/*
		island.Seed = 666666 // XXX dont hard code
		fmt.Printf("Resuming island with seed %d\n", island.Seed)
		err = db.Find(bson.M{"_id": island.Seed}).One(&island)
		if err != nil {
			//panic(err)
			fmt.Printf("ERROR resuming, creating new world")
			NewIsland()
		} else {
			island.LoadSpecs()
		}
	*/
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

				generalLog.Info("---SAVE at", time.Now())

				mongoSession.Refresh()

				coll := mongoSession.DB(dbName).C("world")

				_, err := coll.UpsertId(island.Seed, island)
				if err != nil {
					generalLog.Info("ERROR saving db:", err)
					mongoSession.Refresh()

					_, err = coll.UpsertId(island.Seed, &island)
					if err != nil {
						generalLog.Info("FATAL ERROR, failed twice saving db")
					}
				}

				generalLog.Info("---DONE at", time.Now())

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func getHTTPRouter() *ace.Ace {

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
		serveWebsocket(c.Writer, c.Request)
	})

	r.Static("/js", "./public/js")
	r.Static("/css", "./public/css")
	r.Static("/fonts", "./public/fonts")
	r.Static("/img", "./public/img")
	r.Static("/audio", "./public/audio")

	generalLog.Info("http server started, listening on port", appPort)

	return r
}

// returns a map in Tiled json format, recognized by phaser.io
func getFullIslandController(c *ace.C) {

	c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(islandMap)
}

func getTexturePackCharacterController(c *ace.C) {

	ss, err := ParseSpritesetDefinition("resources/assets/tilesets/oddball/characters.yml")
	if err != nil {
		panic(err)
	}

	tp := GenerateTexturePacker(ss)
	c.JSON(http.StatusOK, tp)
}

func getTexturePackItemController(c *ace.C) {

	ss, err := ParseSpritesetDefinition("resources/assets/tilesets/oddball/items.yml")
	if err != nil {
		panic(err)
	}

	tp := GenerateTexturePacker(ss)
	c.JSON(http.StatusOK, tp)
}

func getTexturePackGround2Controller(c *ace.C) {

	ss, err := ParseSpritesetDefinition("resources/assets/tilesets/oddball/ground2.yml")
	if err != nil {
		panic(err)
	}

	tp := GenerateTexturePacker(ss)
	c.JSON(http.StatusOK, tp)
}
