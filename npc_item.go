package rogue

import "fmt"

// returns index in inventory of something edible
func (n *Obj) tryFindItemTypeInInventory(t string) (int, error) {

	if len(n.Inventory) == 0 {
		return -1, fmt.Errorf("Inventory is empty")
	}

	for idx, it := range n.Inventory {
		if it.Type == t {
			return idx, nil
		}
	}

	return -1, fmt.Errorf("%s has no items of %s in inventory", n.Name, t)
}

func (n *Obj) addToInventory(itemName string) {

	spec := island.getNpcSpecFromName(itemName)
	o := island.getNpcFromSpec(spec)

	n.Inventory = append(n.Inventory, *o)
}

func (n *Obj) removeFromInventory(index int) Obj {

	val := n.Inventory[index]

	copy(n.Inventory[index:], n.Inventory[index+1:])
	n.Inventory[len(n.Inventory)-1] = Obj{}
	n.Inventory = n.Inventory[:len(n.Inventory)-1]

	return val
}

func (n *Obj) hasItemTypeInInventory(t string) bool {

	_, err := n.tryFindItemTypeInInventory(t)
	if err != nil {
		return false
	}
	return true
}

func (n *Obj) hasItemInInventory(itemName string) bool {

	if len(n.Inventory) == 0 {
		return false
	}

	for _, it := range n.Inventory {
		if it.Name == itemName {
			return true
		}
	}
	return false
}
