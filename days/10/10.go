// +build exclude

package day

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 10

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0
	t := ParseTopo(input.Lines())
	slog.Debug("", "topo", t)
	tXs, tYs := t.Trailheads()
	pXs, pYs := t.Peaks()
	// slog.Debug("", "trailheads", t.StringMarked(tXs, tYs), "peaks", t.StringMarked(pXs, pYs))
	slog.Debug("", "both", t.StringMarked2(tXs, tYs, pXs, pYs))

	// startX, startY := tXs[0], tYs[0]
	// goalX, goalY := pXs[0], pYs[0]
	// found := t.Path(startX, startY, goalX, goalY)
	// slog.Debug("path?", "from", [2]int{startX, startY}, "to", [2]int{goalX, goalY}, "found", found)
	// _ = found

	for i := 0; i < len(tXs); i++ {
		for j := 0; j < len(pXs); j++ {
			if t.Path(tXs[i], tYs[i], pXs[j], pYs[j]) {
				slog.Debug("path found", "from", [2]int{tXs[i], tYs[i]}, "to", [2]int{pXs[j], pYs[j]})
				result++
			}
		}
	}
	return result
}

type Topo struct {
	Data [][]int
}

func (t Topo) String() string {
	sb := strings.Builder{}
	sb.WriteString("\n")
	for _, r := range t.Data {
		for _, c := range r {
			sb.WriteString(strconv.Itoa(c))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

const (
	MARKED = "\033[1;35m"
	SECOK  = "\033[1;32m"
	GRAY   = "\033[1;30m"
	RESET  = "\033[0m"
)

func (t *Topo) StringMarked(xs, ys []int) string {
	sb := strings.Builder{}
	sb.WriteString("\n")
	for y, r := range t.Data {
		for x, c := range r {
			marked := false
			for i := 0; i < len(xs); i++ {
				if xs[i] == x && ys[i] == y {
					marked = true
					break
				}
			}
			if marked {
				sb.WriteString(MARKED)
			} else {
				sb.WriteString(GRAY)
			}
			sb.WriteString(strconv.Itoa(c))
			sb.WriteString(RESET)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (t *Topo) StringMarked2(xs, ys, xs2, ys2 []int) string {
	sb := strings.Builder{}
	sb.WriteString("\n")
	for y, r := range t.Data {
		for x, c := range r {
			marked := false
			for i := 0; i < len(xs); i++ {
				if xs[i] == x && ys[i] == y {
					marked = true
					break
				}
			}
			if marked {
				sb.WriteString(MARKED)
			} else {
				sb.WriteString(GRAY)
			}
			marked2 := false
			for i := 0; i < len(xs2); i++ {
				if xs2[i] == x && ys2[i] == y {
					marked2 = true
					break
				}
			}
			if marked2 {
				sb.WriteString(SECOK)
			}
			sb.WriteString(strconv.Itoa(c))
			sb.WriteString(RESET)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func ParseTopo(inp []string) Topo {
	d := make([][]int, 0)
	for _, l := range inp {
		if len(l) == 0 {
			continue
		}
		r := make([]int, 0)
		for _, c := range l {
			number := int(c - '0')
			r = append(r, number)
		}
		d = append(d, r)
	}

	return Topo{
		Data: d,
	}
}

func (t *Topo) Trailheads() ([]int, []int) {
	xs := make([]int, 0)
	ys := make([]int, 0)

	for y, r := range t.Data {
		for x, c := range r {
			if c == 0 {
				xs = append(xs, x)
				ys = append(ys, y)
			}
		}
	}

	return xs, ys
}

func (t *Topo) Peaks() ([]int, []int) {
	xs := make([]int, 0)
	ys := make([]int, 0)

	for y, r := range t.Data {
		for x, c := range r {
			if c == 9 {
				xs = append(xs, x)
				ys = append(ys, y)
			}
		}
	}

	return xs, ys
}

func (t *Topo) Path(startX, startY, goalX, goalY int) bool {
	//dfs
	stack := make([][2]int, 0)
	stack = append(stack, [2]int{startX, startY})
	visited := make(map[[2]int]struct{})
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, ok := visited[current]; ok {
			continue
		}
		visitedForMarkedX := make([]int, 0)
		visitedForMarkedY := make([]int, 0)
		for k := range visited {
			visitedForMarkedX = append(visitedForMarkedX, k[0])
			visitedForMarkedY = append(visitedForMarkedY, k[1])
		}
		visitedForMarkedX = append(visitedForMarkedX, current[0])
		visitedForMarkedY = append(visitedForMarkedY, current[1])
		// slog.Debug("visiting", "current", current, "goal", [2]int{goalX, goalY}, "stack", stack, "topo", t.StringMarked2(visitedForMarkedX, visitedForMarkedY, []int{goalX}, []int{goalY}))
		visited[current] = struct{}{}
		if current[0] == goalX && current[1] == goalY {
			return true
		}
		//we can move in the direction of the value of that field is +1 of current field
		//up
		if current[1] > 0 && t.Data[current[1]-1][current[0]] == t.Data[current[1]][current[0]]+1 {
			stack = append(stack, [2]int{current[0], current[1] - 1})
		}
		//down
		if current[1] < len(t.Data)-1 && t.Data[current[1]+1][current[0]] == t.Data[current[1]][current[0]]+1 {
			stack = append(stack, [2]int{current[0], current[1] + 1})
		}
		//left
		if current[0] > 0 && t.Data[current[1]][current[0]-1] == t.Data[current[1]][current[0]]+1 {
			stack = append(stack, [2]int{current[0] - 1, current[1]})
		}
		//right
		if current[0] < len(t.Data[0])-1 && t.Data[current[1]][current[0]+1] == t.Data[current[1]][current[0]]+1 {
			stack = append(stack, [2]int{current[0] + 1, current[1]})
		}
	}
	return false
}
