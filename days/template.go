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
	var err error
	for l := range input.AugmentedLineStream() {
		log.Info(l.T)
	}
	return 0
}
