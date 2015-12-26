package rogue

import (
	"fmt"
	"math/rand"

	log "github.com/Sirupsen/logrus"
)

// Tick until it returns false
func (n *Obj) Tick() bool {
	n.Age++

	log.Debug("[tick]", n.Name, n.Age)

	if n.isAboveMaxAge() {
		log.Infof("%s dies of old age", n.Name)
		return false
	}

	if n.Type == "tree" {
		treeSpec := island.getNpcSpecFromName(n.Name)

		for _, drop := range treeSpec.Drops {

			roll := float64(rand.Intn(100)) // between 0-99
			//log.Debugf("Rolled %f for check if %s is spawned, %f chance", roll, drop.Name, drop.Chance)

			if roll <= drop.Chance {
				log.Printf("%s drops a %s", n.Name, drop.Name)
				island.addNpcFromName(drop.Name, n.Position)
			}
		}
	}

	return n.npcTick()
}

func (n *Obj) npcTick() bool {

	if n.Class != "npc" {
		return true
	}

	n.Hunger++
	n.Thirst++
	n.Tiredness++

	if n.isSleeping() {
		if n.CurrentAction.Name != "sleep" {
			panic(fmt.Errorf("sleeping and doing something that requires being awake: %s", n.CurrentAction.Name))
		}
		n.performCurrentAction()
		return true
	}

	fireplaces := island.withinRadiusOfType("fireplace", 30, n.Position)
	if n.isCold() && len(fireplaces) > 0 {
		// XXX move to most nearby fireplace

		fmt.Printf("XXXX move to pos %v", fireplaces[0].Position)
		// n.planAction("travel by foot", fireplaces[0].Position)
	}

	if n.isTired() && !n.hasPlanned("sleep") {
		log.Printf("%s is feeling tired (%d tiredness, cap = %d)", n.Name, n.Tiredness, n.tirednessCap())
		n.planAction("sleep")
	}

	n.hungerThirstTick()

	if !n.isTired() && !n.isHungry() && !n.isThirsty() {
		// when basic needs is resolved, randomly decide to do
		// something that would help improve situation for the npc
		if n.Race == "rabbit" {
			if len(island.withinRadiusOfName("small hole", 30, n.Position)) == 0 {
				n.planAction("dig small hole")
			}
		}

		if n.Type == "humanoid" {

			if !n.hasItemTypeInInventory("wood") {
				n.planAction("find fire wood")
			}

			if island.canBuildAt(n.Position) {
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
			}
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
	return true
}

func (n *Obj) hungerThirstTick() {

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

}

// shuffle slice, without allocations
func shuffleActionSlice(p []actionSpec) {

	for i := range p {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}
