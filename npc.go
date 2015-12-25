package rogue

import (
	"io/ioutil"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
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
	Inventory      []Item

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

	var npcs npcListYaml
	err = yaml.Unmarshal(data, &npcs)
	if err != nil {
		panic(err)
	}

	//spew.Dump(npcList)
	log.Infof("Read %d entries from %s", len(npcs.All), defFileName)

	return npcs.All
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
