package rollingparticle

import "math/rand"

// based on https://dance-of-death-worldgen.googlecode.com/svn/trunk/src/com/nolithius/dodworldgen/maps/RollingParticleMap.as
// http://www.nolithius.com/game-development/world-generation-breakdown

type point struct {
	X int
	Y int
}

// New returns a rolling particle covered area
func New(seed int64, width int, height int, particleLength int, innerBlur float64, outerBlur float64) [][]byte {

	edgeBiasWidth := width / 8
	edgeBiasHeight := height / 8

	rand.Seed(seed)

	tiles := make2DByteSlice(width, height)

	iterations := width * height * 2

	for i := 0; i < iterations; i++ {

		// start near the center
		sourceX := rand.Intn(width-edgeBiasWidth) + edgeBiasHeight
		sourceY := rand.Intn(height-edgeBiasHeight) + edgeBiasHeight

		for length := 0; length < particleLength; length++ {

			// between -1 and 1
			sourceX += rand.Intn(2) - 1
			sourceY += rand.Intn(2) - 1

			if sourceX < 1 || sourceX > width-2 || sourceY < 1 || sourceY > height-2 {
				break
			}

			for _, hood := range getNeighborhood(sourceX, sourceY, width, height) {
				if tiles[hood.Y][hood.X] < tiles[sourceY][sourceX] {
					sourceX = hood.X
					sourceY = hood.Y
					break
				}
			}

			tiles[sourceY][sourceX] += 3
		}
	}

	blurEdges(tiles, width, height, innerBlur, outerBlur)

	return tiles
}

/**
 * Get the Moore neighborhood (3x3, 8 surrounding tiles, minus the center tile).
 * @param int x The x position of the center of the neighborhood.
 * @param int y The y position of the center of the neighborhood.
 * @return []point An array of neighbor Points, shuffled.
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

	// return the neighborhood in no particular order
	shufflePointSlice(res)

	return res
}

// shuffle slice, without allocations
func shufflePointSlice(p []point) {

	for i := range p {
		j := rand.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
}

// "Blur" the edges of the tile array to ensure no hard edges.
func blurEdges(tiles [][]byte, width int, height int, innerBlur float64, outerBlur float64) {
	for iy := 0; iy < height; iy++ {
		for ix := 0; ix < width; ix++ {
			// Multiply the outer edge and the second outer edge by some
			// constants to ensure the world does not touch the edges.
			if ix == 0 || ix == width-1 || iy == 0 || iy == height-1 {
				tiles[iy][ix] = byte(float64(tiles[iy][ix]) * innerBlur)
			} else if ix == 1 || ix == width-2 || iy == 1 || iy == height-2 {
				tiles[iy][ix] = byte(float64(tiles[iy][ix]) * outerBlur)
			}
		}
	}
}

// returns a 2d slice in [height][width]
func make2DByteSlice(width int, height int) [][]byte {
	// allocate 2d slice
	m := make([][]byte, height)
	for i := range m {
		m[i] = make([]byte, width)
	}
	return m
}
