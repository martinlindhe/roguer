package rogue

import "math"

// Point ...
type Point struct {
	X int
	Y int
}

func (n *Obj) distanceTo(pos Point) float64 {

	xd := float64(n.Position.X - pos.X)
	yd := float64(n.Position.Y - pos.Y)
	return math.Hypot(xd, yd)
}
