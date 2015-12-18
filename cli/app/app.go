package main

import "github.com/martinlindhe/rogue"

func main() {

	seed := int64(123456)

	isl := rogue.GenerateIsland(seed, 120, 60)

	isl.WriteHeightMapAsImage("test.png")

}
