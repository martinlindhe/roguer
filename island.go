package rogue

import (
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	"github.com/martinlindhe/rogue/rollingparticle"
	"github.com/ojrac/opensimplex-go"
)

// Island ...
type Island struct {
	Width     int
	Height    int
	Seed      int64
	HeightMap [][]uint
	Spawns    []WorldObject
	Age       int64
}

// Tick executes one tick on each spawn in the zone
func (i *Island) Tick() {

	//log.Printf("World tick %d", i.Age)
	for _, o := range i.Spawns {
		o.Tick()
	}
	i.Age++
}

// Add ...
func (i *Island) Add(o WorldObject) {
	i.Spawns = append(i.Spawns, o)
}

// PrintSpawns ...
func (i *Island) PrintSpawns() {

	log.Printf("showing %d spawns:", len(i.Spawns))

	for _, sp := range i.Spawns {
		// XXX need to cast to  instance of the object to call .Name, .Pos
		spawn := reflect.ValueOf(sp)
		log.Println(spawn)
		//log.Printf("%s at %s", spawn.Name, spawn.Position)
	}
}

type npcYaml struct {
	All []npcSpecYaml `json:"all"` // Affects YAML field names too.
}

type npcSpecYaml struct {
	Type     string   `json:"type"`
	Name     []string `json:"name"`
	Quantity int      `json:"qty"`
}

// FillWithCritters ...
func (i *Island) FillWithCritters() {

	// XXX parse yamls

	data, err := ioutil.ReadFile("./data/npc.yml")
	if err != nil {
		panic(err)
	}

	var npcList npcYaml
	err = yaml.Unmarshal(data, &npcList)
	if err != nil {
		panic(err)
	}

	spew.Dump(npcList)

	for _, npcSpec := range npcList.All {
		spew.Dump(npcSpec)
	}
	// XXXXX generate critters based on yaml data

	dwarfs := 1
	//log.Infof("Adding %d dwarfs", dwarfs)
	for n := 0; n < dwarfs; n++ {
		var dwarf dwarf
		dwarf.Defaults()
		dwarf.Position = i.randomPointAboveWater()
		i.Add(&dwarf)
	}

	rabbits := 1
	//log.Infof("Adding %d rabbits", rabbits)
	for n := 0; n < rabbits; n++ {
		var rabbit rabbit
		rabbit.Defaults()
		rabbit.Position = i.randomPointAboveWater()
		i.Add(&rabbit)
	}
}

func (i *Island) randomPointAboveWater() Point {

	p := Point{uint16(rand.Intn(i.Width)), uint16(rand.Intn(i.Height))}

	// above ground
	if i.HeightMap[p.Y][p.X] > ShallowWater {
		return p
	}

	return i.randomPointAboveWater()
}

// GenerateIsland returns a new island
func GenerateIsland(seed int64, width int, height int) Island {

	particleLength := 8
	innerBlur := 0.85
	outerBlur := 0.60
	roller := rollingparticle.New(seed, width, height, particleLength, innerBlur, outerBlur)

	/*
		rollerImage := Slice2DAsImage(&roller, width, height)
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

	island := Island{
		Width:     width,
		Height:    height,
		Seed:      seed,
		HeightMap: m}

	return island
}

// ...
const (
	DeepWater    = 80
	ShallowWater = 90
	Beach        = 95
)

// ColoredHeightMapAsImage ...
func (i *Island) ColoredHeightMapAsImage() image.Image {

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{i.Width, i.Height}})

	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			b := i.HeightMap[y][x]

			var col color.RGBA
			switch {
			case b <= DeepWater:
				col = color.RGBA{0x26, 0x2f, 0x71, 0xff} // deep water

			case b <= ShallowWater:
				col = color.RGBA{0x46, 0x4D, 0x85, 0xff} // shallow water

			case b <= Beach:
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
