package main

import (
	"flag"
	"log/slog"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/prettyslog"

	_ "github.com/dikkadev/aoc24/days/01"
	_ "github.com/dikkadev/aoc24/days/02"
	_ "github.com/dikkadev/aoc24/days/03"
	_ "github.com/dikkadev/aoc24/days/04"
	_ "github.com/dikkadev/aoc24/days/05"
)

var (
	small     bool
	dayNumber uint
)

func main() {
	// handler := prettyslog.NewPrettyslogHandler("AOC", prettyslog.WithSource(true))
	handler := prettyslog.NewPrettyslogHandler("AOC")
	slog.SetDefault(slog.New(handler))

	flag.BoolVar(&small, "s", false, "Use small input")
	flag.UintVar(&dayNumber, "d", 0, "Day to run")
	flag.Parse()

	if small {
		handler := prettyslog.NewPrettyslogHandler("AOC", prettyslog.WithLevel(slog.LevelDebug))
		slog.SetDefault(slog.New(handler))
	}

	slog.Info("Starting")
	day := days.Days[dayNumber]
	if day == nil {
		slog.Error("Day not found", "day", dayNumber)
		return
	}

	err := day.PrepeareInputs()
	if err != nil {
		slog.Error("Failed to prepare inputs", "error", err)
		return
	}

	slog.Info("Solving", "day", dayNumber, "small", small)

	result := day.Solve(small)

	slog.Info("Result", "result", result)

}
