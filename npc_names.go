package rogue

import (
	"math/rand"
	"strings"
)

func (n *dwarf) generateName() string {
	// generate a dwarfish name
	a := []string{"ga", "gi", "go"}
	b := []string{"m", "n", "r", "in"}
	c := []string{"li", "dil", "la", "di"}

	res := a[rand.Intn(len(a))] + b[rand.Intn(len(b))] + c[rand.Intn(len(c))]
	return strings.Title(res)
}
