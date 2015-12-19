package rogue

import (
	"fmt"
	"math/rand"
	"reflect"
)

// Point ...
type Point struct {
	X uint16
	Y uint16
}

type worldObject interface {
	Tick()
	Defaults()
}

// Doing states
const (
	doingNothing = 0
	doingSleeping
	doingEating
	doingDrinking
	doingMoving
	doingForaging
)

type npc struct {
	Level          int
	XP             int
	Age            int
	Name           string
	Position       Point
	CurrentAction  Action
	PlannedActions []Action
	Inventory      []worldObject

	// the lower value, the less hungry npc is
	Hunger    int
	Thirst    int
	Tiredness int
}

type plant struct {
	Name     string
	Position Point
	Age      int
}

type sweetPotato struct {
	plant
}

type dwarf struct {
	npc
}

type rabbit struct {
	npc
}

func (n *npc) Defaults() {
	// init non-zero values
	n.Level = 1
}

func (n *plant) Defaults() {
	// init non-zero values
	n.Name = "sdsdfgsdfg"
}

func (n *plant) Tick() {
	n.Age++
}

func (n *npc) hungerCap() int {
	// XXX
	return n.Level * 100
}

func (n *npc) thirstCap() int {
	// XXX
	return n.Level * 100
}

func (n *npc) tirednessCap() int {
	// XXX
	return n.Level * 100
}

// check if npc already has planned to do a
func (n *npc) hasPlanned(a Action) bool {

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

func (n *npc) Tick() {

	n.Age++
	n.Hunger++
	n.Tiredness++
	n.Thirst++

	//fmt.Println("[tick]", n.Name, n.Age)

	if n.Tiredness > n.tirednessCap() && !n.hasPlanned(&sleep{}) {
		fmt.Println(n.Name, "is feeling tired")
		n.PlannedActions = append(n.PlannedActions, &sleep{})
	}

	if n.Hunger > n.hungerCap() && !n.hasPlanned(&lookForFood{}) {
		fmt.Println(n.Name, "is feeling hungry")
		n.PlannedActions = append(n.PlannedActions, &lookForFood{})
	}
	if n.Thirst > n.thirstCap() && !n.hasPlanned(&lookForWater{}) {
		fmt.Println(n.Name, "is feeling thirsty")
		n.PlannedActions = append(n.PlannedActions, &lookForWater{})
	}

	// select one action to be doing next
	if n.CurrentAction == nil && len(n.PlannedActions) > 0 {
		// shuffle action list, so behaviour is more random
		if len(n.PlannedActions) > 1 {
			shuffleActionSlice(n.PlannedActions)
		}

		// pick something
		n.CurrentAction = n.PlannedActions[0]
		n.PlannedActions = n.PlannedActions[1:]

		fmt.Println(n.Name, "decided to", reflect.TypeOf(n.CurrentAction))
	}

	if n.CurrentAction != nil {
		if n.CurrentAction.Perform(n) == true {
			fmt.Println(n.Name, "finished performing", reflect.TypeOf(n.CurrentAction))
			n.CurrentAction = nil
		}
	}
}

// shuffle slice, without allocations
func shuffleActionSlice(p []Action) {

	for i := range p {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}
