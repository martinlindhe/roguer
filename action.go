package rogue

import log "github.com/Sirupsen/logrus"

// Action ...
type Action interface {
	// returns true when finished performing action
	Perform(npc *Npc) bool
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

type sleep struct {
	timeSpent int
}

func (a *sleep) Perform(n *Npc) bool {
	energyGain := 4

	log.Printf("%s is sleeping. tiredness = %d", n.Name, n.Tiredness)
	a.timeSpent++
	n.Tiredness -= energyGain

	if n.Tiredness <= 0 {
		//log.Printf("%s woke up. tiredness = %d", n.Name, n.Tiredness)
		n.Tiredness = 0
		return true
	}

	if a.timeSpent > 60 {
		// never sleep more than 60 ticks
		return true
	}

	return false
}

type lookForFood struct {
	timeSpent int
}

func (a *lookForFood) Perform(n *Npc) bool {

	log.Println(n.Name, "is looking for food", a.timeSpent)

	// TODO something more advanced for looking for food
	a.timeSpent++
	if a.timeSpent > 5 {

		item := island.randomItemOfType("food")
		log.Printf("%s found a %s", n.Name, item.Name)
		n.Inventory = append(n.Inventory, item)
		return true
	}

	return false
}

type lookForWater struct {
	timeSpent int
}

func (a *lookForWater) Perform(n *Npc) bool {
	log.Println(n.Name, "is looking for water", a.timeSpent)

	// TODO something more advanced for looking for water
	a.timeSpent++
	if a.timeSpent > 5 {

		item := island.randomItemOfType("drink")
		log.Printf("%s found a %s", n.Name, item.Name)
		n.Inventory = append(n.Inventory, item)

		return true
	}

	return false
}
