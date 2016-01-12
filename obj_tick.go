package rogue

import (
	"fmt"
	"math/rand"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

// Announce to nearby players that something happened
func (n *Obj) Announce(format string, a ...interface{}) {

	str := fmt.Sprintf(format, a...)

	for _, pl := range island.Players {
		if pl.Spawn.Position.isNearby(n.Position) {
			log.Printf("XXXX tell %s: %s", pl.Name, str)
			b := []byte(str)
			pl.Socket.WriteMessage(websocket.TextMessage, b)
		}
	}
}

// Tick until it returns false
func (n *Obj) Tick() bool {
	n.Age++

	// log.Println("[tick]", n.Name, n.Age)

	if n.isAboveMaxAge() {
		n.Announce("%s dies of old age", n.Name)
		return false
	}

	n.treeTick()

	if n.Type == "fireplace" && n.Activated {
		n.Announce("%s is burning (%d energy left)", n.Name, n.Energy)
		n.Energy--
		if n.Energy <= 0 {
			n.Energy = 0
			n.Activated = false
			n.Announce("%s burned out", n.Name)
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
			n.Announce("%s drops a %s", n, drop.Name)

			spawnPos, err := n.Position.randomNearby()
			if err == nil {
				if n.Position != spawnPos {
					n.Announce("%s lands at %s, from %s", drop.Name, spawnPos, n)
				}

				island.addNpcFromName(drop.Name, spawnPos)
			} else {
				log.Errorf("XXX failed to find pos nearby %s", n)
			}
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

		nearbyFireplaces := n.Position.spawnsByType("fireplace", 1)
		if len(nearbyFireplaces) > 0 {

			fireplace := nearbyFireplaces[0]
			if fireplace.isActivated() {
				prevColdness := n.Coldness
				n.Coldness -= 100
				if n.Coldness < 0 {
					n.Coldness = 0
				}
				n.Announce("%s is getting warmed up by the %s (coldness -%d)", n, fireplace, prevColdness-n.Coldness)
			} else {

				// NOTE: some max capacity for the fireplace is required
				if fireplace.Energy < 1000 {
					itemIdx, err := n.tryFindItemTypeInInventory("wood")
					if err == nil {
						item := n.removeFromInventory(itemIdx)

						n.Announce("%s is putting %s in the %s", n, item.Name, fireplace)
						// NOTE: to simplify, we just get the energy from the wood directly
						fireplace.Energy += item.Energy
					}
				}

				if fireplace.Energy > 0 {
					n.Announce("%s lights the %s", n, fireplace)
					fireplace.Activate()

					// stay here for a bit
					n.planAction("wait")
				}
			}
		}

		if !n.hasPlannedType("travel") && n.hasItemTypeInInventory("wood") {
			fireplaces := n.Position.spawnsByType("fireplace", 30)

			if len(fireplaces) > 0 {
				if n.distanceTo(fireplaces[0].Position) > 1 {
					n.Announce("%s is freezing, moving to nearest fireplace at %v", n.Name, fireplaces[0].Position)
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

	n.survivalPlanningTick()

	// select one action to be doing next
	if n.CurrentAction == nil && len(n.PlannedActions) > 0 {
		// shuffle actions
		if len(n.PlannedActions) > 1 {
			shuffleActionSlice(n.PlannedActions)
		}

		// pick first
		n.CurrentAction = &n.PlannedActions[0]
		n.PlannedActions = n.PlannedActions[1:]

		n.Announce("%s started to %s", n.Name, n.CurrentAction.Name)
	}

	n.performCurrentAction()
	return true
}

func (n *Obj) survivalPlanningTick() {

	if !n.isTired() && !n.isHungry() && !n.isThirsty() && !n.isCold() && !n.hasPlannedType("travel") {
		// when basic needs is resolved, randomly decide to do
		// something that would help improve situation for the npc
		if n.Race == "rabbit" {
			if len(n.Position.spawnsByType("small hole", 30)) == 0 {
				n.planAction("dig small hole", n.Position)
				return
			}
		}

		if n.Type == "humanoid" {

			if island.canBuildAt(n.Position) && !n.hasPlannedType("build") {
				if len(n.Position.spawnsByType("fireplace", 30)) == 0 {
					// XXX if more than 1 humanoid nearby, instead build a larger fireplace
					n.planAction("build small fireplace", n.Position)
					return
				}
				if n.Home == nil && len(n.Position.spawnsByType("shelter", 30)) == 0 {
					// XXX if more than 1 humanoid nearby, instead build a small hut
					n.planAction("build small shelter", n.Position)
					return
				}

				if len(n.Position.spawnsByType("fireplace", 30)) > 0 &&
					len(n.Position.spawnsByType("shelter", 30)) > 0 {

					// basic survival is satisifed, lets build a cooking pit
					if len(n.Position.spawnsByType("cooking", 30)) == 0 {
						n.planAction("build cooking pit", n.Position)
						return
					}

					// build a hut if we already have a small shelter
					if n.Home != nil && n.Home.Name == "small shelter" && len(n.Position.spawnsByName("small hut", 30)) == 0 {
						n.planAction("build small hut", n.Position)
						return
					}
				}

				if len(n.Position.spawnsByName("farmland", 1)) == 0 {
					n.planAction("build farmland", n.Position)
					return
				}

				if len(n.Position.spawnsByName("apple tree", 30)) == 0 {
					// XXX require having a apple seed
					// XXX require having a garden, plant there
					n.planAction("plant apple tree", n.Position)
					return
				}
			}
		}
	}
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
			n.Announce("%s is feeling tired, decided to sleep (%d tiredness, cap = %d)", n.Name, n.Tiredness, n.tirednessCap())
			n.planAction("sleep")
			return true
		}

		nearbyShelters := n.Position.spawnsByType(shelterType, 0)
		if len(nearbyShelters) > 0 {
			n.Announce("%s is feeling tired, decided to sleep at %s (%d tiredness, cap = %d)", n.Name, nearbyShelters[0].Name, n.Tiredness, n.tirednessCap())
			n.planAction("sleep")
			return true
		}

		shelters := n.Position.spawnsByType(shelterType, 30)
		if len(shelters) == 0 {
			n.Announce("%s is feeling tired, decided to sleep (%d tiredness, cap = %d)", n.Name, n.Tiredness, n.tirednessCap())
			n.planAction("sleep")
			return true
		}

		n.Announce("%s is feeling tired, decided to go to %s for sleeping", n.Name, shelters[0].Name)
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
			n.Announce("%s ate %s (-%d hunger)", n.Name, item.Name, energyDiff)
			return true
		}

		if n.isHungry() && !n.hasPlanned("find food") {
			n.Announce("%s is feeling hungry (%d hunger)", n.Name, n.Hunger)
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
			n.Announce("%s drank %s (-%d thirst)", n.Name, item.Name, energyDiff)
			return true
		}
		if n.isThirsty() && !n.hasPlanned("find water") {
			n.Announce("%s is feeling thirsty (%d thirst)", n.Name, n.Thirst)
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
