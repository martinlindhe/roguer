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

	n.treeTick()

	if n.Type == "fireplace" && n.Activated {
		log.Printf("%s is burning (%d energy left)", n.Name, n.Energy)
		n.Energy--
		if n.Energy <= 0 {
			n.Energy = 0
			n.Activated = false
			log.Printf("%s burned out", n.Name)
		}
	}

	return n.npcTick()
}

func (n *Obj) treeTick() {
	if n.Type != "tree" {
		return
	}

	treeSpec := island.getNpcSpecFromName(n.Name)

	for _, drop := range treeSpec.Drops {

		roll := float64(rand.Intn(100)) // between 0-99
		//log.Debugf("Rolled %f for check if %s is spawned, %f chance", roll, drop.Name, drop.Chance)

		if roll <= drop.Chance {
			log.Printf("%s drops a %s", n.Name, drop.Name)

			spawnPos := n.Position.randomNearby()
			if n.Position != spawnPos {
				log.Debugf("%s lands at %s, from %s", drop.Name, spawnPos, n)
			}

			island.addNpcFromName(drop.Name, spawnPos)
		}
	}
}

func (n *Obj) npcTick() bool {

	if n.Class != "npc" {
		return true
	}

	n.Hunger++
	n.Thirst++

	if n.isSleeping() {
		if n.CurrentAction.Name != "sleep" {
			panic(fmt.Errorf("sleeping and doing something that requires being awake: %s", n.CurrentAction.Name))
		}
		n.performCurrentAction()
		return true
	}

	n.Tiredness++

	if n.isCold() && !n.hasPlannedType("travel") && n.Type == "humanoid" {

		if !n.hasItemTypeInInventory("wood") && !n.hasPlannedType("wait") {
			n.planAction("find fire wood")
		}

		nearbyFireplaces := island.withinRadiusOfType("fireplace", 1, n.Position)
		if len(nearbyFireplaces) > 0 {

			fireplace := nearbyFireplaces[0]
			if fireplace.isActivated() {
				n.Coldness -= 100
				if n.Coldness < 0 {
					n.Coldness = 0
				}
				log.Printf("%s is getting warmed up by the %s (coldness %d)", n, fireplace, n.Coldness)
			} else {

				// NOTE: some max capacity for the fireplace is required
				if fireplace.Energy < 1000 {
					itemIdx, err := n.tryFindItemTypeInInventory("wood")
					if err == nil {
						item := n.removeFromInventory(itemIdx)

						log.Printf("%s is putting %s in the %s", n, item.Name, fireplace)
						// NOTE: to simplify, we just get the energy from the wood directly
						fireplace.Energy += item.Energy
					}
				}

				if fireplace.Energy > 0 {
					log.Printf("%s lights the fireplace", n.Name)
					fireplace.Activate()

					// stay here for a bit
					n.planAction("wait")
				}
			}
		}

		if !n.hasPlannedType("travel") {
			fireplaces := island.withinRadiusOfType("fireplace", 30, n.Position)

			if len(fireplaces) > 0 {
				if n.distanceTo(fireplaces[0].Position) > 1 {
					log.Printf("%s is freezing, moving to nearest fireplace at %v", n.Name, fireplaces[0].Position)
					n.planAction("walk", fireplaces[0].Position)
				}
			}
		}
	}

	if n.hungerThirstTick() {
		return true
	}

	if n.tiredTick() {
		return true
	}

	if !n.isTired() && !n.isHungry() && !n.isThirsty() && !n.hasPlannedType("travel") {
		// when basic needs is resolved, randomly decide to do
		// something that would help improve situation for the npc
		if n.Race == "rabbit" {
			if len(island.withinRadiusOfName("small hole", 30, n.Position)) == 0 {
				n.planAction("dig small hole")
			}
		}

		if n.Type == "humanoid" {

			if island.canBuildAt(n.Position) {
				if len(island.withinRadiusOfType("fireplace", 30, n.Position)) == 0 {
					// XXX if more than 1 humanoid nearby, instead build a larger fireplace
					n.planAction("build small fireplace")
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

func (n *Obj) preferredShelterType() string {
	shelterType := ""
	if n.Type == "humanoid" {
		shelterType = "shelter"
	} else if n.Type == "rodent" {
		shelterType = "burrow"
	}
	return shelterType
}

func (n *Obj) tiredTick() bool {

	shelterType := n.preferredShelterType()

	// if next to shelter, sleep. if shelter nearby, go there and then sleep
	if n.isTired() && !n.hasPlannedType("sleep") && !n.hasPlannedType("travel") {

		if shelterType == "" {
			log.Printf("%s is feeling tired, decided to sleep (%d tiredness, cap = %d)", n.Name, n.Tiredness, n.tirednessCap())
			n.planAction("sleep")
			return true
		}

		nearbyShelters := island.withinRadiusOfType(shelterType, 0, n.Position)
		if len(nearbyShelters) > 0 {
			log.Printf("%s is feeling tired, decided to sleep at %s (%d tiredness, cap = %d)", n.Name, nearbyShelters[0].Name, n.Tiredness, n.tirednessCap())
			n.planAction("sleep")
			return true
		}

		shelters := island.withinRadiusOfType(shelterType, 30, n.Position)
		if len(shelters) == 0 {
			log.Printf("%s is feeling tired, decided to sleep (%d tiredness, cap = %d)", n.Name, n.Tiredness, n.tirednessCap())
			n.planAction("sleep")
			return true
		}

		log.Printf("%s is feeling tired, decided to go to %s for sleeping", n.Name, shelters[0].Name)
		n.planAction("walk", shelters[0].Position)
	}

	return false
}

func (n *Obj) hungerThirstTick() bool {

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
			return true
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
			return true
		}
		if n.isThirsty() && !n.hasPlanned("find water") {
			log.Printf("%s is feeling thirsty (%d thirst)", n.Name, n.Thirst)
			n.planAction("find water")
		}
	}
	return false
}

// shuffle slice, without allocations
func shuffleActionSlice(p []actionSpec) {

	for i := range p {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}
