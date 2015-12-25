package rogue

import (
	"io/ioutil"
	"math/rand"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

type itemList struct {
	All []Item `json:"all"`
}

// Item ...
type Item struct {
	Age    int
	Name   string `json:"name"`
	Type   string `json:"type"`
	Energy int    `json:"energy"`
}

func (i *Island) randomItemOfType(t string) Item {
	var m []Item

	for _, it := range i.ItemSpecs {
		if it.Type == t {
			m = append(m, it)
		}
	}

	return m[rand.Intn(len(m))]
}

func parseItemsDefinition(defFileName string) []Item {

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		panic(err)
	}

	var items itemList
	err = yaml.Unmarshal(data, &items)
	if err != nil {
		panic(err)
	}

	log.Infof("Read %d entries from %s", len(items.All), defFileName)
	return items.All
}
