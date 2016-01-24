package rogue

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

// Announce to nearby players that something happened
func (o *Obj) Announce(format string, a ...interface{}) {

	str := fmt.Sprintf(format, a...)

	generalLog.Debug(o.Name, " announces: ", str)

	for _, pl := range o.Island.Players {
		if pl.Spawn.Position.isNearby(o.Position) {
			//log.Printf("tell %s: %s", pl.Name, str)

			res := messageResponse{Type: "msg", Message: str, Time: o.Island.Age.Current()}

			b, _ := json.Marshal(res)
			pl.Socket.WriteMessage(websocket.TextMessage, b)
		}
	}
}

// Tick until it returns false
func (o *Obj) Tick() bool {
	o.Age.Tick()

	//generalLog.Debug("[tick]", n.Name, n.Age)
	//log.Info("[tick]", n.Name, n.Age)

	if o.isAboveMaxAge() {
		o.Announce("%s dies of old age", o.Name)
		return false
	}

	if o.Hunger > o.diesOfHungerCap() {
		o.Announce("%s dies of hunger", o.Name)
		return false
	}

	if o.Thirst > o.diesOfThirstCap() {
		o.Announce("%s dies of thirst", o.Name)
		return false
	}

	o.treeTick()

	if o.Type == "fireplace" && o.Activated {
		o.Announce("%s is burning (%d energy left)", o.Name, o.Energy)
		o.Energy--
		if o.Energy <= 0 {
			o.Energy = 0
			o.Activated = false
			o.Announce("%s burned out", o.Name)
		}
	}

	return o.npcTick()
}

func (o *Obj) treeTick() {

	if o.Type != "tree" {
		return
	}

	treeSpec := o.Island.getNpcSpecFromName(o.Name)

	for _, drop := range treeSpec.Drops {

		roll := float64(rand.Intn(100)) // between 0-99
		//log.Debugf("Rolled %f for check if %s is spawned, %f chance", roll, drop.Name, drop.Chance)

		if roll <= drop.Chance {
			o.Announce("%s falls from %s", drop.Name, o.Name)

			spawnPos, err := o.randomNearby()
			if err == nil {
				o.Island.addNpcFromName(drop.Name, spawnPos)
			} else {
				generalLog.Error("Failed to find pos nearby", o)
			}
		}
	}
}

func (o *Obj) npcTick() bool {

	if o.Class != "npc" {
		return true
	}

	if o.isSleeping() {
		if o.CurrentAction.Name != "sleep" {
			panic(fmt.Errorf("sleeping and doing something that requires being awake: %s", o.CurrentAction.Name))
		}
		o.performCurrentAction()
		return true
	}

	o.Tiredness++

	if o.isCold() && !o.hasPlannedType("travel") && o.Type == "humanoid" {

		if !o.hasItemTypeInInventory("wood") && !o.hasPlannedType("wait") {
			o.planAction("find fire wood")
		}

		nearbyFireplaces := o.spawnsByType("fireplace", 1)
		if len(nearbyFireplaces) > 0 {

			fireplace := nearbyFireplaces[0]
			if fireplace.isActivated() {
				prevColdness := o.Coldness
				o.Coldness -= 100
				if o.Coldness < 0 {
					o.Coldness = 0
				}
				o.Announce("%s is getting warmed up by the %s (coldness -%d)", o.Name, fireplace, prevColdness-o.Coldness)
			} else {

				// NOTE: some max capacity for the fireplace is required
				if fireplace.Energy < 1000 {
					itemIdx, err := o.tryFindItemTypeInInventory("wood")
					if err == nil {
						item := o.removeFromInventory(itemIdx)

						o.Announce("%s is putting %s in the %s", o.Name, item.Name, fireplace)
						// NOTE: to simplify, we just get the energy from the wood directly
						fireplace.Energy += item.Energy
					}
				}

				if fireplace.Energy > 0 {
					o.Announce("%s lights the %s", o.Name, fireplace)
					fireplace.Activate()

					// stay here for a bit
					o.planAction("wait")
				}
			}
		}

		if !o.hasPlannedType("travel") && o.hasItemTypeInInventory("wood") {
			fireplaces := o.spawnsByType("fireplace", 30)

			if len(fireplaces) > 0 {
				if o.distanceTo(&fireplaces[0].Position) > 1 {
					o.Announce("%s is freezing, moving to nearest fireplace at %v", o.Name, fireplaces[0].Position)
					o.planAction("walk", fireplaces[0].Position)
				}
			}
		}
	}

	if o.hungerThirstTick() {
		return true
	}

	if o.tiredTick() {
		return true
	}

	o.survivalPlanningTick()

	// select one action to be doing next
	if o.CurrentAction == nil && len(o.PlannedActions) > 0 {
		// shuffle actions
		if len(o.PlannedActions) > 1 {
			shuffleActionSlice(o.PlannedActions)
		}

		// pick first
		o.CurrentAction = &o.PlannedActions[0]
		o.PlannedActions = o.PlannedActions[1:]

		o.Announce("%s started to %s", o.Name, o.CurrentAction.Name)
	}

	o.performCurrentAction()
	return true
}

