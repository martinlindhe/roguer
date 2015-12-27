package rogue

import "math"

// Point ...
type Point struct {
	X float64
	Y float64
}

func (n *Obj) distanceTo(pos Point) float64 {

	xd := n.Position.X - pos.X
	yd := n.Position.Y - pos.Y
	return math.Hypot(xd, yd)
}
