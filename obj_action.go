package rogue

import (
	"fmt"
	"math"
	"math/rand"

	log "github.com/Sirupsen/logrus"
)

// check if npc already has planned to do a
func (n *Obj) hasPlanned(actionName string) bool {

	if n.isPerforming(actionName) {
		return true
	}

	for _, v := range n.PlannedActions {
		if v.Name == actionName {
			return true
		}
	}
	return false
}

func (n *Obj) isPerforming(actionName string) bool {
	if n.CurrentAction != nil && n.CurrentAction.Name == actionName {
		return true
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
				fmt.Printf("XXX %s aborting, already has planned %s\n", n.Name, actionName)
				return
			}
		case Point:
			dst = it.(Point)
		default:
			panic(t)
		}
	}

	a := island.findActionByName(actionName)

	a.Destination = &dst
	n.PlannedActions = append(n.PlannedActions, a)

	if a.Destination.empty() {
		n.Announce("%s decided to %s", n, a.Name)
	} else {
		n.Announce("%s decided to %s (%s)", n, a.Name, a.Destination)
	}

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
		status = n.performTravel(n.CurrentAction.Energy)

	case "wait":
		status = n.performWait()

	default:
		panic(fmt.Errorf("Unknown action type: %s", n.CurrentAction.Type))
	}

	if status == true {
		n.Announce("%s finished %s", n.Name, n.CurrentAction.Name)
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

func (n *Obj) performWait() bool {

	n.Announce("%s is waiting", n.Name)
	n.CurrentAction.Duration--

	if n.CurrentAction.Duration < 0 {
		n.Announce("%s finished waiting", n.Name)
		return true
	}

	return false
}

func (n *Obj) performTravel(energy int) bool {

	if energy == 0 {
		panic("travel: energy is 0")
	}

	// move closer to dst
	deltaX := n.CurrentAction.Destination.X - n.Position.X
	deltaY := n.CurrentAction.Destination.Y - n.Position.Y

	angle := math.Atan2(deltaY, deltaX)
	distance := float64(energy)

	newX := n.Position.X + math.Cos(angle)*distance
	newY := n.Position.Y + math.Sin(angle)*distance

	oldPos := n.Position

	moved := false
	if math.Floor(n.Position.X) != math.Floor(n.CurrentAction.Destination.X) {
		n.Position.X = newX
		moved = true
	}
	if math.Floor(n.Position.Y) != math.Floor(n.CurrentAction.Destination.Y) {
		n.Position.Y = newY
		moved = true
	}

	if moved {
		n.Announce("%s is performing %s from %v to %v  with step %f dst= %v", n.Name, n.CurrentAction.Name, oldPos, n.Position, distance, n.CurrentAction.Destination)
	}

	if n.Position.intMatches(n.CurrentAction.Destination) {
		return true
	}

	return false
}

func (n *Obj) performSleep() bool {

	shelterType := n.preferredShelterType()

	mult := 1
	energy := mult

	if shelterType != "" {
		shelters := n.Position.spawnsByType(shelterType, 0)
		if len(shelters) > 0 {
			// give bonus from nearby shelter
			mult = shelters[0].Energy

			n.Announce("%s gets sleeping bonus %d from %s", n.Name, mult, shelters[0].Name)
		}
		energy = mult * n.CurrentAction.Energy
	}

	n.Announce("%s is sleeping (tiredness = %d, energy gain = %d)", n.Name, n.Tiredness, energy)
	n.CurrentAction.Duration--
	n.Tiredness -= energy

	if n.Tiredness <= 0 {
		n.Tiredness = 0
		n.Announce("%s woke up, no longer tired", n.Name)
		return true
	}

	if n.CurrentAction.Duration < 0 {
		// XXX some rested-bonus buff?
		n.Announce("%s woke up, slept through full duration", n.Name)
		return true
	}

	return false
}

func (n *Obj) performForage() bool {

	n.Announce("%s is performing %s", n.Name, n.CurrentAction.Name)

	p := Point{0, 0}

	if *n.CurrentAction.Destination == p {
		// XXX

		list := n.Position.spawnsByType(n.CurrentAction.Result, 30)
		if len(list) > 0 {

			rnd := list[rand.Intn(len(list))]

			n.CurrentAction.Destination = &rnd.Position
			n.Announce("%s decided to go pick up %s", n, rnd)
		}
	} else {
		// progress towards target
		check := n.performTravel(1) // XXX 1=walking speed

		// look for food at current spot
		list := n.Position.spawnsByType(n.CurrentAction.Result, 0.9)

		for _, it := range list {
			n.Announce("%s picked up %s", n.Name, it.Name)
			n.addItemToInventory(*it)

			// remove spawn from world
			island.removeSpawn(it)
		}

		// if nothing left on dst point, consider it a success!
		dstList := n.CurrentAction.Destination.spawnsByType(n.CurrentAction.Result, 0.9)
		if len(dstList) == 0 {
			return true
		}

		if check {
			// destination reached
			return true
		}
	}

	// TODO dont re-visit previously foraged places

	n.CurrentAction.Duration--
	if n.CurrentAction.Duration < 0 {
		log.Errorf("%s gave up foraging before dst reached!", n.Name)
		return true
	}

	return false
}

func (n *Obj) performBuild() bool {

	// if not at destination, move there
	// XXX 1=walking speed
	if !n.performTravel(1) {
		return false
	}

	n.Announce("%s is performing %s, duration left %d", n.Name, n.CurrentAction.Name, n.CurrentAction.Duration)

	n.CurrentAction.Duration--
	if n.CurrentAction.Duration < 0 {
		spec := island.getNpcSpecFromName(n.CurrentAction.Result)

		o := island.getNpcFromSpec(spec)
		o.Position = *n.CurrentAction.Destination
		island.addSpawn(o)

		// if object is a shelter, make it my home
		if spec.Type == "shelter" || spec.Type == "burrow" {
			n.Announce("%s has declared %s their home", n, o)
			n.Home = o
		}

		return true
	}

	return false
}
