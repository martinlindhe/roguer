package rogue

import (
	"fmt"
	"math/rand"
	"strings"
)

func (o *Obj) generateName() string {
	if o.Name != "" {
		return o.Name
	}

	if o.Race == "dwarf" {
		return randomDwarfName()
	}

	panic(fmt.Errorf("unknown npc race: %s", o.Race))
}

// generate a dwarfish name
func randomDwarfName() string {

	a := []string{"ga", "gi", "go"}
	b := []string{"m", "n", "r", "in"}
	c := []string{"li", "dil", "la", "di"}

	res := a[rand.Intn(len(a))] + b[rand.Intn(len(b))] + c[rand.Intn(len(c))]
	return strings.Title(res)
}
