package main

import "github.com/martinlindhe/rogue"

func main() {

	seed := int64(123456)

	island := rogue.GenerateIsland(seed, 120, 60)

	island.WriteHeightMapAsImage("test.png")

}
