package rogue

import (
	"fmt"
	"math/rand"
	"reflect"

	log "github.com/Sirupsen/logrus"
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
	//n.npc.Defaults()  // NOTE: plant is currently base class
	n.Name = "sdsdfgsdfg"
	log.Printf("plant defaults")
}

func (n *sweetPotato) Defaults() {
	n.plant.Defaults()
	n.Name = "sweet potato"
	fmt.Println("potato def: ", n.Name)
}

func (n *rabbit) Defaults() {
	n.Npc.Defaults()
	n.Name = "a rabbit"
}

func (n *dwarf) Defaults() {
	n.Npc.Defaults()
	n.Name = n.generateName()
	//log.Printf("dwarf defaults")
}

func (n *plant) Tick() {
	n.Age++
}

func (n *Npc) hungerCap() int {
	// XXX
	return n.Level * 5
}

func (n *Npc) thirstCap() int {
	// XXX
	return n.Level * 100
}

func (n *Npc) tirednessCap() int {
	// XXX
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

// shuffle slice, without allocations
func shuffleActionSlice(p []Action) {

	for i := range p {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}
