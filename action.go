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
}

func (a *sleep) Perform(n *Npc) bool {
	energyGain := 4

	log.Printf("%s is sleeping. tiredness = %d", n.Name, n.Tiredness)
	n.Tiredness -= energyGain
	if n.Tiredness < 0 {
		//log.Printf("%s woke up. tiredness = %d", n.Name, n.Tiredness)
		n.Tiredness = 0
		return true
	}
	return false
}

type lookForFood struct {
	timeSpentLooking int
}

func (a *lookForFood) Perform(n *Npc) bool {

	// XXX auto eat some food in inventory instead of looking for food, if possible

	log.Println(n.Name, "is looking for food", a.timeSpentLooking)

	// TODO something more advanced for looking for food
	a.timeSpentLooking++
	if a.timeSpentLooking > 5 {

		food := getRandomFoodFrom(&n.Position)

		log.Printf("%s found something to eat: %s", n.Name, food.Name)
		n.Inventory = append(n.Inventory, &food)

		// XXX reduce hunger by some amount from the food eaten
		n.Hunger = 0
		return true
	}

	return false
}

func getRandomFoodFrom(p *Point) WorldObjectInstance {

	// XXX?!?!
	var food sweetPotato
	food.Defaults()

	return food
}

type lookForWater struct {
	timeSpentLooking int
}

func (a *lookForWater) Perform(n *Npc) bool {
	log.Println(n.Name, "is looking for water", a.timeSpentLooking)
	// TODO something more advanced for looking for food
	a.timeSpentLooking++
	if a.timeSpentLooking > 5 {
		// XXX reduce thirst by some amount from the water drunk
		n.Thirst = 0
		return true
	}

	return false
}
