package rogue

import "fmt"

// returns index in inventory of something edible
func (n *Npc) tryPickSomethingToEatFromInventory() (int, error) {
	// XXX find something edible in inventory, or nil

	if len(n.Inventory) == 0 {
		return 0, fmt.Errorf("Inventory is empty")
	}

	// XXX iterate and check types
	return 0, nil
}

func (n *Npc) removeFromInventory(index int) {

	copy(n.Inventory[index:], n.Inventory[index+1:])
	n.Inventory[len(n.Inventory)-1] = Item{}
	n.Inventory = n.Inventory[:len(n.Inventory)-1]
}
