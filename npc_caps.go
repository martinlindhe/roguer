package rogue

func (n *Npc) hungerCap() int {
	return n.Level * 5
}

func (n *Npc) thirstCap() int {
	return n.Level * 100
}

func (n *Npc) tirednessCap() int {
	return n.Level * 5
}
