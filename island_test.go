package rogue

import "testing"

func BenchmarkGenerateIsland(b *testing.B) {
	for n := 0; n < b.N; n++ {

		GenerateIsland(666, 800, 400)
	}
}
