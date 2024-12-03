//go:build ignore

package day

import (
	"log/slog"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 0

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			continue
		}
		slog.Debug(l.T)
	}
	return result
}
