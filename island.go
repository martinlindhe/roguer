package rogue

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/martinlindhe/rogue/rollingparticle"
	"github.com/ojrac/opensimplex-go"
)

// Island ...
type Island struct {
	Width     int
	Height    int
	Seed      int64
	Age       int64
	HeightMap [][]uint
	Spawns    []*Npc

	// lookup lists:
	ItemSpecs []Item
	npcSpecs  []npcSpec
}

var island Island // singelton

// InitIsland inits the singelton
func InitIsland() {
	// XXX load existing world from disk
	seed := int64(666666)
	log.Infof("Generating island with seed %d ...", seed)
	generateIsland(seed, 220, 140)

	island.fillWithCritters()
	log.Info("Done generating island")

	// store island to disk as png
	islandColImage := island.ColoredHeightMapAsImage()
	islandColImageName := fmt.Sprintf("./public/img/islands/%d.png", seed)
	islandColImgFile, _ := os.Create(islandColImageName)
	png.Encode(islandColImgFile, islandColImage)
	/*
		islandImage := island.HeightMapAsImage()
		islandImgFile, _ := os.Create("island.png")
		png.Encode(islandImgFile, islandImage)
	*/
}

// Add ...
func (i *Island) addSpawn(o *Npc) {
	i.Spawns = append(i.Spawns, o)
}

// Tick executes one tick on each spawn in the zone
func (i *Island) Tick() {

	i.Age++
	log.Printf("World tick %d", i.Age)

	for _, o := range i.Spawns {
		o.Tick()
	}
}

// generate critters based on data file
func (i *Island) fillWithCritters() {

	// log.Infof("Looking at %d blueprints for npcs", len(island.npcSpecs))

	for _, npcSpec := range island.npcSpecs {
		log.Infof("Adding %d %s", npcSpec.Quantity, npcSpec.Type)
		for n := 0; n < npcSpec.Quantity; n++ {
			o := new(Npc)

			if len(npcSpec.Name) == 0 {
				// if name field is unset, run a generator based on npc type
				o.Name = generateNpcName(npcSpec.Type)

			} else {
				// pick one name by random
				o.Name = npcSpec.Name[rand.Intn(len(npcSpec.Name))]
			}

			o.Level = 1
			o.Type = npcSpec.Type
			o.Position = i.randomPointAboveWater()
			i.addSpawn(o)
		}
	}
}

func (i *Island) randomPointAboveWater() Point {

	p := Point{uint16(rand.Intn(i.Width)), uint16(rand.Intn(i.Height))}

	// above ground
	if i.HeightMap[p.Y][p.X] > shallowWater {
		return p
	}

	return i.randomPointAboveWater()
}

// GenerateIsland sets the island singelton
func generateIsland(seed int64, width int, height int) {

	particleLength := 8
	innerBlur := 0.85
	outerBlur := 0.60
	roller := rollingparticle.New(seed, width, height, particleLength, innerBlur, outerBlur)

	/*
		rollerImage := slice2DAsImage(&roller, width, height)
		rollerImgFile, _ := os.Create("roller.png")
		png.Encode(rollerImgFile, rollerImage)
	*/

	m := make2DUintSlice(width, height)

	noise := opensimplex.NewWithSeed(seed)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			// combine some sizes of noise
			fBig := noise.Eval2(float64(x)*0.005, float64(y)*0.005)
			fMid := noise.Eval2(float64(x)*0.01, float64(y)*0.01)
			fSmall := noise.Eval2(float64(x)*0.02, float64(y)*0.02)
			fMini := noise.Eval2(float64(x)*0.04, float64(y)*0.04)

			f := (fBig + fMid + fSmall + fMini) / 4

			// scale from -1.0-1.0 to 0.0-1.0
			f = (f + 1.0) / 2.0

			// scale from 0.0-1.0 to 0-255
			b := uint(0)
			if f == 1.0 {
				b = 255
			} else {
				b = uint(math.Floor(f * 256.0))
			}

			// combine with rolling particle
			opacity := 0.5
			b = uint((1-opacity)*float64(b) + opacity*float64(roller[y][x]))

			m[y][x] = b
		}
	}

	island = Island{
		Width:     width,
		Height:    height,
		Seed:      seed,
		HeightMap: m}

	// load all possible world items and npcs
	island.ItemSpecs = parseItemsDefinition("data/items.yml")
	island.npcSpecs = parseNpcsDefinition("data/npc.yml")
}

// ...
const (
	deepWater    = 80
	shallowWater = 90
	beach        = 95
)

// ColoredHeightMapAsImage ...
func (i *Island) ColoredHeightMapAsImage() image.Image {

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{i.Width, i.Height}})

	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			b := i.HeightMap[y][x]

			var col color.RGBA
			switch {
			case b <= deepWater:
				col = color.RGBA{0x26, 0x2f, 0x71, 0xff} // deep water

			case b <= shallowWater:
				col = color.RGBA{0x46, 0x4D, 0x85, 0xff} // shallow water

			case b <= beach:
				col = color.RGBA{0xD4, 0xBC, 0x6A, 0xff} // beach

			case b <= 150:
				col = color.RGBA{0x2D, 0x88, 0x2D, 0xff} // grass (green)

			case b <= 230:
				col = color.RGBA{0x00, 0x4E, 0x00, 0xff} // forest (dark green)

			case b <= 240:
				col = color.RGBA{0x4B, 0x2D, 0x12, 0xff} // hills (brown)

			default:
				col = color.RGBA{0xF2, 0xED, 0xE6, 0xff} // gray (mountains)
			}

			img.Set(x, y, col)
		}
	}

	return img
}
