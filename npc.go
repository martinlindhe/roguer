package rogue

import (
	"io/ioutil"
	"math"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

// Point ...
type Point struct {
	X int
	Y int
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
	Race     string
	Type     string
	Position Point

	XP             int
	CurrentAction  *actionSpec
	PlannedActions []actionSpec
	Inventory      []Item

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
	Race     string   `json:"race"`
	Name     []string `json:"name"`
	Quantity int      `json:"qty"`
}

func (n *Npc) distanceTo(pos Point) float64 {

	xd := float64(n.Position.X - pos.X)
	yd := float64(n.Position.Y - pos.Y)

	return math.Hypot(xd, yd)
}

func (i *Island) withinRadiusOfName(n string, radius float64, pos Point) []Npc {
	var res []Npc

	for _, npc := range i.Spawns {
		if npc.Name == n && npc.distanceTo(pos) <= radius {
			res = append(res, *npc)
		}
	}

	return res
}

func (i *Island) withinRadiusOfType(t string, radius float64, pos Point) []Npc {
	var res []Npc

	for _, npc := range i.Spawns {

		if npc.Type == t && npc.distanceTo(pos) <= radius {
			res = append(res, *npc)
		}
	}

	return res
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

	if n.CurrentAction != nil && n.CurrentAction.Name == t {
		return true
	}

	for _, v := range n.PlannedActions {
		if v.Name == t {
			return true
		}
	}
	return false
}

func (n *Npc) planAction(actionName string) {

	if n.hasPlanned(actionName) {
		return
	}

	a := island.findActionByName(actionName)
	log.Printf("%s decided to %s", n.Name, a.Name)

	n.PlannedActions = append(n.PlannedActions, a)
}

// return false if there is some construction at pos
func (i *Island) canBuildAt(pos Point) bool {

	for _, sp := range i.Spawns {
		if sp.Position == pos {
			if sp.Type == "tree" || sp.Type == "food producer" || sp.Type == "shelter" || sp.Type == "fireplace" {
				return false
			}
		}
	}

	return true
}
