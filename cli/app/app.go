package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/martinlindhe/rogue"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

func newWindow() (*sdl.Window, *sdl.Renderer) {
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		panic(err)
	}

	return window, renderer
}

func main() {

	//seed := time.Now().Unix()

	seed := int64(1450549167)

	fmt.Printf("Using seed %d\n", seed)

	island := rogue.GenerateIsland(seed, 220, 100)
	// -------

	// debug output:
	islandImage := island.HeightMapAsImage()
	islandColImage := island.ColoredHeightMapAsImage()

	islandImgFile, _ := os.Create("island.png")
	png.Encode(islandImgFile, islandImage)

	islandColImgFile, _ := os.Create("island_col.png")
	png.Encode(islandColImgFile, islandColImage)

	sdl.Init(sdl.INIT_EVERYTHING)

	window, renderer := newWindow()

	defer window.Destroy()
	defer renderer.Destroy()

	// XXX how to load a sdl.Image from my rgb buffer, rather than from disk?
	image, err := img.Load("island_col.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %s\n", err)
		panic(err)
	}
	defer image.Free()

	texture, err := renderer.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	srcRect := sdl.Rect{0, 0, int32(island.Width), int32(island.Height)}
	dstRect := sdl.Rect{0, 0, int32(island.Width), int32(island.Height)}

	renderer.Clear()
	renderer.Copy(texture, &srcRect, &dstRect)
	renderer.Present()

	sdl.Delay(1000 * 2)

	// -------
	island.FillWithCritters()
	for i := 0; i < 10; i++ {
		island.Tick()
	}

	// XXX could this be deferred instead?
	sdl.Quit()
}
