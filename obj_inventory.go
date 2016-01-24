package rogue

import "fmt"

// returns index in inventory of something edible
func (o *Obj) tryFindItemTypeInInventory(t string) (int, error) {

	if len(o.Inventory) == 0 {
		return -1, fmt.Errorf("Inventory is empty")
	}

	for idx, it := range o.Inventory {
		if it.Type == t {
			return idx, nil
		}
	}

	return -1, fmt.Errorf("%s has no items of %s in inventory", o.Name, t)
}

func (o *Obj) addItemToInventory(item *Obj) {

	o.Inventory = append(o.Inventory, item)
}

func (o *Obj) addToInventory(itemName string) {

	spec := o.Island.getNpcSpecFromName(itemName)
	i := o.Island.getNpcFromSpec(spec)

	o.Inventory = append(o.Inventory, i)
}

func (o *Obj) removeFromInventory(index int) *Obj {

	i := o.Inventory[index]
	o.Inventory = append(o.Inventory[:index], o.Inventory[index+1:]...)

	return i
}

func (o *Obj) hasItemTypeInInventory(t string) bool {

	_, err := o.tryFindItemTypeInInventory(t)
	if err != nil {
		return false
	}
	return true
}

func (o *Obj) hasItemInInventory(itemName string) bool {

	if len(o.Inventory) == 0 {
		return false
	}

	for _, it := range o.Inventory {
		if it.Name == itemName {
			return true
		}
	}
	return false
}
