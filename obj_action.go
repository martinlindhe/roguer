package rogue

import (
	"fmt"
	"math/rand"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
)

// check if npc already has planned to do a
func (n *Obj) hasPlanned(actionName string) bool {

	if n.CurrentAction != nil && n.CurrentAction.Name == actionName {
		return true
	}

	for _, v := range n.PlannedActions {
		if v.Name == actionName {
			return true
		}
	}
	return false
}

// check if npc already has planned to do a
func (n *Obj) hasPlannedType(actionType string) bool {

	if n.CurrentAction != nil && n.CurrentAction.Type == actionType {
		return true
	}

	for _, v := range n.PlannedActions {
		if v.Type == actionType {
			return true
		}
	}
	return false
}

func (n *Obj) planAction(params ...interface{}) {

	actionName := ""
	var dst Point
	for _, it := range params {
		switch t := it.(type) {
		case string:
			actionName = it.(string)
			if n.hasPlanned(actionName) {
				return
			}
		case Point:
			dst = it.(Point)
		default:
			panic(t)
		}
	}

	a := island.findActionByName(actionName)
	log.Printf("%s decided to %s", n.Name, a.Name)

	a.Destination = dst
	n.PlannedActions = append(n.PlannedActions, a)
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

	case "travel":
		status = n.performTravel()

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

func (n *Obj) performTravel() bool {
	// XXX move closer to dst
	spew.Dump(n.CurrentAction)

	distanceTravelled := float64(n.CurrentAction.Energy / 100)

	deltaX := n.Position.X - n.CurrentAction.Destination.X
	deltaY := n.Position.Y - n.CurrentAction.Destination.Y

	newX := n.Position.X + distanceTravelled*deltaX
	newY := n.Position.Y + distanceTravelled*deltaY

	log.Printf("XXXX %s is performing %s from %v to %f,%f  with step %f dst= %v", n.Name, n.CurrentAction.Name, n.Position, newX, newY, distanceTravelled, n.CurrentAction.Destination)
	return true
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
