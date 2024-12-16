// +build exclude

package day

import (
	"log/slog"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 12

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(in *input.Input, log *slog.Logger) int {
	lines := in.Lines()
	if len(lines) == 0 {
		return 0
	}

	grid := make([][]byte, len(lines))
	var maxCols int
	for i, l := range lines {
		b := []byte(l)
		if len(b) > maxCols {
			maxCols = len(b)
		}
		grid[i] = b
	}

	rows := len(grid)

	visited := make([][]bool, rows)
	for i := 0; i < rows; i++ {
		visited[i] = make([]bool, len(grid[i]))
	}

	dir := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	var totalPrice int

	for r := 0; r < rows; r++ {
		for c := 0; c < len(grid[r]); c++ {
			if visited[r][c] {
				continue
			}
			typ := grid[r][c]
			stack := [][2]int{{r, c}}
			visited[r][c] = true
			var cells [][2]int

			for len(stack) > 0 {
				cr, cc := stack[len(stack)-1][0], stack[len(stack)-1][1]
				stack = stack[:len(stack)-1]
				cells = append(cells, [2]int{cr, cc})

				for _, d := range dir {
					nr, nc := cr+d[0], cc+d[1]
					if nr < 0 || nr >= rows {
						continue
					}
					if nc < 0 || nc >= len(grid[nr]) {
						continue
					}
					if visited[nr][nc] {
						continue
					}
					if grid[nr][nc] == typ {
						visited[nr][nc] = true
						stack = append(stack, [2]int{nr, nc})
					}
				}
			}

			area := len(cells)
			var perimeter int
			for _, cell := range cells {
				cr, cc := cell[0], cell[1]
				for _, d := range dir {
					nr, nc := cr+d[0], cc+d[1]
					if nr < 0 || nr >= rows || nc < 0 || nc >= len(grid[nr]) || grid[nr][nc] != grid[cr][cc] {
						perimeter++
					}
				}
			}

			price := area * perimeter
			totalPrice += price
			log.Debug("region", slog.String("type", string(typ)), slog.Int("area", area), slog.Int("perim", perimeter), slog.Int("price", price))
		}
	}

	log.Debug("totalPrice", slog.Int("val", totalPrice))
	return totalPrice
}
