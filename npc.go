package rogue

import (
	"io/ioutil"

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

// Npc ...
type Npc struct {
	Level    int
	Age      int
	Name     string
	Type     string
	Position Point

	XP                       int
	TimeSpentOnCurrentAction int
	CurrentAction            string
	PlannedActions           []string
	Inventory                []Item

	// the lower value, the less hungry npc is
	Hunger    int
	Thirst    int
	Tiredness int
}

type npcList struct {
	All []npcSpec `json:"all"`
}

type npcSpec struct {
	Type     string   `json:"type"`
	Name     []string `json:"name"`
	Quantity int      `json:"qty"`
}

func (i *Island) withinRadius(n string, radius int, pos Point) []Npc {
	// XXXX
	return nil
}

func parseNpcsDefinition(defFileName string) []npcSpec {

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		panic(err)
	}

	var npcs npcList
	err = yaml.Unmarshal(data, &npcs)
	if err != nil {
		panic(err)
	}

	log.Infof("Read %d entries from %s", len(npcs.All), defFileName)
	return npcs.All
}

// check if npc already has planned to do a
func (n *Npc) hasPlanned(t string) bool {

	if n.CurrentAction == t {
		return true
	}

	for _, v := range n.PlannedActions {
		if v == t {
			return true
		}
	}
	return false
}

func (n *Npc) planAction(t string) {
	n.PlannedActions = append(n.PlannedActions, t)
}
