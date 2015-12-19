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
	PlannedActions []Action

	// the lower value, the less hungry npc is
	Hunger    int
	Thirst    int
	Tiredness int
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

func (n *Npc) thirstCap() int {
	// XXX
	return n.Level * 10
}

func (n *Npc) tirednessCap() int {
	// XXX
	return n.Level * 10
}

// check if npc already has planned to do a
func (n *Npc) hasPlanned(a Action) bool {
	for _, v := range n.PlannedActions {
		if v == a {
			return true
		}
	}
	return false
}

func (n *Npc) Tick() {
	n.Age++

	n.Hunger++
	n.Tiredness++
	n.Thirst++

	fmt.Println("[tick]", n.Name, n.Age)

	if n.Tiredness > n.tirednessCap() && !n.hasPlanned(Sleep{}) {
		fmt.Println(n.Name, "is feeling tired")
		n.PlannedActions = append(n.PlannedActions, Sleep{})
	}

	if n.Hunger > n.hungerCap() && !n.hasPlanned(LookForFood{}) {
		fmt.Println(n.Name, "is feeling hungry")
		n.PlannedActions = append(n.PlannedActions, LookForFood{})
	}
	if n.Thirst > n.thirstCap() && !n.hasPlanned(LookForWater{}) {
		fmt.Println(n.Name, "is feeling thirsty")
		n.PlannedActions = append(n.PlannedActions, LookForWater{})
	}

	if n.CurrentlyDoing == DoingNothing {
		// XXX shuffle action list, so behaviour is more random
		// XXX if action in queue, make it "currently doing"
	}

	fmt.Println(n.PlannedActions)
}
