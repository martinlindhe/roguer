package rogue

import (
	"io/ioutil"

	"github.com/qiniu/log"
	"gopkg.in/yaml.v2"
)

// Obj is a in-game object, such as a npc or a item
type Obj struct {
	Level    int
	Age      int
	Name     string
	Race     string
	Type     string
	Class    string
	Position Point
	Energy   int
	Weight   int

	XP             int
	CurrentAction  *actionSpec
	PlannedActions []actionSpec
	Inventory      []Obj

	// the lower value, the less hungry npc is
	Hunger    int
	Thirst    int
	Tiredness int
	Coldness  int

	// for objects
	Activated bool
}

type objList struct {
	All []objSpec `json:"all"`
}

type objSpec struct {
	Type   string     `json:"type"`
	Class  string     `json:"class"`
	Race   string     `json:"race"`
	Name   string     `json:"name"`
	Energy int        `json:"energy"`
	Weight int        `json:"weight"`
	Drops  []dropSpec `json:"drops"`
}

type dropSpec struct {
	Name   string  `json:"name"`
	Chance float64 `json:"chance"`
}

func parseObjectsDefinition(defFileName string) []objSpec {

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		panic(err)
	}

	var npcs objList
	err = yaml.Unmarshal(data, &npcs)
	if err != nil {
		panic(err)
	}

	log.Infof("Read %d entries from %s", len(npcs.All), defFileName)
	return npcs.All
}
