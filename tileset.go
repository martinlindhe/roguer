package rogue

import (
	"io/ioutil"

	"github.com/gobuild/log"
	"gopkg.in/yaml.v2"
)

type tilesetSpec struct {
	Props tilesetProps   `json:"props"`
	Tiles map[string]int `json:"tiles"`
}

type tilesetProps struct {
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	TileMap string `json:"tilemap"`
}

func parseTilesetDefinition(defFileName string) (tilesetSpec, error) {

	var specs tilesetSpec

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
