package rogue

// return false if there is some construction at pos
func (i *Island) canBuildAt(pos Point) bool {

	for _, sp := range i.Spawns {
		if sp.Position == pos {
			if sp.Type == "tree" || sp.Type == "food producer" || sp.Type == "shelter" || sp.Type == "fireplace" {
				return false
			}
		}
	}

	return true
}
