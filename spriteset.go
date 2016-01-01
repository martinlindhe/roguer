package rogue

import (
	"io/ioutil"

	"github.com/gobuild/log"
	"gopkg.in/yaml.v2"
)

type spritesheetSpec struct {
	Props tilesetProps   `json:"props"`
	Tiles map[string]int `json:"tiles"`
}

func parseSpritesetDefinition(defFileName string) (spritesheetSpec, error) {

	var specs spritesheetSpec

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		return specs, err
	}

	err = yaml.Unmarshal(data, &specs)
	if err != nil {
		return specs, err
	}

	log.Infof("Read %d entries from %s", len(specs.Tiles), defFileName)
	return specs, nil
}

// TexturePack ...
type TexturePack struct {
	Frames []textureFrame `json:"frames"`
	Meta   textureMeta    `json:"meta"`
}

type textureMeta struct {
	App         string      `json:"app"`
	Version     string      `json:"version"`
	Image       string      `json:"image"`
	Format      string      `json:"format"`
	Size        textureSize `json:"size"`
	Scale       float64     `json:"scale"`
	SmartUpdate string      `json:"smartupdate"`
}

type textureSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type textureDimension struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type textureFrame struct {
	FileName         string           `json:"filename"`
	Frame            textureDimension `json:"frame"`
	Rotated          bool             `json:"rotated"`
	Trimmed          bool             `json:"trimmed"`
	SpriteSourceSize textureDimension `json:"spriteSourceSize"`
	SourceSize       textureSize      `json:"sourceSize"`
}

func generateTexturePacker(spec spritesheetSpec) TexturePack {

	imgWidth, imgHeight := getImageDimension(spec.Props.TileMap)

	tilesPerRow := imgWidth / spec.Props.Width

	var texturePack TexturePack

	texturePack.Meta.App = "rogue"
	texturePack.Meta.Version = "1.0"
	texturePack.Meta.Image = spec.Props.TileMap
	texturePack.Meta.Format = "RGBA8888"
	texturePack.Meta.Size = textureSize{imgWidth, imgHeight}
	texturePack.Meta.Scale = 1

	for name, tileIdx := range spec.Tiles {

		x := (tileIdx % tilesPerRow) * spec.Props.Width
		y := (tileIdx / tilesPerRow) * spec.Props.Height

		var frame textureFrame
		frame.FileName = name
		frame.Frame = textureDimension{x, y, spec.Props.Width, spec.Props.Height}
		frame.SpriteSourceSize = textureDimension{0, 0, spec.Props.Width, spec.Props.Height}
		frame.SourceSize = textureSize{spec.Props.Width, spec.Props.Height}
		texturePack.Frames = append(texturePack.Frames, frame)
	}

	return texturePack
}
