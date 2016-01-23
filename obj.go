package rogue

import (
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

// Obj is a in-game object, such as a npc or a item
type Obj struct {
	Id       int64
	Level    int
	Age      GameTime
	Name     string
	Race     string
	Type     string
	Class    string
	Sprite   string
	Position Point
	Energy   int
	Weight   int

	XP             int
	Home           *Obj
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

func (o Obj) String() string {
	// XXX
	return fmt.Sprintf("%s (%s)", o.Name, o.Position)
}

type objList struct {
	All []objSpec `json:"all"`
}

type objSpec struct {
	Type   string     `json:"type"`
	Class  string     `json:"class"`
	Race   string     `json:"race"`
	Name   string     `json:"name"`
	Sprite string     `json:"sprite"`
	Energy int        `json:"energy"`
	Weight int        `json:"weight"`
	Drops  []dropSpec `json:"drops"`
}

type dropSpec struct {
	Name   string  `json:"name"`
	Chance float64 `json:"chance"`
}

func parseObjectsDefinition(defFileName string) ([]objSpec, error) {

	var npcs objList

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		return npcs.All, err
	}

	err = yaml.Unmarshal(data, &npcs)
	if err != nil {
		return npcs.All, err
	}

	log.Infof("Read %d entries from %s", len(npcs.All), defFileName)
	return npcs.All, nil
}
