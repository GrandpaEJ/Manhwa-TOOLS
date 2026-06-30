package vision

// traceMaskContour finds the ordered perimeter of the mask using Moore Neighborhood contour following.
func traceMaskContour(mask [][]bool, xMin, yMin, xMax, yMax int) [][2]int {
	start := [2]int{-1, -1}
findStart:
	for y := yMin; y <= yMax; y++ {
		for x := xMin; x <= xMax; x++ {
			if mask[y][x] {
				start = [2]int{x, y}
				break findStart
			}
		}
	}
	if start[0] == -1 {
		return nil
	}

	var contour [][2]int
	curr := start
	prev := [2]int{start[0] - 1, start[1]}
	dirs := [8][2]int{{-1, 0}, {-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}}

	for {
		contour = append(contour, curr)
		startDir := 0
		dx, dy := prev[0]-curr[0], prev[1]-curr[1]
		for i, d := range dirs {
			if d[0] == dx && d[1] == dy {
				startDir = (i + 1) % 8
				break
			}
		}

		found := false
		for i := 0; i < 8; i++ {
			dirIdx := (startDir + i) % 8
			nx, ny := curr[0]+dirs[dirIdx][0], curr[1]+dirs[dirIdx][1]
			if nx >= 0 && nx < 160 && ny >= 0 && ny < 160 && mask[ny][nx] {
				prev = curr
				curr = [2]int{nx, ny}
				found = true
				break
			}
		}
		if !found || (curr[0] == start[0] && curr[1] == start[1]) || len(contour) > 1000 {
			break
		}
	}

	if len(contour) > 60 {
		step := len(contour) / 60
		var simplified [][2]int
		for i := 0; i < len(contour); i += step {
			simplified = append(simplified, contour[i])
		}
		return simplified
	}
	return contour
}
