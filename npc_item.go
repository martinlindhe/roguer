package rogue

import "fmt"

// returns index in inventory of something edible
func (n *Npc) tryFindItemTypeInInventory(t string) (int, error) {
	// XXX find something edible in inventory, or nil

	if len(n.Inventory) == 0 {
		return 0, fmt.Errorf("Inventory is empty")
	}

	for idx, it := range n.Inventory {
		if it.Type == t {
			return idx, nil
		}
	}

	// XXX iterate and check types
	return 0, nil
}

func (n *Npc) removeFromInventory(index int) Item {

	x := n.Inventory[index]

	copy(n.Inventory[index:], n.Inventory[index+1:])
	n.Inventory[len(n.Inventory)-1] = Item{}
	n.Inventory = n.Inventory[:len(n.Inventory)-1]

	return x
}
