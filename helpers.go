package rogue

// returns a 2d slice in [height][width]
func make2DByteSlice(width int, height int) [][]byte {
	// allocate 2d slice
	m := make([][]byte, height)
	for i := range m {
		m[i] = make([]byte, width)
	}
	return m
}
