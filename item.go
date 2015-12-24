package rogue

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

type item struct {
	Age    int
	Name   string `json:"name"`
	Type   string `json:"type"`
	Energy int    `json:"energy"`
}

type itemList struct {
	All []item `json:"all"` // Affects YAML field names too.
}

func getItemsFromDefinition(defFileName string) []item {

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

	var res []item
	// generate critters based on yaml data
	for _, itemSpec := range items.All {
		log.Infof("Adding %s: %s", itemSpec.Type, itemSpec.Name)
		var o item

		o.Name = itemSpec.Name
		o.Type = itemSpec.Type

		res = append(res, o)
	}

	return res
}
