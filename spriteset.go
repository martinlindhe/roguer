package rogue

import (
	"io/ioutil"

	"github.com/gobuild/log"
	"gopkg.in/yaml.v2"
)

func parseSpritesetDefinition(defFileName string) (map[string]int, error) {

	var specs map[string]int

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		return specs, err
	}

	err = yaml.Unmarshal(data, &specs)
	if err != nil {
		return specs, err
	}

	log.Infof("Read %d entries from %s", len(specs), defFileName)
	return specs, nil
}
