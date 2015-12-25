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
			log.Printf("XXX selecting %s for random roll", it.Name)
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

	//spew.Dump(npcList)
	log.Infof("Processing %d entries from %s", len(items.All), defFileName)

	var res []Item

	for _, itemSpec := range items.All {
		// log.Infof("Adding %s: %s", itemSpec.Type, itemSpec.Name)
		var o Item
		o.Name = itemSpec.Name
		o.Type = itemSpec.Type
		res = append(res, o)
	}

	return res
}
