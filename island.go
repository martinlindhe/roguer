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
	ItemSpecs   []Item
	npcSpecs    []npcSpec
	actionSpecs []actionSpec
}

var island *Island // singelton

// InitIsland inits the singelton
func InitIsland() {
	// XXX load existing world from disk
	seed := int64(666666)
	log.Infof("Generating island with seed %d ...", seed)
	island = generateIsland(seed, 220, 140)

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
	log.Debugf("World tick %d", i.Age)

	for idx, o := range i.Spawns {
		check := o.Tick()
		if check == false {
			log.Infof("Removing spawn %s", o.Name)
			i.Spawns = append(i.Spawns[:idx], i.Spawns[idx+1:]...)
		}
	}
}

// generate critters based on data file
func (i *Island) fillWithCritters() {

	for _, npcSpec := range island.npcSpecs {
		log.Infof("Adding %d %s", npcSpec.Quantity, npcSpec.Type)
		for n := 0; n < npcSpec.Quantity; n++ {
			i.addNpcFromSpec(npcSpec, i.randomPointAboveWater())
		}
	}
}

func (i *Island) getNpcSpecFromName(n string) npcSpec {
	for _, npcSpec := range island.npcSpecs {
		if len(npcSpec.Name) > 0 && npcSpec.Name[0] == n {
			return npcSpec
		}
	}

	panic(fmt.Errorf("npc spec by name not found: %s", n))
}

func (i *Island) getNpcSpecFromRace(n string) npcSpec {
	for _, npcSpec := range island.npcSpecs {
		if npcSpec.Race == n {
			return npcSpec
		}
	}

	panic(fmt.Errorf("npc spec by race not found: %s", n))
}

func (i *Island) addNpcFromName(n string, pos Point) {

	island.addNpcFromSpec(island.getNpcSpecFromName(n), pos)
}

func (i *Island) addNpcFromRace(n string, pos Point) {

	island.addNpcFromSpec(island.getNpcSpecFromRace(n), pos)
}

func (i *Island) addNpcFromSpec(spec npcSpec, pos Point) {
	o := new(Npc)

	o.Level = 1
	o.Race = spec.Race
	o.Type = spec.Type
	o.Class = spec.Class
	o.Position = pos

	if len(spec.Name) == 0 {
		// if name field is unset, run a generator based on npc type
		o.Name = o.generateName()

	} else {
		// pick one name by random
		o.Name = spec.Name[rand.Intn(len(spec.Name))]
	}

	i.addSpawn(o)
}

func (i *Island) randomPointAboveWater() Point {

	p := Point{rand.Intn(i.Width), rand.Intn(i.Height)}

	// above ground
	if i.HeightMap[p.Y][p.X] > shallowWater {
		return p
	}

	return i.randomPointAboveWater()
}

func generateIsland(seed int64, width int, height int) *Island {

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

	is := &Island{
		Width:     width,
		Height:    height,
		Seed:      seed,
		HeightMap: m}

	// load all possible world items, NPC:s and actions
	is.ItemSpecs = parseItemsDefinition("data/items.yml")
	is.npcSpecs = parseNpcsDefinition("data/npc.yml")
	is.actionSpecs = parseActionsDefinition("data/actions.yml")

	return is
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
