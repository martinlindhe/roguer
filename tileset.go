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

// PrecalcTilemap ...
func PrecalcTilemap() []byte {

	tiles, err := parseGroundTilesetDefinition("resources/assets/tilesets/oddball/ground.yml")
	if err != nil {
		panic(err)
	}

	imgWidth, imgHeight := getImageDimension(tiles.Props.TileMap)

	tileMap := TiledMapJSON{
		Version:     1,
		Width:       island.Width,
		Height:      island.Height,
		TileWidth:   tiles.Props.Width,
		TileHeight:  tiles.Props.Height,
		Orientation: "orthogonal"}

	layer := TiledMapLayer{
		Data:    island.HeightsAsFlatTilemap(),
		Width:   island.Width,
		Height:  island.Height,
		Visible: true,
		Opacity: 1,
		Type:    "tilelayer",
		Name:    "layer1"}

	tileMap.Layers = append(tileMap.Layers, layer)

	// NOTE: need to specify a tile in phaser later, Name and Image
	// must have same value (phaser 2.4.4, dec 2015)
	tileset := TiledTileSet{
		FirstGid:    0,
		Name:        "island_tiles",
		Image:       "island_tiles",
		ImageWidth:  imgWidth,
		ImageHeight: imgHeight,
		TileWidth:   tiles.Props.Width,
		TileHeight:  tiles.Props.Height}

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
