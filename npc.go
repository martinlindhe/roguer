package rogue

import (
	"io/ioutil"
	"reflect"

	"github.com/qiniu/log"
	"gopkg.in/yaml.v2"
)

// Point ...
type Point struct {
	X uint16
	Y uint16
}

// WorldObject ...
type WorldObject interface {
	Tick()
}

// WorldObjectInstance ...
type WorldObjectInstance struct {
	Level    int
	Age      int
	Name     string
	Type     string
	Position Point
}

// Npc ...
type Npc struct {
	WorldObjectInstance
	XP             int
	CurrentAction  Action
	PlannedActions []Action
	Inventory      []WorldObject

	// the lower value, the less hungry npc is
	Hunger    int
	Thirst    int
	Tiredness int
}

type npcListYaml struct {
	All []npcSpecYaml `json:"all"` // Affects YAML field names too.
}

type npcSpecYaml struct {
	Type     string   `json:"type"`
	Name     []string `json:"name"`
	Quantity int      `json:"qty"`
}

func getNpcsFromDefinition(defFileName string) []npcSpecYaml {

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		panic(err)
	}

	var npcList npcListYaml
	err = yaml.Unmarshal(data, &npcList)
	if err != nil {
		panic(err)
	}

	//spew.Dump(npcList)
	log.Infof("Processing %d entries from %s", len(npcList.All), defFileName)

	return npcList.All
}

type plant struct {
	WorldObjectInstance
}

type edible struct {
	Energy int
}

type sweetPotato struct {
	plant
	edible
}

type dwarf struct {
	Npc
}

type rabbit struct {
	Npc
}

func (n *plant) Tick() {
	n.Age++
}

func (n *Npc) hungerCap() int {
	return n.Level * 5
}

func (n *Npc) thirstCap() int {
	return n.Level * 100
}

func (n *Npc) tirednessCap() int {
	return n.Level * 5
}

func (n *Npc) pickSomethingToEat() WorldObject {
	// XXX find something edible in inventory, or nil
	if len(n.Inventory) == 0 {
		return nil
	}

	return n.Inventory[0]
}

// check if npc already has planned to do a
func (n *Npc) hasPlanned(a Action) bool {

	t := reflect.TypeOf(a)

	if reflect.TypeOf(n.CurrentAction) == t {
		return true
	}

	for _, v := range n.PlannedActions {
		if reflect.TypeOf(v) == t {
			return true
		}
	}
	return false
}
