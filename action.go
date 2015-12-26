package rogue

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

// Action ...
type Action interface {
	// returns true when finished performing action
	Perform(npc *Obj) bool
}

type actionList struct {
	All []actionSpec `json:"all"`
}

type actionSpec struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Result    []string `json:"result"`
	Duration  int      `json:"duration"`
	Energy    int      `json:"energy"`
	TimeSpent int
}

func (n *Obj) performCurrentAction() {
	if n.CurrentAction == nil {
		return
	}

	status := false
	switch n.CurrentAction.Type {
	case "sleep":
		status = n.performSleep()

	case "forage":
		status = n.performForage()

	case "build":
		status = n.performBuild()

	default:
		panic(fmt.Errorf("Unknown action type: %s", n.CurrentAction.Type))
	}

	if status == true {
		log.Println(n.Name, "finished", n.CurrentAction.Name)
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

func (n *Obj) performSleep() bool {

	mult := 1
	if len(island.withinRadiusOfType("shelter", 0, n.Position)) > 0 {
		// XXX make use of sleeping bag or other shelter, and gain energy bonus
		log.Printf("XXX %s get sleeping bonus from nearby shelter", n.Name)
		os.Exit(0)
		mult = 3
	}
	energy := mult * n.CurrentAction.Energy

	log.Debugln("%s is sleeping. tiredness = %d. energy gain = %d", n.Name, n.Tiredness, energy)
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

func (n *Obj) performForage() bool {

	log.Debugln(n.Name, "is performing", n.CurrentAction.Name)

	// TODO actually move around, and dont re-visit previously foraged places
	n.CurrentAction.Duration--
	if n.CurrentAction.Duration < 0 {

		rnd := n.CurrentAction.Result[rand.Intn(len(n.CurrentAction.Result))]

		log.Printf("%s found a %s", n.Name, rnd)
		n.addToInventory(rnd)
		return true
	}

	return false
}

func (n *Obj) performBuild() bool {

	log.Debugln(n.Name, "is performing", n.CurrentAction.Name)

	n.CurrentAction.Duration--
	if n.CurrentAction.Duration < 0 {
		rnd := n.CurrentAction.Result[rand.Intn(len(n.CurrentAction.Result))]
		island.addNpcFromName(rnd, n.Position)
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
