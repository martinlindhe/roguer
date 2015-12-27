package rogue

import (
	"fmt"
	"math"
	"math/rand"
	"os"

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
	log.Printf("%s decided to %s", n.Name, a.Name)

	a.Destination = &dst
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
		status = n.performTravel(n.CurrentAction.Energy)

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

	if math.Floor(n.Position.X) != math.Floor(n.CurrentAction.Destination.X) {
		n.Position.X = newX
	}
	if math.Floor(n.Position.Y) != math.Floor(n.CurrentAction.Destination.Y) {
		n.Position.Y = newY
	}

	log.Printf("%s is performing %s from %v to %v  with step %f dst= %v", n.Name, n.CurrentAction.Name, oldPos, n.Position, distance, n.CurrentAction.Destination)

	if n.Position.intMatches(n.CurrentAction.Destination) {
		return true
	}

	return false
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

	p := Point{0, 0}

	if *n.CurrentAction.Destination == p {
		// XXX

		list := island.withinRadiusOfType(n.CurrentAction.Result, 30, n.Position)
		if len(list) > 0 {

			rnd := list[rand.Intn(len(list))]

			n.CurrentAction.Destination = &rnd.Position
			log.Printf("%s decided to go pick up %s at %v", n.Name, rnd.Name, rnd.Position)

		}
	} else {
		// progress towards target
		check := n.performTravel(1) // XXX 1=walking speed

		// look for food at current spot
		list := island.withinRadiusOfType(n.CurrentAction.Result, 0, n.Position)

		for _, it := range list {
			log.Printf("%s picked up %s", n.Name, it.Name)
			n.addItemToInventory(*it)

			// remove spawn from world
			island.removeSpawn(it)
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

	log.Debugln(n.Name, "is performing", n.CurrentAction.Name)

	n.CurrentAction.Duration--
	if n.CurrentAction.Duration < 0 {
		island.addNpcFromName(n.CurrentAction.Result, n.Position)
		return true
	}

	return false
}
