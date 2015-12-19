package rogue

import "fmt"

// Action ...
type Action interface {
	// returns true when finished performing action
	Perform(npc *npc) bool
}

type sleep struct {
}

func (a sleep) Perform(n *npc) bool {
	fmt.Println(n.Name, "is sleeping XXX")
	n.Tiredness -= 10
	if n.Tiredness < 0 {
		n.Tiredness = 0
		return true
	}
	return false
}

type lookForFood struct {
}

func (a lookForFood) Perform(n *npc) bool {
	fmt.Println(n.Name, "is looking for food XXX")
	return false
}

type lookForWater struct {
}

func (a lookForWater) Perform(n *npc) bool {
	fmt.Println(n.Name, "is looking for water XXX")
	return false
}
