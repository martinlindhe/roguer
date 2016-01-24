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
	game             *Game
	appPort          = 3322
	mainloopInterval = 100 * time.Millisecond
	gameTickIRL      = 3 * time.Second // 1 game tick = 3 real world seconds
	snapshotInterval = 10*gameTickIRL + 1
	dbName           = "roguer"
	enableAutosave   = false
)

// Game ...
type Game struct {
	mongoSession *mgo.Session
	input        *input
	Island       *Island
	islandMap    []byte
}

// BootServer ...
func BootServer() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	//termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	game, err = initServer()
	if err != nil {
		panic(err)
	}
	defer game.mongoSession.Close()

	// Init input
	game.input.start()
	defer game.input.stop()

	game.serverLoop()
}

// newGame creates a new Game, along with a input handler.
// Returns a pointer to the new Game.
func newGame() (*Game, error) {
	mongoSession, err := initMongo()
	g := Game{
		input:        newInput(),
		mongoSession: mongoSession,
		Island:       NewIsland(),
	}

	g.Island.spawnGravel()
	g.Island.spawnTrees()
	g.Island.fillWithCritters()

	g.precalcTilemap()

	if err != nil {
		return &g, err
	}

	return &g, nil
}

func initServer() (*Game, error) {

	game, err := newGame()
	if err != nil {
		return game, err
	}

	registerAutosaver()

	r := getHTTPRouter()
	listenAt := fmt.Sprintf(":%d", appPort)

	go r.Run(listenAt)

	return game, nil
}

func (g *Game) serverLoop() {

	imgcat.CatFile("public/img/islands/current.png", os.Stdout)

	// main loop
	var cnt time.Duration
	c := time.Tick(mainloopInterval)
	for range c {

		generalLog.repaintMostRecent()
		/*
			if !handleEvents() {
				break
			}
		*/
		select {
		case ev := <-g.input.eventQ:
			if ev.Key == g.input.endKey {
				fmt.Println("breaking main loop")
				return
			}
		}

		cnt += mainloopInterval
		if cnt >= gameTickIRL {
			cnt = 0
			// progress game world
			g.Island.Tick()
		}
	}
}

/*
func resumeIsland() {
	// XXX loading island half-works, disabled for now //jan 2016

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
}
*/
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
				/*
					generalLog.Info("---SAVE at", time.Now())

					game.mongoSession.Refresh()

					coll := game.mongoSession.DB(dbName).C("world")

					_, err := coll.UpsertId(island.Seed, island)
					if err != nil {
						generalLog.Info("ERROR saving db:", err)
						game.mongoSession.Refresh()

						_, err = coll.UpsertId(island.Seed, &island)
						if err != nil {
							generalLog.Info("FATAL ERROR, failed twice saving db")
						}
					}

					generalLog.Info("---DONE at", time.Now())
				*/
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

	generalLog.Info("http server started, listening on port ", appPort)

	return r
}

// returns a map in Tiled json format, recognized by phaser.io
func getFullIslandController(c *ace.C) {

	c.Writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(game.islandMap)
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
