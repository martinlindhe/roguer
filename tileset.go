package rogue

import (
	"encoding/json"
	"image"
	"io/ioutil"
	"os"

	"github.com/gobuild/log"
	"gopkg.in/yaml.v2"
)

type tilesetSpec struct {
	Props        tilesetProps `json:"props"`
	DeepWater    []int        `json:"deepwater"`
	Beach        []int        `json:"beach"`
	Grass        []int        `json:"grass"`
	ShallowWater []int        `json:"shallowwater"`
	Lava         []int        `json:"lava"`
	Dirt         []int        `json:"dirt"`
	Wall         []int        `json:"wall"`
}

type tilesetProps struct {
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	TileMap string `json:"tilemap"`
}

func parseGroundTilesetDefinition(defFileName string) (tilesetSpec, error) {

	var specs tilesetSpec

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		return specs, err
	}

	err = yaml.Unmarshal(data, &specs)
	if err != nil {
		return specs, err
	}

	log.Infof("Read %s", defFileName)
	return specs, nil
}

func PrecalcTilemap() []byte {

	// XXX also contain collision data

	tiles, err := parseGroundTilesetDefinition("resources/assets/tilesets/oddball/ground.yml")
	if err != nil {
		panic(err)
	}

	imgWidth, imgHeight := getImageDimension(tiles.Props.TileMap)

	var tileMap PhaserTileMap
	tileMap.Version = 1
	tileMap.Width = island.Width
	tileMap.Height = island.Height
	tileMap.TileWidth = tiles.Props.Width
	tileMap.TileHeight = tiles.Props.Height
	tileMap.Orientation = "orthogonal"

	var layer PhaserTileLayer
	layer.Data = island.HeightsAsFlatTilemap()
	layer.Width = island.Width
	layer.Height = island.Height
	layer.Visible = true
	layer.Opacity = 1
	layer.Type = "tilelayer"
	layer.Name = "layer1"
	tileMap.Layers = append(tileMap.Layers, layer)

	var tileset PhaserTileSet
	tileset.FirstGid = 0
	// NOTE: need to specify a tile in phaser later, .Name and .Image must be the same value (phaser 2.4.4, dec 2015)
	tileset.Name = "island_tiles"
	tileset.Image = "island_tiles"
	tileset.ImageWidth = imgWidth
	tileset.ImageHeight = imgHeight
	tileset.TileWidth = tiles.Props.Width
	tileset.TileHeight = tiles.Props.Height
	tileMap.TileSets = append(tileMap.TileSets, tileset)

	b, _ := json.Marshal(tileMap)

	return b
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Println(imagePath, err)
	}
	return image.Width, image.Height
}