func (o *Obj) survivalPlanningTick() {

	if !o.isTired() && !o.isHungry() && !o.isThirsty() && !o.isCold() && !o.hasPlannedType("travel") {
		// when basic needs is resolved, randomly decide to do
		// something that would help improve situation for the npc
		if o.Race == "rabbit" {
			if len(o.spawnsByType("small hole", 30)) == 0 {
				o.planAction("dig small hole", o.Position)
				return
			}
		}

		if o.Type == "humanoid" {

			if o.Island.canBuildAt(o.Position) && !o.hasPlannedType("build") {
				if len(o.spawnsByType("fireplace", 30)) == 0 {
					// XXX if more than 1 humanoid nearby, instead build a larger fireplace
					o.planAction("build small fireplace", o.Position)
					return
				}
				if o.Home == nil && len(o.spawnsByType("shelter", 30)) == 0 {
					// XXX if more than 1 humanoid nearby, instead build a small hut
					o.planAction("build small shelter", o.Position)
					return
				}

				if len(o.spawnsByType("fireplace", 30)) > 0 &&
					len(o.spawnsByType("shelter", 30)) > 0 {

					// basic survival is satisifed, lets build a cooking pit
					if len(o.spawnsByType("cooking", 30)) == 0 {
						o.planAction("build cooking pit", o.Position)
						return
					}

					// build a hut if we already have a small shelter
					if o.Home != nil && o.Home.Name == "small shelter" && len(o.spawnsByName("small hut", 30)) == 0 {
						o.planAction("build small hut", o.Position)
						return
					}
				}

				if len(o.spawnsByName("farmland", 1)) == 0 {
					o.planAction("build farmland", o.Position)
					return
				}

				if len(o.spawnsByName("apple tree", 30)) == 0 {
					// XXX require having a apple seed
					// XXX require having a garden, plant there
					o.planAction("plant apple tree", o.Position)
					return
				}
			}
		}
	}
}

func (o *Obj) preferredShelterType() string {

	shelterType := ""
	if o.Type == "humanoid" {
		shelterType = "shelter"
	} else if o.Type == "rodent" {
		shelterType = "burrow"
	}
	return shelterType
}

func (o *Obj) tiredTick() bool {

	shelterType := o.preferredShelterType()

	// if next to shelter, sleep. if shelter nearby, go there and then sleep
	if o.isTired() && !o.hasPlannedType("sleep") && !o.hasPlannedType("travel") {

		if shelterType == "" {
			o.Announce("%s is feeling tired, decided to sleep (%d tiredness, cap = %d)", o.Name, o.Tiredness, o.tirednessCap())
			o.planAction("sleep")
			return true
		}

		nearbyShelters := o.spawnsByType(shelterType, 0)
		if len(nearbyShelters) > 0 {
			o.Announce("%s is feeling tired, decided to sleep at %s (%d tiredness, cap = %d)", o.Name, nearbyShelters[0].Name, o.Tiredness, o.tirednessCap())
			o.planAction("sleep")
			return true
		}

		shelters := o.spawnsByType(shelterType, 30)
		if len(shelters) == 0 {
			o.Announce("%s is feeling tired, decided to sleep (%d tiredness, cap = %d)", o.Name, o.Tiredness, o.tirednessCap())
			o.planAction("sleep")
			return true
		}

		o.Announce("%s is feeling tired, decided to go to %s for sleeping", o.Name, shelters[0].Name)
		o.planAction("walk", shelters[0].Position)
	}

	return false
}

func (o *Obj) hungerThirstTick() bool {

	o.Hunger++
	o.Thirst++

	if o.isHungry() {

		// auto eat some food in inventory instead of looking for food, if possible
		itemIdx, err := o.tryFindItemTypeInInventory("food")
		if err == nil {
			item := o.removeFromInventory(itemIdx)

			prevHunger := o.Hunger

			// eat item: reduce hunger by some amount from the food eaten
			o.Hunger -= item.Energy
			if o.Hunger < 0 {
				o.Hunger = 0
			}

			energyDiff := prevHunger - o.Hunger
			o.Announce("%s ate %s (-%d hunger)", o.Name, item.Name, energyDiff)
			return true
		}

		if o.isHungry() && !o.hasPlanned("find food") {
			o.Announce("%s is feeling hungry (%d hunger)", o.Name, o.Hunger)
			o.planAction("find food")
		}
	}

	if o.isThirsty() {

		// auto eat some food in inventory instead of looking for food, if possible
		itemIdx, err := o.tryFindItemTypeInInventory("drink")
		if err == nil {
			item := o.removeFromInventory(itemIdx)

			prevThirst := o.Thirst

			// eat item: reduce hunger by some amount from the food eaten
			o.Thirst -= item.Energy
			if o.Thirst < 0 {
				o.Thirst = 0
			}

			energyDiff := prevThirst - o.Thirst
			o.Announce("%s drank %s (-%d thirst)", o.Name, item.Name, energyDiff)
			return true
		}
		if o.isThirsty() && !o.hasPlanned("find water") {
			o.Announce("%s is feeling thirsty (%d thirst)", o.Name, o.Thirst)
			o.planAction("find water")
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
