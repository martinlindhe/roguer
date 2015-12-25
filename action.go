package rogue

import log "github.com/Sirupsen/logrus"

// Action ...
type Action interface {
	// returns true when finished performing action
	Perform(npc *Npc) bool
}

type actionList struct {
	All []actionYaml `json:"all"`
}

type actionYaml struct {
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	TimeSpent int
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

func (n *Npc) performSleep() bool {
	energyGain := 4

	log.Printf("%s is sleeping. tiredness = %d", n.Name, n.Tiredness)
	n.TimeSpentOnCurrentAction++
	n.Tiredness -= energyGain

	if n.Tiredness <= 0 {
		//log.Printf("%s woke up. tiredness = %d", n.Name, n.Tiredness)
		n.Tiredness = 0
		return true
	}

	if n.TimeSpentOnCurrentAction > 30 {
		// never sleep more than 30 ticks
		return true
	}

	return false
}

func (n *Npc) performFindFood() bool {

	log.Println(n.Name, "is looking for food")

	// TODO something more advanced for looking for food
	n.TimeSpentOnCurrentAction++
	if n.TimeSpentOnCurrentAction > 5 {

		item := island.randomItemOfType("food")
		log.Printf("%s found a %s", n.Name, item.Name)
		n.Inventory = append(n.Inventory, item)
		return true
	}

	return false
}

func (n *Npc) performFindWater() bool {
	log.Println(n.Name, "is looking for water")

	// TODO something more advanced for looking for water
	n.TimeSpentOnCurrentAction++
	if n.TimeSpentOnCurrentAction > 5 {

		item := island.randomItemOfType("drink")
		log.Printf("%s found a %s", n.Name, item.Name)
		n.Inventory = append(n.Inventory, item)

		return true
	}

	return false
}
