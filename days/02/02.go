// +build exclude

package day

import (
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 2

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	safe := 0
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			continue
		}

		report := NewReport(l.T)
		if report.IsSafe() {
			safe++
		}
		slog.Debug("", "report", report, "safe?", report.IsSafe())
	}
	return safe
}

type Report struct {
	List []int
}

func NewReport(line string) Report {
	fields := strings.Fields(line)
	temp := make([]int, 0)
	for _, f := range fields {
		if len(f) == 0 {
			continue
		}

		n, err := strconv.Atoi(f)
		if err != nil {
			panic(err)
		}
		temp = append(temp, n)
	}
	return Report{
		List: temp,
	}
}

func (r *Report) IsSafe() bool {
	allInc, allDec := true, true
	prev := r.List[0]
	for i, l := range r.List {
		if i == 0 {
			continue
		}

		if l < prev {
			allInc = false
		}
		if l > prev {
			allDec = false
		}

		diff := int(math.Abs(float64(l - prev)))
		if diff < 1 || diff > 3 {
			return false
		}

		prev = l
	}

	if !(allInc || allDec) {
		return false
	}

	return true
}
