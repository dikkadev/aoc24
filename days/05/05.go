// +build exclude

package day

import (
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 5

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0

	rules := make([]Rule, 0)
	updates := make([]Update, 0)

	inUpdateSection := false
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			inUpdateSection = true
			continue
		}
		// slog.Debug(l.T)
		if !inUpdateSection {
			r := ParseRule(l.T)
			rules = append(rules, r)
			// slog.Debug("", "rule", r)
		} else {
			u := ParseUpdate(l.T)
			updates = append(updates, u)
			// slog.Debug("", "update", u)
		}
	}

	correctUpdates := make([]Update, 0)
	for _, u := range updates {
		if u.IsCorrect(rules) {
			correctUpdates = append(correctUpdates, u)
			slog.Debug("", "correct", u)
			result += u.Middle()
		}
	}

	return result
}

type Rule struct {
	Left  int
	Right int
}

func ParseRule(inp string) Rule {
	split := strings.Split(inp, "|")
	if len(split) != 2 {
		panic("Not two numbers here " + inp)
	}

	s, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(split[1])
	if err != nil {
		panic(err)
	}

	return Rule{
		Left:  s,
		Right: b,
	}
}

func (r Rule) Holds(series Update) bool {
	idxLeft := slices.Index(series, r.Left)
	if idxLeft == -1 {
		return true
	}

	idxRight := slices.Index(series, r.Right)
	if idxRight == -1 {
		return true
	}

	return idxLeft < idxRight
}

type Update []int

func ParseUpdate(inp string) Update {
	split := strings.Split(inp, ",")
	if len(split) < 2 {
		panic("Not an update str " + inp)
	}

	result := make([]int, 0)

	for _, s := range split {
		n, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		result = append(result, n)
	}

	return result
}

func (u Update) IsCorrect(rules []Rule) bool {
	for _, r := range rules {
		if !r.Holds(u) {
			return false
		}
	}
	return true
}

func (u Update) Middle() int {
	return u[len(u)/2]
}
