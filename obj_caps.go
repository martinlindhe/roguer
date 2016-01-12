package rogue

import "fmt"

func (n *Obj) ageCap() int64 {
	// XXX what is a sane age?
	return 100000
}

func (n *Obj) hungerCap() int {
	return 100 + (n.Level * 5)
}

func (n *Obj) thirstCap() int {
	return 100 + (n.Level * 5)
}

func (n *Obj) coldnessCap() int {
	return 100 + (n.Level * 5)
}

func (n *Obj) tirednessCap() int {
	return 10000 + (n.Level * 100)
}

func (n *Obj) isHungry() bool {
	if n.Hunger >= n.hungerCap() {
		return true
	}
	return false
}

func (n *Obj) isThirsty() bool {
	if n.Thirst >= n.thirstCap() {
		return true
	}
	return false
}

func (n *Obj) isCold() bool {
	if n.Coldness >= n.coldnessCap() {
		return true
	}
	return false
}

func (n *Obj) isTired() bool {
	if n.Tiredness >= n.tirednessCap() {
		return true
	}
	return false
}

func (n *Obj) isSleeping() bool {
	if n.CurrentAction != nil && n.CurrentAction.Name == "sleep" {
		return true
	}
	return false
}

func (n *Obj) isAboveMaxAge() bool {
	if n.Age.Current() >= n.ageCap() {
		return true
	}
	return false
}

// used by eg. fireplace. if it is activated = it is burning
func (n *Obj) isActivated() bool {
	return n.Activated
}

// Activate ...
func (n *Obj) Activate() {
	if n.Activated == true {
		panic(fmt.Errorf("obj %s is already activated", n.Name))
	}

	n.Activated = true
}
