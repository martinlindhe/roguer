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

func (p *Point) randomNearby() (Point, error) {
	// select 3x3 square of positions around n.pos, pick one at random (never p)
	var m []Point

	for y := p.Y - 1; y <= p.Y+1; y++ {
		for x := p.X - 1; x <= p.X+1; x++ {
			if y >= 0 && y < float64(island.Height) && x >= 0 && x < float64(island.Width) {
				pp := Point{x, y}
				if pp != *p && island.isAboveWater(pp) {
					m = append(m, pp)
				}
			}
		}
	}

	if len(m) == 0 {
		empty := Point{}
		return empty, fmt.Errorf("Cant find nearby points to %s", p)
	}

	// select something by random
	return m[rand.Intn(len(m))], nil
}

func (p *Point) isNearby(pos Point) bool {
	distance := float64(5)
	absX := math.Abs(p.X - pos.X)
	absY := math.Abs(p.Y - pos.Y)
	if absX < distance && absY < distance {
		return true
	}
	return false
}

func (pos *Point) spawnsByName(n string, radius float64) []Obj {

	var res []Obj
	for _, o := range island.Spawns {
		if o.Name == n && o.distanceTo(*pos) <= radius {
			res = append(res, *o)
		}
	}
	return res
}

func (pos *Point) spawnsByType(t string, radius float64) []*Obj {

	var res []*Obj
	for _, o := range island.Spawns {

		// log.Printf("XXX %s at dist %f (radius %f)", npc, npc.distanceTo(pos), radius)
		if o.Type == t && o.distanceTo(*pos) <= radius {
			res = append(res, o)
		}
	}
	return res
}
