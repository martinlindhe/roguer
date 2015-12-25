package rogue

func (n *Npc) hungerCap() int {
	return 100 + (n.Level * 5)
}

func (n *Npc) thirstCap() int {
	return 100 + (n.Level * 5)
}

func (n *Npc) tirednessCap() int {
	return 10000 + (n.Level * 100)
}

func (n *Npc) isHungry() bool {
	if n.Hunger > n.hungerCap() {
		return true
	}
	return false
}

func (n *Npc) isThirsty() bool {
	if n.Thirst > n.thirstCap() {
		return true
	}
	return false
}

func (n *Npc) isTired() bool {
	if n.Tiredness > n.tirednessCap() {
		return true
	}
	return false
}

func (n *Npc) isSleeping() bool {
	if n.CurrentAction == "sleep" {
		return true
	}
	return false
}
