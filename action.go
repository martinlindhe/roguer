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
	if n.CurrentAction == nil {
		return
	}

	// XXX react on TYPE, be more generic:
	status := false
	switch n.CurrentAction.Name {
	case "find food":
		status = n.performFindFood()
	case "find water":
		status = n.performFindWater()
	case "sleep":
		status = n.performSleep()
	case "dig small hole":
		status = n.performDigHole()
	case "build small fireplace":
		status = n.performBuildFireplace()
	default:
		panic(fmt.Errorf("Cant perform unknown action: %s", n.CurrentAction))
	}

	if status == true {
		log.Println(n.Name, "finished performing", n.CurrentAction)
		n.CurrentAction = nil
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

	mult := 1
	if len(island.withinRadiusOfType("shelter", 0, n.Position)) > 0 {
		// XXX make use of sleeping bag or other shelter, and gain energy bonus
		log.Printf("XXX %s get sleeping bonus from nearby shelter", n.Name)
		os.Exit(0)
		mult = 3
	}
	energy := mult * n.CurrentAction.Energy

	log.Printf("%s is sleeping. tiredness = %d. energy gain = %d", n.Name, n.Tiredness, energy)
	n.CurrentAction.Duration--
	n.Tiredness -= energy

	if n.Tiredness <= 0 {
		n.Tiredness = 0
		log.Printf("%s woke up, no longer tired", n.Name)
		return true
	}

	if n.CurrentAction.Duration < 0 {
		// XXX some rested-bonus buff?
		log.Printf("%s woke up, slept through full duration", n.Name)
		return true
	}

	return false
}

func (n *Npc) performFindFood() bool {

	log.Println(n.Name, "is looking for food")

	// TODO something more advanced for looking for food
	n.CurrentAction.Duration--
	if n.CurrentAction.Duration < 0 {

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
	n.CurrentAction.Duration--
	if n.CurrentAction.Duration < 0 {

		item := island.randomItemOfType("drink")
		log.Printf("%s found a %s", n.Name, item.Name)
		n.Inventory = append(n.Inventory, item)

		return true
	}

	return false
}

func (n *Npc) performDigHole() bool {

	log.Println(n.Name, "is digging a hole")

	n.CurrentAction.Duration--
	if n.CurrentAction.Duration < 0 {
		island.addNpcFromName("rabbit hole", n.Position)
		return true
	}

	return false
}

func (n *Npc) performBuildFireplace() bool {

	log.Println(n.Name, "is building a fireplace")

	n.CurrentAction.Duration--
	if n.CurrentAction.Duration < 0 {
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
