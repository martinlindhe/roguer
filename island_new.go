package rogue

import (
	"image/png"
	"math"
	"math/rand"
	"os"

	"github.com/martinlindhe/roguer/rollingparticle"
	"github.com/ojrac/opensimplex-go"
)

// NewIsland inits the singelton
func NewIsland() *Island {

	// XXX load existing world from disk
	seed := int64(666666)
	generalLog.Info("Generating island with seed ", seed, " ...")

	island := generateIsland(seed, 220, 140)

	generalLog.Info("Done generating island")

	// store island to disk as png
	islandColImage := island.ColoredHeightMapAsImage()
	//islandColImageName := fmt.Sprintf("public/img/islands/%d.png", seed)
	islandColImageName := "public/img/islands/current.png"
	islandColImgFile, err := os.Create(islandColImageName)
	if err != nil {
		panic(err)
	}
	png.Encode(islandColImgFile, islandColImage)

	return island
}

// create some small rocks spread out over surface
func (i *Island) spawnGravel() {
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			pos := Point{X: float64(x), Y: float64(y)}
			if i.isAboveWater(pos) {
				// 10% chance to add a rock
				if rand.Intn(100) < 10 {
					name := ""
					// XXX make it less likely to place a large rock
					switch rand.Intn(3) {
					case 0:
						name = "small rock"
					case 1:
						name = "medium rock"
					case 2:
						name = "large rock"
					default:
						panic("")
					}
					i.addNpcFromName(name, pos)
				}
			}
		}
	}
}

// create trees all over the island
func (i *Island) spawnTrees() {

	cnt := 0
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			pos := Point{X: float64(x), Y: float64(y)}
			if i.isAboveWater(pos) {
				// add one tree in 2% of the time
				if rand.Intn(100) < 2 {
					name := ""
					// XXX make it less likely to place oak trees ?
					switch rand.Intn(3) {
					case 0:
						name = "oak tree"
					case 1:
						name = "apple tree"
					case 2:
						name = "birch tree"
					default:
						panic("")
					}
					i.addNpcFromName(name, pos)
					cnt++
				}
			}
		}
	}
	generalLog.Info("spawned ", cnt, " trees")
}

// generate critters based on data file
func (i *Island) fillWithCritters() {

	dwarf := i.getNpcSpecFromRace("dwarf")
	for n := 0; n < 5; n++ {
		i.addNpcFromSpec(dwarf, i.RandomPointAboveWater())
	}

	rabbit := i.getNpcSpecFromRace("rabbit")
	for n := 0; n < 5; n++ {
		i.addNpcFromSpec(rabbit, i.RandomPointAboveWater())
	}
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

	m := make2DIntSlice(width, height)

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
			b := int(0)
			if f == 1.0 {
				b = 255
			} else {
				b = int(math.Floor(f * 256.0))
			}

			// combine with rolling particle
			opacity := 0.5
			b = int((1-opacity)*float64(b) + opacity*float64(roller[y][x]))

			m[y][x] = b
		}
	}

	i := Island{
		Width:     width,
		Height:    height,
		Seed:      seed,
		HeightMap: m}

	i.LoadSpecs()

	return &i
}

// LoadSpecs loads all possible world items, NPC:s and actions
func (i *Island) LoadSpecs() {
	i.npcSpecs, _ = parseObjectsDefinition("data/objs.yml")
	i.actionSpecs, _ = parseActionsDefinition("data/actions.yml")
}
