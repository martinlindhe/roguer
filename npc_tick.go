package rogue

import (
	"log"
	"reflect"
)

// Tick ...
func (n *WorldObjectInstance) Tick() {
	n.Age++
}

// Tick ...
func (n *Npc) Tick() {
	n.WorldObjectInstance.Tick() // call base

	n.Hunger++
	n.Tiredness++
	n.Thirst++

	//fmt.Println("[tick]", n.Name, n.Age)

	if n.Tiredness > n.tirednessCap() && !n.hasPlanned(&sleep{}) {
		log.Printf("%s is feeling tired. tiredness = %d, cap = %d", n.Name, n.Tiredness, n.tirednessCap())
		n.PlannedActions = append(n.PlannedActions, &sleep{})
	}

	if n.Hunger > n.hungerCap() && !n.hasPlanned(&lookForFood{}) {
		log.Println(n.Name, "is feeling hungry")
		n.PlannedActions = append(n.PlannedActions, &lookForFood{})
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
