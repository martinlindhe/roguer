package rogue

import (
	"fmt"
	"math"
	"math/rand"
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

func (n *Point) empty() bool {
	// XXX
	if n.X == 0 && n.Y == 0 {
		return true
	}
	return false
}

func (p *Point) intMatches(t *Point) bool {
	if math.Floor(p.X) == math.Floor(t.X) && math.Floor(p.Y) == math.Floor(t.Y) {
		return true
	}
	return false
}

func (p *Point) randomNearby() Point {
	// select 3x3 square of positions around n.pos, pick one at random
	var m []Point

	for y := p.Y - 1; y <= p.Y+1; y++ {
		for x := p.X - 1; x <= p.X+1; x++ {
			pp := Point{x, y}
			if island.isAboveWater(pp) {
				m = append(m, pp)
			}
		}
	}

	// select something by random
	return m[rand.Intn(len(m))]
}
