package rogue

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
)

// Tick base class
func (n *WorldObjectInstance) Tick() {
	n.Age++
}

// Tick ...
func (n *Npc) Tick() {
	n.WorldObjectInstance.Tick()

	n.Hunger++
	n.Tiredness++
	n.Thirst++

	//fmt.Println("[tick]", n.Name, n.Age)

	if n.Tiredness > n.tirednessCap() && !n.hasPlanned(&sleep{}) {
		log.Printf("%s is feeling tired. tiredness = %d, cap = %d", n.Name, n.Tiredness, n.tirednessCap())
		n.PlannedActions = append(n.PlannedActions, &sleep{})
	}

	if n.Hunger > n.hungerCap() {

		// auto eat some food in inventory instead of looking for food, if possible
		item := n.pickSomethingToEat()
		if item != nil {

			// XXX EAT ETC
			n.Hunger = 0

			/*
					// XXX unpack interface
					food := item.(WorldObject)

					log.Printf("%s picked something to eat from inventory: %s with energy %d", n.Name, food.Name, food.Energy)
					// reduce hunger by some amount from the food eaten

				n.Hunger -= food.Energy
				if n.Hunger < 0 {
					n.Hunger = 0
				}
			*/
		}

		if n.Hunger > n.hungerCap() && !n.hasPlanned(&lookForFood{}) {
			log.Println(n.Name, "is feeling hungry")
			n.PlannedActions = append(n.PlannedActions, &lookForFood{})
		}
	}

	if n.Thirst > n.thirstCap() && !n.hasPlanned(&lookForWater{}) {
		log.Println(n.Name, "is feeling thirsty")
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

		log.Println(n.Name, "decided to", reflect.TypeOf(n.CurrentAction))
	}

	if n.CurrentAction != nil {
		if n.CurrentAction.Perform(n) == true {
			log.Println(n.Name, "finished performing", reflect.TypeOf(n.CurrentAction))
			n.CurrentAction = nil
		}
	}
}
