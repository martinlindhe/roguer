package rogue

func (n *Obj) ageCap() int {
	// XXX what is a sane age?
	return 100000
}

func (n *Obj) hungerCap() int {
	return 100 + (n.Level * 5)
}

func (n *Obj) thirstCap() int {
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
	if n.Age >= n.ageCap() {
		return true
	}
	return false
}
