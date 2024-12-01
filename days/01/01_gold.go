package day

import (
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 1

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {

	left := List{
		Raw:    make([]int, 0),
		Sorted: make([]int, 0),
	}
	right := List{
		Raw:    make([]int, 0),
		Sorted: make([]int, 0),
	}

	for line := range input.AugmentedLineStream() {
		split := strings.Fields(line.T)
		if len(split) != 2 {
			continue
		}
		leftStr := strings.TrimSpace(split[0])
		rightStr := strings.TrimSpace(split[1])

		l, err := strconv.Atoi(leftStr)
		if err != nil {
			panic(err)
		}
		r, err := strconv.Atoi(rightStr)
		if err != nil {
			panic(err)
		}

		left.Raw = append(left.Raw, l)
		left.Sorted = append(left.Sorted, l)
		right.Raw = append(right.Raw, r)
		right.Sorted = append(right.Sorted, r)
	}

	sort.Ints(left.Sorted)
	sort.Ints(right.Sorted)

	slog.Debug("", "left", left, "right", right)

	// pairs := make([]Pair, len(left.Sorted))
	// for i := 0; i < len(left.Sorted); i++ {
	// 	pairs[i] = Pair{
	// 		Left:  left.Sorted[i],
	// 		Right: right.Sorted[i],
	// 	}
	// }
	//
	// slog.Debug("", "pairs", pairs)

	result := 0

	for _, x := range left.Raw {
		cnt := right.Count(x)
		score := x * cnt
		result += score
	}

	return result
}

type List struct {
	Raw    []int
	Sorted []int
}

func (l *List) Count(x int) int {
	var result int

	for _, y := range l.Raw {
		if x == y {
			result++
		}
	}

	return result
}

type Pair struct {
	Left  int
	Right int
}
