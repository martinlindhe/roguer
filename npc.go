package rogue

import (
	"math/rand"
	"reflect"
)

// Point ...
type Point struct {
	X uint16
	Y uint16
}

// WorldObject ...
type WorldObject interface {
	Tick()
	Defaults()
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

// Defaults ...
func (n *WorldObjectInstance) Defaults() {
	// init non-zero values
	n.Level = 1
	//log.Debug("npc defaults")
}

func (n *plant) Defaults() {
	n.WorldObjectInstance.Defaults()
}

func (n *sweetPotato) Defaults() {
	n.plant.Defaults()
	n.Name = "sweet potato"
	n.Energy = 5
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

// shuffle slice, without allocations
func shuffleActionSlice(p []Action) {

	for i := range p {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}
