package rogue

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/gorilla/websocket"

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
	shallowWater = 104
	beach        = 112
	grass        = 118
)

// NewPlayer creates a new player
func (i *Island) NewPlayer(name string, socket *websocket.Conn) (Point, string) {

	pos := island.RandomPointAboveWater()
	token := newJwt()

	spawn := i.getNpcFromRace("dwarf")
	spawn.Class = "player"
	spawn.Name = name
	spawn.Position = pos
	island.Spawns = append(island.Spawns, spawn)

	var player Player
	player.Name = name
	player.Token = token
	player.Spawn = spawn
	player.Socket = socket
	island.Players = append(island.Players, player)

	return pos, token
}

// ContinuePlayer resumes a game, given a valid token
func (i *Island) ContinuePlayer(token string, socket *websocket.Conn) (*Point, string, error) {

	for i, pl := range island.Players {
		if pl.Token == token {
			// give the client a fresh token
			island.Players[i].Token = newJwt()
			island.Players[i].Socket = socket
			return &pl.Spawn.Position, island.Players[i].Token, nil
		}
	}

	return nil, token, fmt.Errorf("Session not found")
}

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

	return i.addNpcFromSpec(i.getNpcSpecFromName(n), pos)
}

func (i *Island) addNpcFromRace(n string, pos Point) {

	i.addNpcFromSpec(i.getNpcSpecFromRace(n), pos)
}

func (i *Island) getNpcFromRace(race string) *Obj {
	spec := i.getNpcSpecFromRace(race)
	return i.getNpcFromSpec(spec)
}

func (i *Island) getNpcFromSpec(spec objSpec) *Obj {
	o := new(Obj)

	o.Level = 1
	o.Race = spec.Race
	o.Type = spec.Type
	o.Class = spec.Class
	o.Sprite = spec.Sprite
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

// RandomPointAboveWater ...
func (i *Island) RandomPointAboveWater() Point {

	p := Point{float64(rand.Intn(i.Width)), float64(rand.Intn(i.Height))}

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

			default:
				col = color.RGBA{0x2D, 0x88, 0x2D, 0xff} // grass (green)
			}

			img.Set(x, y, col)
		}
	}

	return img
}

// HeightsAsFlatTilemap ...
func (i *Island) HeightsAsFlatTilemap() []int {

	tiles, err := parseGroundTilesetDefinition("resources/assets/tilesets/oddball/ground.yml")
	if err != nil {
		panic(err)
	}

	res := make([]int, island.Width*island.Height)

	for y := 0; y < island.Height; y++ {
		for x := 0; x < island.Width; x++ {
			num := 0

			b := island.HeightMap[y][x]

			switch {
			case b < deepWater:
				num = fromSliceByScale(b, 0, deepWater, tiles.DeepWater)

			case b < shallowWater:
				num = fromSliceByScale(b, deepWater, shallowWater, tiles.ShallowWater)

			case b < beach:
				num = fromSliceByScale(b, shallowWater, beach, tiles.Beach)

			default:
				num = fromSliceByScale(b, beach, grass, tiles.Grass)
			}

			res[y*island.Width+x] = num
		}
	}

	return res
}

func scale(valueIn float64, baseMin float64, baseMax float64, limitMin float64, limitMax float64) float64 {
	return ((limitMax - limitMin) * (valueIn - baseMin) / (baseMax - baseMin)) + limitMin
}

func fromSliceByScale(b int, min int, max int, tiles []int) int {

	num := int(scale(float64(b), float64(min), float64(max), 0, float64(len(tiles))))

	if num >= len(tiles) {
		fmt.Printf("ERROR: num is %d, max is %d\n", num, len(tiles)-1)
		num = len(tiles) - 1
	}
	return tiles[num]
}

// LocalSpawns exposes public info about the spawn to the player
type LocalSpawns struct {
	Name   string
	Class  string
	Sprite string
	X      float64
	Y      float64
}

// DescribeLocalArea returns all spawns near pos
func (i *Island) DescribeLocalArea(pos Point) []LocalSpawns {
	var res []LocalSpawns

	for _, sp := range island.Spawns {
		if sp.Position.isNearby(pos) {
			var ls LocalSpawns
			ls.Name = sp.Name
			ls.Class = sp.Class
			ls.Sprite = sp.Sprite
			ls.X = sp.Position.X
			ls.Y = sp.Position.Y

			res = append(res, ls)
		}
	}

	return res
}
