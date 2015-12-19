package rogue

import "fmt"

type Point struct {
	X uint16
	Y uint16
}

type WorldObject interface {
	Tick()
	Defaults()
}

const (
	DoingNothing = 0
	DoingSleeping
	DoingEating
	DoingDrinking
	DoingMoving
	DoingForaging
)

type Npc struct {
	Level          int
	XP             int
	Age            int
	Name           string
	Position       Point
	CurrentlyDoing int

	// the lower value, the less hungry npc is
	Hunger int
	Thirst int

	// the higher value, the faster npc get hungry
	HungerMod int
}

type Dwarf struct {
	Npc
}

type Rabbit struct {
	Npc
}

func (n *Npc) Defaults() {
	// init non-zero values
	n.Level = 1
}

func (n *Npc) hungerCap() int {
	// XXX
	return n.Level * 10
}

func (n *Npc) Tick() {
	n.Age++

	n.Hunger++

	fmt.Println("[tick]", n.Name, n.Age)

	if n.Hunger > n.hungerCap() {
		fmt.Println("HUNGRY!", n.Hunger, n.hungerCap())

		// XXX enqueue action
	}

	if n.CurrentlyDoing == DoingNothing {
		// XXX shuffle action list, so behaviour is more random
		// XXX if action in queue, make it "currently doing"
	}
}
