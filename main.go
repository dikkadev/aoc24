package main

import (
	"flag"
	"log/slog"
	"os"
	"runtime/pprof"

	_ "net/http/pprof"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/prettyslog"

	_ "github.com/dikkadev/aoc24/days/01"
	_ "github.com/dikkadev/aoc24/days/02"
	_ "github.com/dikkadev/aoc24/days/03"
	_ "github.com/dikkadev/aoc24/days/04"
	_ "github.com/dikkadev/aoc24/days/05"
	_ "github.com/dikkadev/aoc24/days/06"
	_ "github.com/dikkadev/aoc24/days/07"
	_ "github.com/dikkadev/aoc24/days/08"
	_ "github.com/dikkadev/aoc24/days/09"
	_ "github.com/dikkadev/aoc24/days/10"
	_ "github.com/dikkadev/aoc24/days/11"
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

		for _, d := range days.Days {
			if d != nil {
				d.SetLogger(slog.New(handler))
			}
		}
	}

	slog.Info("Starting")

	f, err := os.Create("cpu.prof")
	if err != nil {
		slog.Error("Could not create CPU profile", "error", err)
		return
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		slog.Error("Could not start CPU profile", "error", err)
		return
	}
	defer pprof.StopCPUProfile()

	day := days.Days[dayNumber]
	if day == nil {
		slog.Error("Day not found", "day", dayNumber)
		return
	}

	err = day.PrepeareInputs()
	if err != nil {
		slog.Error("Failed to prepare inputs", "error", err)
		return
	}

	slog.Info("Solving", "day", dayNumber, "small", small)

	result := day.Solve(small)

	slog.Info("Result", "result", result)

}
