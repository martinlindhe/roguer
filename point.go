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

func (p Point) String() string {
	return fmt.Sprintf("%f,%f", p.X, p.Y)
}

func (p Point) Equals(p2 Point) bool {
	if p.X == p2.X && p.Y == p2.Y {
		return true
	}
	return false
}

func (n *Obj) distanceTo(pos Point) float64 {

	xd := n.Position.X - pos.X
	yd := n.Position.Y - pos.Y
	return math.Hypot(xd, yd)
}

func (p *Point) empty() bool {
	// XXX
	if p.X == 0 && p.Y == 0 {
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

// select 3x3 square of positions around n.pos, pick one at random (never p)
func (p *Point) randomNearby() (Point, error) {

	var m []Point

	for y := p.Y - 1; y <= p.Y+1; y++ {
		for x := p.X - 1; x <= p.X+1; x++ {
			if y >= 0 && y < float64(island.Height) && x >= 0 && x < float64(island.Width) {
				pp := Point{x, y}
				if !pp.Equals(*p) && island.isAboveWater(pp) {
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

func (p *Point) spawnsByName(n string, radius float64) []Obj {

	var res []Obj
	for _, o := range island.Spawns {
		if o.Name == n && o.distanceTo(*p) <= radius {
			res = append(res, *o)
		}
	}
	return res
}

func (p *Point) spawnsByType(t string, radius float64) []*Obj {

	var res []*Obj
	for _, o := range island.Spawns {
		if o.Type == t && o.distanceTo(*p) <= radius {
			res = append(res, o)
		}
	}
	//fmt.Printf("spawnsByType radius %f from %s match %s: found %d\n", radius, p, t, len(res))
	return res
}
