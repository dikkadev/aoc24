package days

import (
	"fmt"
	"log/slog"

	"github.com/dikkadev/aoc24/input"
	"github.com/dikkadev/prettyslog"
)

type Solver func(*input.Input, *slog.Logger) int

type Day struct {
	Number uint

	SmallInput *input.Input
	RealInput  *input.Input

	solve Solver

	logger *slog.Logger
}

var Days = make([]*Day, 24, 24)

func RegisterDay(number uint, solver Solver) error {
	if d := Days[number]; d != nil {
		return fmt.Errorf("day %d already registered", number)
	}

	handler := prettyslog.NewPrettyslogHandler(fmt.Sprintf("D%02d", number))
	logger := slog.New(handler)

	day := Day{
		Number: number,
		solve:  solver,
		logger: logger,
	}

	Days[number] = &day
	logger.Info("Registered", "day", number)

	return nil
}

func (d *Day) PrepeareInputs() error {
	if d.SmallInput == nil {
		smallInput, err := input.NewInputForDay(input.SmallFileName(d.Number))
		if err != nil {
			return err
		}
		d.SmallInput = smallInput
	}

	if d.RealInput == nil {
		realInput, err := input.NewInputForDay(input.RealFileName(d.Number))
		if err != nil {
			return err
		}
		d.RealInput = realInput
	}

	return nil
}

func (d *Day) Solve(small bool) int {
	var input *input.Input
	if small {
		input = d.SmallInput
	} else {
		input = d.RealInput
	}

	return d.solve(input, d.logger)
}
