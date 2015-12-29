package rogue

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	log "github.com/Sirupsen/logrus"
)

// Island ...
type Island struct {
	Width     int
	Height    int
	Seed      int64
	Age       int64
	HeightMap [][]int
	Spawns    []*Obj
	Players   []Player

	// lookup lists:
	npcSpecs    []objSpec
	actionSpecs []actionSpec
}

// height constants
const (
	deepWater    = 80
	shallowWater = 90
	beach        = 95
	grass        = 150
	forest       = 230
	hills        = 240
)

// Add ...
func (i *Island) addSpawn(o *Obj) {
	i.Spawns = append(i.Spawns, o)
}

func (i *Island) removeSpawn(o *Obj) {

	removeIdx := -1

	for idx, sp := range i.Spawns {
		if sp == o {
			removeIdx = idx
			break
		}
	}

	if removeIdx == -1 {
		panic("removeSpawn failed")
	}

	i.Spawns = append(i.Spawns[:removeIdx], i.Spawns[removeIdx+1:]...)
}

// Tick executes one tick on each spawn in the zone
func (i *Island) Tick() {

	i.Age++
	log.Debugf("World tick %d", i.Age)

	for _, o := range i.Spawns {
		check := o.Tick()
		if check == false {
			log.Infof("Removing spawn %s", o.Name)
			i.removeSpawn(o)
		}
	}
}

func (i *Island) getNpcSpecFromName(n string) objSpec {
	for _, spec := range island.npcSpecs {
		if spec.Name == n {
			return spec
		}
	}

	panic(fmt.Errorf("npc spec by name not found: %s", n))
}

func (i *Island) getNpcSpecFromRace(n string) objSpec {
	for _, spec := range island.npcSpecs {
		if spec.Race == n {
			return spec
		}
	}

	panic(fmt.Errorf("npc spec by race not found: %s", n))
}

func (i *Island) addNpcFromName(n string, pos Point) *Obj {

	return island.addNpcFromSpec(island.getNpcSpecFromName(n), pos)
}

func (i *Island) addNpcFromRace(n string, pos Point) {

	island.addNpcFromSpec(island.getNpcSpecFromRace(n), pos)
}

func (i *Island) getNpcFromSpec(spec objSpec) *Obj {
	o := new(Obj)

	o.Level = 1
	o.Race = spec.Race
	o.Type = spec.Type
	o.Class = spec.Class
	o.Energy = spec.Energy
	o.Weight = spec.Weight

	if spec.Name == "" {
		// if name field is unset, let the npc generate a name
		o.Name = o.generateName()
	} else {
		o.Name = spec.Name
	}

	return o
}

func (i *Island) addNpcFromSpec(spec objSpec, pos Point) *Obj {

	o := i.getNpcFromSpec(spec)
	o.Position = pos
	i.addSpawn(o)
	return o
}

func (i *Island) RandomPointAboveWater() Point {

	p := Point{float64(rand.Intn(i.Width)), float64(rand.Intn(i.Height))}

	// above ground
	if i.isAboveWater(p) {
		return p
	}

	return i.RandomPointAboveWater()
}

func (i *Island) isAboveWater(p Point) bool {
	if i.HeightMap[int(p.Y)][int(p.X)] > shallowWater {
		return true
	}
	return false
}

func (i *Island) withinRadiusOfName(n string, radius float64, pos Point) []Obj {

	var res []Obj
	for _, npc := range i.Spawns {
		if npc.Name == n && npc.distanceTo(pos) <= radius {
			res = append(res, *npc)
		}
	}
	return res
}

func (i *Island) withinRadiusOfType(t string, radius float64, pos Point) []*Obj {

	var res []*Obj
	for _, npc := range i.Spawns {

		// log.Printf("XXX %s at dist %f (radius %f)", npc, npc.distanceTo(pos), radius)
		if npc.Type == t && npc.distanceTo(pos) <= radius {

			res = append(res, npc)
		}
	}
	return res
}

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

			case b <= grass:
				col = color.RGBA{0x2D, 0x88, 0x2D, 0xff} // grass (green)

			case b <= forest:
				col = color.RGBA{0x00, 0x4E, 0x00, 0xff} // forest (dark green)

			case b <= hills:
				col = color.RGBA{0x4B, 0x2D, 0x12, 0xff} // hills (brown)

			default:
				col = color.RGBA{0xF2, 0xED, 0xE6, 0xff} // gray (mountains)
			}

			img.Set(x, y, col)
		}
	}

	return img
}

func mapHeightToTileNumber(b int) int {
	switch {
	case b <= deepWater:
		return 0 // deep water

	case b <= shallowWater:
		return 1 // shallow water

	case b <= beach:
		return 2 // beach

	case b <= grass:
		return 3 // grass (green)

	case b <= forest:
		return 4 // forest (dark green)

	case b <= hills:
		return 5 // hills (brown)

	default:
		return 6 // gray (mountains)
	}
}

func (i *Island) HeightsAsFlatTilemap() []int {
	res := make([]int, island.Width*island.Height)

	for y := 0; y < island.Height; y++ {
		for x := 0; x < island.Width; x++ {
			res[y*island.Width+x] = mapHeightToTileNumber(island.HeightMap[y][x])
		}
	}

	return res
}

// this describes a little piece of the world map
type LocalArea struct {
	center Point
	Spawns []spawnInfo
}

// expose "public" info about the spawn to the player
type spawnInfo struct {
	Name string
	X    float64
	Y    float64
}

func (i *Island) DescribeLocalArea(pos Point) LocalArea {
	var res LocalArea

	// find all spawns near pos
	for _, sp := range island.Spawns {
		if sp.Position.isNearby(pos) {
			var si spawnInfo
			si.Name = sp.Name
			si.X = sp.Position.X
			si.Y = sp.Position.Y

			res.Spawns = append(res.Spawns, si)
		}
	}

	return res
}
