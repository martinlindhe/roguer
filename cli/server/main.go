package main

import (
	"fmt"
	"time"

	"github.com/martinlindhe/roguer"
	"github.com/nsf/termbox-go"
)

var (
	logMessages MessageList
)

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	//termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	//log.SetLevel(log.DebugLevel)

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

	logMessages.Info("roguer server started, listening on", listenAt)

	go r.Run(listenAt)

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
		logMessages.repaintMostRecent()
	}
}
