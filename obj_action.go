package rogue

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/gobuild/log"
)

// check if npc already has planned to do a
func (o *Obj) hasPlanned(actionName string) bool {

	if o.isPerforming(actionName) {
		return true
	}

	for _, v := range o.PlannedActions {
		if v.Name == actionName {
			return true
		}
	}
	return false
}

func (o *Obj) isPerforming(actionName string) bool {
	if o.CurrentAction != nil && o.CurrentAction.Name == actionName {
		return true
	}
	return false
}

// check if npc already has planned to do a
func (o *Obj) hasPlannedType(actionType string) bool {

	if o.CurrentAction != nil && o.CurrentAction.Type == actionType {
		return true
	}

	for _, v := range o.PlannedActions {
		if v.Type == actionType {
			return true
		}
	}
	return false
}

func (o *Obj) planAction(params ...interface{}) {

	actionName := ""
	var dst Point
	for _, it := range params {
		switch t := it.(type) {
		case string:
			actionName = it.(string)
			if o.hasPlanned(actionName) {
				// fmt.Printf("XXX %s aborting, already has planned %s\n", o.Name, actionName)
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
	o.PlannedActions = append(o.PlannedActions, a)

	if a.Destination.empty() {
		o.Announce("%s decided to %s", o, a.Name)
	} else {
		o.Announce("%s decided to %s (%s)", o, a.Name, a.Destination)
	}

}

func (o *Obj) performCurrentAction() {
	if o.CurrentAction == nil {
		return
	}

	status := false
	switch o.CurrentAction.Type {
	case "sleep":
		status = o.performSleep()

	case "forage":
		status = o.performForage()

	case "build":
		status = o.performBuild()

	case "travel":
		status = o.performTravel(o.CurrentAction.Energy)

	case "wait":
		status = o.performWait()

	default:
		panic(fmt.Errorf("Unknown action type: %s", o.CurrentAction.Type))
	}

	if status == true {
		o.Announce("%s finished %s", o.Name, o.CurrentAction.Name)
		o.CurrentAction = nil
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

func (o *Obj) performWait() bool {

	o.Announce("%s is waiting", o.Name)
	o.CurrentAction.Duration--

	if o.CurrentAction.Duration < 0 {
		o.Announce("%s finished waiting", o.Name)
		return true
	}

	return false
}

func (o *Obj) performTravel(energy int) bool {

	if energy == 0 {
		log.Error("travel: energy is 0")
		return true
	}

	// move closer to dst
	deltaX := o.CurrentAction.Destination.X - o.Position.X
	deltaY := o.CurrentAction.Destination.Y - o.Position.Y

	angle := math.Atan2(deltaY, deltaX)
	distance := float64(energy)

	newX := o.Position.X + math.Cos(angle)*distance
	newY := o.Position.Y + math.Sin(angle)*distance

	oldPos := o.Position

	moved := false
	if math.Floor(o.Position.X) != math.Floor(o.CurrentAction.Destination.X) {
		o.Position.X = newX
		moved = true
	}
	if math.Floor(o.Position.Y) != math.Floor(o.CurrentAction.Destination.Y) {
		o.Position.Y = newY
		moved = true
	}

	if moved {
		o.Announce("%s is performing %s from %v to %v  with step %f dst= %v", o.Name, o.CurrentAction.Name, oldPos, o.Position, distance, o.CurrentAction.Destination)
	}

	if o.Position.intMatches(o.CurrentAction.Destination) {
		return true
	}

	return false
}

func (o *Obj) performSleep() bool {

	shelterType := o.preferredShelterType()

	mult := 1
	energy := mult

	if shelterType != "" {
		shelters := o.Position.spawnsByType(shelterType, 0)
		if len(shelters) > 0 {
			// give bonus from nearby shelter
			mult = shelters[0].Energy

			o.Announce("%s gets sleeping bonus %d from %s", o.Name, mult, shelters[0].Name)
		}
		energy = mult * o.CurrentAction.Energy
	}

	o.Announce("%s is sleeping (tiredness = %d, energy gain = %d)", o.Name, o.Tiredness, energy)
	o.CurrentAction.Duration--
	o.Tiredness -= energy

	if o.Tiredness <= 0 {
		o.Tiredness = 0
		o.Announce("%s woke up, no longer tired", o.Name)
		return true
	}

	if o.CurrentAction.Duration < 0 {
		// XXX some rested-bonus buff?
		o.Announce("%s woke up, slept through full duration", o.Name)
		return true
	}

	return false
}

func (o *Obj) performForage() bool {

	o.Announce("%s is performing %s", o.Name, o.CurrentAction.Name)

	p := Point{0, 0}

	if *o.CurrentAction.Destination == p {
		// XXX

		list := o.Position.spawnsByType(o.CurrentAction.Result, 30)
		if len(list) > 0 {

			rnd := list[rand.Intn(len(list))]

			o.CurrentAction.Destination = &rnd.Position
			o.Announce("%s decided to go pick up %s", o, rnd)
		}
	} else {
		// progress towards target
		check := o.performTravel(1) // XXX 1=walking speed

		// look for food at current spot
		list := o.Position.spawnsByType(o.CurrentAction.Result, 0.9)

		for _, it := range list {
			o.Announce("%s picked up %s", o.Name, it.Name)
			o.addItemToInventory(it)

			// remove spawn from world
			island.removeSpawn(it)
		}

		// if nothing left on dst point, consider it a success!
		dstList := o.CurrentAction.Destination.spawnsByType(o.CurrentAction.Result, 0.9)
		if len(dstList) == 0 {
			return true
		}

		if check {
			// destination reached
			return true
		}
	}

	// TODO dont re-visit previously foraged places

	o.CurrentAction.Duration--
	if o.CurrentAction.Duration < 0 {
		generalLog.Error(o.Name, "gave up foraging before dst reached!")
		return true
	}

	return false
}

func (o *Obj) performBuild() bool {

	// if not at destination, move there
	// XXX 1=walking speed
	if !o.performTravel(1) {
		return false
	}

	o.Announce("%s is performing %s, duration left %d", o.Name, o.CurrentAction.Name, o.CurrentAction.Duration)

	o.CurrentAction.Duration--
	if o.CurrentAction.Duration < 0 {
		spec := island.getNpcSpecFromName(o.CurrentAction.Result)

		home := island.getNpcFromSpec(spec)
		home.Position = *o.CurrentAction.Destination
		island.addSpawn(home)

		// if object is a shelter, make it my home
		if spec.Type == "shelter" || spec.Type == "burrow" {
			o.Announce("%s has declared %s their home", o, home)
			o.Home = home
		}

		return true
	}

	return false
}
