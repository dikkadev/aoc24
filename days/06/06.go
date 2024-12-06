// +build exclude

package day

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 6

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0
	mapInp := make([]string, 0)
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			continue
		}
		// slog.Debug(l.T)
		mapInp = append(mapInp, l.T)
	}

	m, g := ParseMap(mapInp)
	slog.Debug("", "guard", g.String())

	const HEADLESS = true

	for {
		if !HEADLESS {
			Display(m, g)
			moveCursorBack(m)
		}
		left := g.Walk(m, HEADLESS)
		if left {
			break
		}
	}
	if !HEADLESS {
		Display(m, g)
	}

	for _, row := range m {
		for _, s := range row {
			if s.Visited {
				result++
			}
		}
	}

	return result
}

type Space struct {
	Blocked bool
	Visited bool

	x, y int
}

func NewSpace(x, y int, b, v bool) *Space {
	return &Space{
		Blocked: b,
		Visited: v,
		x:       x,
		y:       y,
	}
}

type Map [][]*Space

type Heading rune

const (
	NORTH Heading = '^'
	SOUTH Heading = 'V'
	WEST  Heading = '<'
	EAST  Heading = '>'
)

// walk forward until bonks head, turns right)
type Guard struct {
	X, Y int
	Head Heading
}

func (g *Guard) Turn() {
	switch g.Head {
	case NORTH:
		g.Head = EAST
	case SOUTH:
		g.Head = WEST
	case WEST:
		g.Head = NORTH
	case EAST:
		g.Head = SOUTH
	}
}

func (g *Guard) WouldBonk(m Map) bool {
	switch g.Head {
	case NORTH:
		return m[g.Y-1][g.X].Blocked
	case SOUTH:
		return m[g.Y+1][g.X].Blocked
	case WEST:
		return m[g.Y][g.X-1].Blocked
	case EAST:
		return m[g.Y][g.X+1].Blocked
	}
	return false
}

func (g *Guard) Walk(m Map, headless bool) bool {
	var wouldLeave bool
	for !g.WouldBonk(m) {
		m[g.Y][g.X].Visited = true
		switch g.Head {
		case NORTH:
			g.Y--
		case SOUTH:
			g.Y++
		case WEST:
			g.X--
		case EAST:
			g.X++
		}
		wouldLeave = g.X-1 < 0 || g.X+1 >= len(m[0]) || g.Y-1 < 0 || g.Y+1 >= len(m)
		if wouldLeave {
			m[g.Y][g.X].Visited = true
			break
		}
		if !headless {
			Display(m, *g)
			moveCursorBack(m)
			time.Sleep(80 * time.Millisecond)
		}
	}

	if !wouldLeave {
		g.Turn()
	}
	return wouldLeave
}

func (g *Guard) String() string {
	return fmt.Sprintf("Guard{X: %d, Y: %d, Head: %s}", g.X, g.Y, string(g.Head))
}

func ParseMap(inp []string) (Map, Guard) {
	data := make([][]*Space, len(inp), len(inp))
	for i := range data {
		data[i] = make([]*Space, len(inp[0]), len(inp[0]))
	}

	var g Guard

	for row, l := range inp {
		for col, r := range l {
			var space *Space
			switch r {
			case '#':
				space = NewSpace(col, row, true, false)
			case 'V', '^', '>', '<':
				space = NewSpace(col, row, false, false)
				slog.Debug("found guard", "row", row, "col", col, "r", r)
				g = Guard{
					X:    col,
					Y:    row,
					Head: Heading(r),
				}
			default:
				space = NewSpace(col, row, false, false)

			}
			data[row][col] = space
		}
	}

	return Map(data), g
}

const (
	GRAY   = "\033[1;30m"
	PRUPLE = "\033[1;35m"
	TEAL   = "\033[1;36m"

	RESET = "\033[0m"
	CLEAR = "\033[2J"
)

func Display(m Map, g Guard) {
	// print(CLEAR)
	for y, col := range m {
		for x, s := range col {
			if g.X == x && g.Y == y {
				print(PRUPLE)
				print(string(g.Head))
				print(RESET)
				continue
			}
			if s.Visited {
				print(TEAL)
				print("X")
				print(RESET)
				continue
			}
			if s.Blocked {
				print(RESET)
				print("#")
				continue
			}
			print(GRAY)
			print(".")
			print(RESET)
		}
		print("\n")
	}
	print(RESET)
	print("\n")
}

func moveCursorBack(m Map) {
	print("\033[", len(m)+1, "A")
	print("\033[0G")
}
