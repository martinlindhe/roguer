package rogue

import (
	"fmt"
	"math/rand"

	log "github.com/Sirupsen/logrus"
)

// Tick ...
func (n *Npc) Tick() {
	n.Age++

	n.Hunger++
	n.Tiredness++
	n.Thirst++

	log.Debug("[tick]", n.Name, n.Age)

	if n.isSleeping() {
		if n.CurrentAction.Name != "sleep" {
			// XXX this should never happen
			panic(fmt.Errorf("sleeping and doing something that requires being awake: %s", n.CurrentAction.Name))
		}
		n.performCurrentAction()
		return
	}

	if n.isTired() && !n.hasPlanned("sleep") {
		log.Printf("%s is feeling tired (%d tiredness, cap = %d)", n.Name, n.Tiredness, n.tirednessCap())
		n.planAction("sleep")
	}

	if n.isHungry() {

		// auto eat some food in inventory instead of looking for food, if possible
		itemIdx, err := n.tryFindItemTypeInInventory("food")
		if err == nil {
			item := n.removeFromInventory(itemIdx)

			prevHunger := n.Hunger

			// eat item: reduce hunger by some amount from the food eaten
			n.Hunger -= item.Energy
			if n.Hunger < 0 {
				n.Hunger = 0
			}

			energyDiff := prevHunger - n.Hunger
			log.Printf("%s ate %s (-%d hunger)", n.Name, item.Name, energyDiff)
			return
		}

		if n.isHungry() && !n.hasPlanned("find food") {
			log.Printf("%s is feeling hungry (%d hunger)", n.Name, n.Hunger)
			n.planAction("find food")
		}
	}

	if n.isThirsty() {

		// auto eat some food in inventory instead of looking for food, if possible
		itemIdx, err := n.tryFindItemTypeInInventory("drink")
		if err == nil {
			item := n.removeFromInventory(itemIdx)

			prevThirst := n.Thirst

			// eat item: reduce hunger by some amount from the food eaten
			n.Thirst -= item.Energy
			if n.Thirst < 0 {
				n.Thirst = 0
			}

			energyDiff := prevThirst - n.Thirst
			log.Printf("%s drank %s (-%d thirst)", n.Name, item.Name, energyDiff)
			return
		}
		if n.isThirsty() && !n.hasPlanned("find water") {
			log.Printf("%s is feeling thirsty (%d thirst)", n.Name, n.Thirst)
			n.planAction("find water")
		}
	}

	if !n.isTired() && !n.isHungry() && !n.isThirsty() {
		// when basic needs is resolved, randomly decide to do
		// something that would help improve situation for the npc
		if n.Race == "rabbit" {
			if len(island.withinRadiusOfName("small hole", 30, n.Position)) == 0 {
				n.planAction("dig small hole")
			}
		}

		if n.Type == "humanoid" {

			if !n.hasItemInInventory("firewood") {
				n.planAction("find firewood")
			}

			//if island.canBuildAt(n.Position) {
			if len(island.withinRadiusOfType("fireplace", 30, n.Position)) == 0 {
				n.planAction("build small fireplace")
				// XXX if more than 1 humanoid nearby, instead build a larger fireplace
			}
			if len(island.withinRadiusOfType("shelter", 30, n.Position)) == 0 {
				// XXX if more than 1 humanoid nearby, instead build a small hut
				n.planAction("build small shelter")
			}
			if len(island.withinRadiusOfName("farmland", 1, n.Position)) == 0 {
				n.planAction("build farmland")
			}
			if len(island.withinRadiusOfName("apple tree", 30, n.Position)) == 0 {
				n.planAction("plant apple tree")
			}
			//}
		}

		if n.Type == "tree" {
			// XXX tick for tree
		}
	}

	// select one action to be doing next
	if n.CurrentAction == nil && len(n.PlannedActions) > 0 {
		// shuffle actions
		if len(n.PlannedActions) > 1 {
			shuffleActionSlice(n.PlannedActions)
		}

		// pick first
		n.CurrentAction = &n.PlannedActions[0]
		n.PlannedActions = n.PlannedActions[1:]

		log.Println(n.Name, "started to", n.CurrentAction.Name)
	}

	n.performCurrentAction()
}

// shuffle slice, without allocations
func shuffleActionSlice(p []actionSpec) {

	for i := range p {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}
