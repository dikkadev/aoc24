package day

import (
	"log/slog"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 4

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0
	longestLine := 0
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			continue
		}
		if len(l.T) > longestLine {
			longestLine = len(l.T)
		}
		// slog.Debug(l.T)
	}
	grid := NewGrid(input.Lines(), longestLine)

	for i := 0; i < len(grid.Data); i++ {
		for j := 0; j < len(grid.Data[i]); j++ {
			// result += grid.Horizontal(j, i)
			// result += grid.Vertical(j, i)
			// result += grid.Diagonal(j, i)
			result += grid.XMas(j, i)
		}
	}

	slog.Debug("", "grid", grid)

	return result
}

type Grid struct {
	Data  [][]rune
	Found [][]bool
}

func (g *Grid) String() string {
	sb := strings.Builder{}
	sb.WriteRune('\n')
	for i := range g.Data {
		for j := range g.Data[i] {
			if g.Found[i][j] {
				sb.WriteRune(g.Data[i][j])
			} else {
				sb.WriteRune('.')
			}
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func NewGrid(inp []string, width int) *Grid {
	actualInput := make([]string, 0)
	for _, l := range inp {
		if len(l) > 0 {
			actualInput = append(actualInput, l)
		}
	}
	inp = actualInput
	g := Grid{
		Data:  make([][]rune, len(inp)),
		Found: make([][]bool, len(inp)),
	}
	for i := range g.Data {
		g.Data[i] = make([]rune, width)
		g.Found[i] = make([]bool, width)
	}
	for i, l := range inp {
		if len(l) == 0 {
			continue
		}
		for j, r := range l {
			g.Data[i][j] = r
		}
	}
	return &g
}

func (g *Grid) RuneAtOffsetFromPoint(r, c, dr, dc int) rune {
	if r+dr < 0 || r+dr >= len(g.Data) {
		return 0
	}
	if c+dc < 0 || c+dc >= len(g.Data[r]) {
		return 0
	}
	return g.Data[r+dr][c+dc]
}

func (g *Grid) XMas(col, row int) int {
	finds := 0

	check := func(letters []rune) bool {
		cnt := 0
		z := g.RuneAtOffsetFromPoint(row, col, 0, 0)
		if z != letters[cnt] {
			return false
		}
		cnt++
		z = g.RuneAtOffsetFromPoint(row, col, 1, 1)
		if z != letters[cnt] {
			return false
		}
		cnt++
		z = g.RuneAtOffsetFromPoint(row, col, 2, 2)
		if z != letters[cnt] {
			return false
		}
		cnt++
		z = g.RuneAtOffsetFromPoint(row, col, 2, 0)
		if z != letters[cnt] {
			return false
		}
		cnt++
		z = g.RuneAtOffsetFromPoint(row, col, 1, 1)
		if z != letters[cnt] {
			return false
		}
		cnt++
		z = g.RuneAtOffsetFromPoint(row, col, 0, 2)
		if z != letters[cnt] {
			return false
		}
		return true
	}

	ack := func() {
		g.Found[row][col] = true
		g.Found[row+1][col+1] = true
		g.Found[row+2][col+2] = true
		g.Found[row+2][col] = true
		g.Found[row+1][col+1] = true
		g.Found[row][col+2] = true
		finds++
	}

	// M on top
	// M.M
	// .A.
	// S.S
	if check([]rune{'M', 'A', 'S', 'M', 'A', 'S'}) {
		ack()
	}

	// M on the left
	// M.S
	// .A.
	// M.S
	if check([]rune{'M', 'A', 'S', 'S', 'A', 'M'}) {
		ack()
	}

	// M on bottom
	// S.S
	// .A.
	// M.M
	if check([]rune{'S', 'A', 'M', 'S', 'A', 'M'}) {
		ack()
	}

	// M on the right
	// S.M
	// .A.
	// S.M
	if check([]rune{'S', 'A', 'M', 'M', 'A', 'S'}) {
		ack()
	}

	return finds
}

func (g *Grid) Horizontal(col, row int) int {
	finds := 0
	//regular
	isRegular := true
	for i := 0; i < 4; i++ {
		z := g.RuneAtOffsetFromPoint(row, col, 0, i)
		if z == 0 {
			isRegular = false
			break
		}
		switch i {
		case 0:
			if z != 'X' {
				isRegular = false
				break
			}
		case 1:
			if z != 'M' {
				isRegular = false
				break
			}
		case 2:
			if z != 'A' {
				isRegular = false
				break
			}
		case 3:
			if z != 'S' {
				isRegular = false
				break
			}
		}
	}
	if isRegular {
		g.Found[row][col] = true
		g.Found[row][col+1] = true
		g.Found[row][col+2] = true
		g.Found[row][col+3] = true
		finds++
	}

	//reverse
	isReverse := true
	for i := 0; i < 4; i++ {
		z := g.RuneAtOffsetFromPoint(row, col, 0, i)
		if z == 0 {
			isReverse = false
			break
		}
		switch i {
		case 0:
			if z != 'S' {
				isReverse = false
				break
			}
		case 1:
			if z != 'A' {
				isReverse = false
				break
			}
		case 2:
			if z != 'M' {
				isReverse = false
				break
			}
		case 3:
			if z != 'X' {
				isReverse = false
				break
			}
		}
	}
	if isReverse {
		g.Found[row][col] = true
		g.Found[row][col+1] = true
		g.Found[row][col+2] = true
		g.Found[row][col+3] = true
		finds++
	}

	return finds
}

func (g *Grid) Vertical(col, row int) int {
	finds := 0
	//top down
	isTopDown := true
	for i := 0; i < 4; i++ {
		z := g.RuneAtOffsetFromPoint(row, col, i, 0)
		if z == 0 {
			isTopDown = false
			break
		}
		switch i {
		case 0:
			if z != 'X' {
				isTopDown = false
				break
			}
		case 1:
			if z != 'M' {
				isTopDown = false
				break
			}
		case 2:
			if z != 'A' {
				isTopDown = false
				break
			}
		case 3:
			if z != 'S' {
				isTopDown = false
				break
			}
		}
	}
	if isTopDown {
		g.Found[row][col] = true
		g.Found[row+1][col] = true
		g.Found[row+2][col] = true
		g.Found[row+3][col] = true
		finds++
	}

	//bottom up
	isBottomUp := true
	for i := 0; i < 4; i++ {
		z := g.RuneAtOffsetFromPoint(row, col, i, 0)
		if z == 0 {
			isBottomUp = false
			break
		}
		switch i {
		case 0:
			if z != 'S' {
				isBottomUp = false
				break
			}
		case 1:
			if z != 'A' {
				isBottomUp = false
				break
			}
		case 2:
			if z != 'M' {
				isBottomUp = false
				break
			}
		case 3:
			if z != 'X' {
				isBottomUp = false
				break
			}
		}
	}
	if isBottomUp {
		g.Found[row][col] = true
		g.Found[row+1][col] = true
		g.Found[row+2][col] = true
		g.Found[row+3][col] = true
		finds++
	}

	return finds
}

func (g *Grid) Diagonal(col, row int) int {
	finds := 0
	//top left to bottom right
	isTopLeftToBottomRight := true
	for i := 0; i < 4; i++ {
		z := g.RuneAtOffsetFromPoint(row, col, i, i)
		if z == 0 {
			isTopLeftToBottomRight = false
			break
		}
		switch i {
		case 0:
			if z != 'X' {
				isTopLeftToBottomRight = false
				break
			}
		case 1:
			if z != 'M' {
				isTopLeftToBottomRight = false
				break
			}
		case 2:
			if z != 'A' {
				isTopLeftToBottomRight = false
				break
			}
		case 3:
			if z != 'S' {
				isTopLeftToBottomRight = false
				break
			}
		}
	}
	if isTopLeftToBottomRight {
		g.Found[row][col] = true
		g.Found[row+1][col+1] = true
		g.Found[row+2][col+2] = true
		g.Found[row+3][col+3] = true
		finds++
	}

	//bottom right to top left
	isBottomRightToTopLeft := true
	for i := 0; i < 4; i++ {
		z := g.RuneAtOffsetFromPoint(row, col, i, i)
		if z == 0 {
			isBottomRightToTopLeft = false
			break
		}
		switch i {
		case 0:
			if z != 'S' {
				isBottomRightToTopLeft = false
				break
			}
		case 1:
			if z != 'A' {
				isBottomRightToTopLeft = false
				break
			}
		case 2:
			if z != 'M' {
				isBottomRightToTopLeft = false
				break
			}
		case 3:
			if z != 'X' {
				isBottomRightToTopLeft = false
				break
			}
		}
	}
	if isBottomRightToTopLeft {
		g.Found[row][col] = true
		g.Found[row+1][col+1] = true
		g.Found[row+2][col+2] = true
		g.Found[row+3][col+3] = true
		finds++
	}

	//top right to bottom left
	isTopRightToBottomLeft := true
	for i := 0; i < 4; i++ {
		z := g.RuneAtOffsetFromPoint(row, col, i, -i)
		if z == 0 {
			isTopRightToBottomLeft = false
			break
		}
		switch i {
		case 0:
			if z != 'X' {
				isTopRightToBottomLeft = false
				break
			}
		case 1:
			if z != 'M' {
				isTopRightToBottomLeft = false
				break
			}
		case 2:
			if z != 'A' {
				isTopRightToBottomLeft = false
				break
			}
		case 3:
			if z != 'S' {
				isTopRightToBottomLeft = false
				break
			}
		}
	}
	if isTopRightToBottomLeft {
		g.Found[row][col] = true
		g.Found[row+1][col-1] = true
		g.Found[row+2][col-2] = true
		g.Found[row+3][col-3] = true
		finds++
	}

	//bottom left to top right
	isBottomLeftToTopRight := true
	for i := 0; i < 4; i++ {
		z := g.RuneAtOffsetFromPoint(row, col, i, -i)
		if z == 0 {
			isBottomLeftToTopRight = false
			break
		}
		switch i {
		case 0:
			if z != 'S' {
				isBottomLeftToTopRight = false
				break
			}
		case 1:
			if z != 'A' {
				isBottomLeftToTopRight = false
				break
			}
		case 2:
			if z != 'M' {
				isBottomLeftToTopRight = false
				break
			}
		case 3:
			if z != 'X' {
				isBottomLeftToTopRight = false
				break
			}
		}
	}
	if isBottomLeftToTopRight {
		g.Found[row][col] = true
		g.Found[row+1][col-1] = true
		g.Found[row+2][col-2] = true
		g.Found[row+3][col-3] = true
		finds++
	}

	return finds
}
