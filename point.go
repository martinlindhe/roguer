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

func (p *Point) intMatches(t Point) bool {
	if int(p.X) == int(t.X) && int(p.Y) == int(t.Y) {
		return true
	}
	return false
}
