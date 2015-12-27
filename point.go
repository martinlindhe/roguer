package rogue

import (
	"fmt"
	"math"
)

// Point ...
type Point struct {
	X float64
	Y float64
}

func (n Point) String() string {
	return fmt.Sprintf("%f,%f", n.X, n.Y)
}

func (n *Obj) distanceTo(pos Point) float64 {

	xd := n.Position.X - pos.X
	yd := n.Position.Y - pos.Y
	return math.Hypot(xd, yd)
}

func (p *Point) intMatches(t Point) bool {
	if math.Floor(p.X) == math.Floor(t.X) && math.Floor(p.Y) == math.Floor(t.Y) {
		return true
	}
	return false
}
