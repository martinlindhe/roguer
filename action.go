package rogue

import "fmt"

// Action ...
type Action interface {
	// returns true when finished performing action
	Perform(npc *npc) bool
}

type sleep struct {
}

func (a *sleep) Perform(n *npc) bool {
	fmt.Println(n.Name, "is sleeping", n.Tiredness)
	n.Tiredness -= 10
	if n.Tiredness < 0 {
		n.Tiredness = 0
		return true
	}
	return false
}

type lookForFood struct {
	timeSpentLooking int
}

func (a *lookForFood) Perform(n *npc) bool {

	// XXX auto eat some food in inventory

	fmt.Println(n.Name, "is looking for food", a.timeSpentLooking)

	// TODO something more advanced for looking for food
	a.timeSpentLooking++
	if a.timeSpentLooking > 5 {
		n.Inventory = append(n.Inventory, &sweetPotato{})

		// XXX reduce hunger by some amount from the food eaten
		n.Hunger = 0
		return true
	}

	return false
}

type lookForWater struct {
	timeSpentLooking int
}

func (a *lookForWater) Perform(n *npc) bool {
	fmt.Println(n.Name, "is looking for water", a.timeSpentLooking)
	// TODO something more advanced for looking for food
	a.timeSpentLooking++
	if a.timeSpentLooking > 5 {
		// XXX reduce thirst by some amount from the water drunk
		n.Thirst = 0
		return true
	}

	return false
}
