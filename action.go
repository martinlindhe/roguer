package rogue

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

// Action ...
type Action interface {
	// returns true when finished performing action
	Perform(npc *Npc) bool
}

type actionList struct {
	All []actionSpec `json:"all"`
}

type actionSpec struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Duration  int    `json:"duration"`
	Energy    int    `json:"energy"`
	TimeSpent int
}

func (n *Npc) performCurrentAction() {
	if len(n.CurrentAction) == 0 {
		return
	}

	status := false
	switch n.CurrentAction {
	case "find-food":
		status = n.performFindFood()
	case "find-water":
		status = n.performFindWater()
	case "sleep":
		status = n.performSleep()
	case "dig-hole":
		status = n.performDigHole()
	case "build-fireplace":
		status = n.performBuildFireplace()
	default:
		panic(fmt.Errorf("Cant perform unknown action: %s", n.CurrentAction))
	}

	if status == true {
		log.Println(n.Name, "finished performing", n.CurrentAction)
		n.CurrentAction = ""
	}
}

func (i *Island) findActionByName(n string) actionSpec {

	for _, spec := range i.actionSpecs {
		if spec.Name == n {
			return spec
		}
	}

	panic(fmt.Errorf("cant find action: %s", n))
}

func (n *Npc) performSleep() bool {

	sleep := island.findActionByName("sleep")

	mult := 1
	if len(island.withinRadiusOfType("shelter", 0, n.Position)) > 0 {
		// XXX make use of sleeping bag or other shelter, and gain energy bonus
		log.Printf("XXX %s get sleeping bonus from nearby shelter", n.Name)
		os.Exit(0)
		mult = 3
	}
	energy := mult * sleep.Energy

	log.Printf("%s is sleeping. tiredness = %d. energy gain = %d", n.Name, n.Tiredness, energy)
	n.TimeSpentOnCurrentAction++
	n.Tiredness -= energy

	if n.Tiredness <= 0 {
		n.Tiredness = 0
		log.Printf("%s woke up, no longer tired", n.Name)
		return true
	}

	if n.TimeSpentOnCurrentAction > sleep.Duration {
		// XXX some rested-bonus buff?
		log.Printf("%s woke up, slept through full duration", n.Name)
		return true
	}

	return false
}

func (n *Npc) performFindFood() bool {

	finder := island.findActionByName("find food")

	log.Println(n.Name, "is looking for food")

	// TODO something more advanced for looking for food
	n.TimeSpentOnCurrentAction++
	if n.TimeSpentOnCurrentAction > finder.Duration {

		item := island.randomItemOfType("food")
		log.Printf("%s found a %s", n.Name, item.Name)
		n.Inventory = append(n.Inventory, item)
		return true
	}

	return false
}

func (n *Npc) performFindWater() bool {

	finder := island.findActionByName("find water")
	log.Println(n.Name, "is looking for water")

	// TODO something more advanced for looking for water
	n.TimeSpentOnCurrentAction++
	if n.TimeSpentOnCurrentAction > finder.Duration {

		item := island.randomItemOfType("drink")
		log.Printf("%s found a %s", n.Name, item.Name)
		n.Inventory = append(n.Inventory, item)

		return true
	}

	return false
}

func (n *Npc) performDigHole() bool {

	finder := island.findActionByName("dig hole")
	log.Println(n.Name, "is digging a hole")

	n.TimeSpentOnCurrentAction++
	if n.TimeSpentOnCurrentAction > finder.Duration {
		island.addNpcFromName("rabbit hole", n.Position)
		return true
	}

	return false
}

func (n *Npc) performBuildFireplace() bool {

	finder := island.findActionByName("build fireplace")
	log.Println(n.Name, "is building a fireplace")

	n.TimeSpentOnCurrentAction++
	if n.TimeSpentOnCurrentAction > finder.Duration {
		island.addNpcFromName("small fireplace", n.Position)
		return true
	}

	return false
}

func parseActionsDefinition(defFileName string) []actionSpec {

	data, err := ioutil.ReadFile(defFileName)
	if err != nil {
		panic(err)
	}

	var actions actionList
	err = yaml.Unmarshal(data, &actions)
	if err != nil {
		panic(err)
	}

	log.Infof("Read %d entries from %s", len(actions.All), defFileName)
	return actions.All
}
