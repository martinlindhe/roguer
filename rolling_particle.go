package rogue

import (
	"fmt"
	"math/rand"
)

/**
 * based on https://dance-of-death-worldgen.googlecode.com/svn/trunk/src/com/nolithius/dodworldgen/maps/RollingParticleMap.as
 * http://www.nolithius.com/game-development/world-generation-breakdown
 */

const (
	PARTICLE_LENGTH = 50

	OUTER_BLUR = 0.65
	INNER_BLUR = 0.90
)

type point struct {
	X int
	Y int
}

func New(width int, height int) [][]byte {

	edgeBiasWidth := width / 8
	edgeBiasHeight := height / 8

	seed := int64(123456)
	rand.Seed(seed)

	tiles := make2DByteSlice(width, height)

	iterations := width * height * 2
	// TODO: if roll results in all zero array, do a re roll, mark random seed as "bad" using redis!

	for i := 0; i < iterations; i++ {

		// Start nearer the center
		sourceX := rand.Intn(width-edgeBiasWidth) + edgeBiasHeight
		sourceY := rand.Intn(height-edgeBiasHeight) + edgeBiasHeight

		for length := 0; length < PARTICLE_LENGTH; length++ {
			sourceX += rand.Intn(2) - 1 // between -1 and 1
			sourceY += rand.Intn(2) - 1

			if sourceX < 1 || sourceX > width-2 || sourceY < 1 || sourceY > height-2 {
				break
			}

			for _, hood := range getNeighborhood(sourceX, sourceY, width, height) {
				if tiles[hood.X][hood.Y] < tiles[sourceX][sourceY] {
					sourceX = hood.X
					sourceY = hood.Y
					break
				}
			}

			tiles[sourceX][sourceY] += 3
		}
	}

	blurEdges(tiles, width, height)

	return tiles
}

/**
 * Get the Moore neighborhood (3x3, 8 surrounding tiles, minus the center tile).
 * @param int $x The x position of the center of the neighborhood.
 * @param int $y The y position of the center of the neighborhood.
 * @return Point[] An array of neighbor Points, shuffled.
 */
func getNeighborhood(x int, y int, width int, height int) []point {
	var res []point

	for a := -1; a <= 1; a++ {
		for b := -1; b <= 1; b++ {
			if a != 0 || b != 0 {
				if x+a >= 0 && x+a < width && y+b >= 0 && y+b < height {
					res = append(res, point{X: x + a, Y: y + b})
				}
			}
		}
	}

	// Return the neighborhood in no particular order
	fmt.Println(res)
	shufflePointSlice(res)
	fmt.Println(res)

	return res
}

// shuffle slice, without allocations
func shufflePointSlice(src []point) {
	for i := range src {
		j := rand.Intn(i + 1)
		src[i], src[j] = src[j], src[i]
	}
	//	return src
}

/**
 * "Blur" the edges of the tile array to ensure no hard edges.
 */
func blurEdges(tiles [][]byte, width int, height int) {
	for iy := 0; iy < height; iy++ {
		for ix := 0; ix < width; ix++ {
			// Multiply the outer edge and the second outer edge by some
			// constants to ensure the world does not touch the edges.
			if ix == 0 || ix == width-1 || iy == 0 || iy == height-1 {
				tiles[iy][ix] = byte(float64(tiles[iy][ix]) * OUTER_BLUR)
			} else if ix == 1 || ix == width-2 || iy == 1 || iy == height-2 {
				tiles[iy][ix] = byte(float64(tiles[iy][ix]) * INNER_BLUR)
			}
		}
	}
}
