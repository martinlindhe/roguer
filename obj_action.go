package rogue

import (
	"fmt"
	"math"
	"math/rand"

	log "github.com/Sirupsen/logrus"
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

	a.Destination = &dst
	n.PlannedActions = append(n.PlannedActions, a)

	if a.Destination.empty() {
		log.Printf("%s decided to %s", n, a.Name)
	} else {
		log.Printf("%s decided to %s (%s)", n, a.Name, a.Destination)
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

func (n *Obj) performWait() bool {

	log.Debugln("%s is waiting", n.Name)
	n.CurrentAction.Duration--

	if n.CurrentAction.Duration < 0 {
		log.Printf("%s finished waiting", n.Name)
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
		log.Printf("%s is performing %s from %v to %v  with step %f dst= %v", n.Name, n.CurrentAction.Name, oldPos, n.Position, distance, n.CurrentAction.Destination)
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
		shelters := island.withinRadiusOfType(shelterType, 0, n.Position)
		if len(shelters) > 0 {
			// give bonus from nearby shelter
			mult = shelters[0].Energy

			log.Printf("%s gets sleeping bonus %d from %s", n.Name, mult, shelters[0].Name)
		}
		energy = mult * n.CurrentAction.Energy
	}

	log.Printf("%s is sleeping (tiredness = %d, energy gain = %d)", n.Name, n.Tiredness, energy)
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

	p := Point{0, 0}

	if *n.CurrentAction.Destination == p {
		// XXX

		list := island.withinRadiusOfType(n.CurrentAction.Result, 30, n.Position)
		if len(list) > 0 {

			rnd := list[rand.Intn(len(list))]

			n.CurrentAction.Destination = &rnd.Position
			log.Printf("%s decided to go pick up %s", n, rnd)
		}
	} else {
		// progress towards target
		check := n.performTravel(1) // XXX 1=walking speed

		// look for food at current spot
		list := island.withinRadiusOfType(n.CurrentAction.Result, 0.9, n.Position)

		for _, it := range list {
			log.Printf("%s picked up %s", n.Name, it.Name)
			n.addItemToInventory(*it)

			// remove spawn from world
			island.removeSpawn(it)
		}

		// if nothing left on dst point, consider it a success!
		dstList := island.withinRadiusOfType(n.CurrentAction.Result, 0.9, *n.CurrentAction.Destination)
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
	if n.performTravel(1) {

		log.Debugln(n.Name, "is performing", n.CurrentAction.Name)

		n.CurrentAction.Duration--
		if n.CurrentAction.Duration < 0 {
			island.addNpcFromName(n.CurrentAction.Result, *n.CurrentAction.Destination)
			return true
		}
	}

	return false
}
