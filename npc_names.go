package rogue

import (
	"fmt"
	"math/rand"
	"strings"
)

func generateNpcName(t string) string {
	if t == "dwarf" {
		return generateDwarfName()
	}

	panic(fmt.Errorf("unknown npc type: %s", t))
}

// generate a dwarfish name
func generateDwarfName() string {

	a := []string{"ga", "gi", "go"}
	b := []string{"m", "n", "r", "in"}
	c := []string{"li", "dil", "la", "di"}

	res := a[rand.Intn(len(a))] + b[rand.Intn(len(b))] + c[rand.Intn(len(c))]
	return strings.Title(res)
}
