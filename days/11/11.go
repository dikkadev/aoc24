//go:build exclude

package day

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 11

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0

	stones := make([]*Stone, 0)

	fields := strings.Fields(input.Lines()[0])
	for _, f := range fields {
		stones = append(stones, ParseStone(f))
	}
	slog.Debug("", "stones", stones)

	insert := func(idx int, st *Stone) {
		stones = append(stones, nil)
		copy(stones[idx+1:], stones[idx:])
		stones[idx] = st
	}

	blink := func() {
		idx := 0
		for idx < len(stones) {
			s := stones[idx]
			if s.Value == 0 {
				stones[idx] = &Stone{Value: 1}
				slog.Debug("A", "stone", s, "result", stones[idx])
				idx++
				continue
			}
			str := strconv.Itoa(s.Value)
			if len(str)%2 == 0 {
				leftStr := str[:len(str)/2]
				rightStr := str[len(str)/2:]
				left, err := strconv.Atoi(leftStr)
				if err != nil {
					panic(err)
				}
				right, err := strconv.Atoi(rightStr)
				if err != nil {
					panic(err)
				}
				stones[idx] = &Stone{Value: left}
				insert(idx+1, &Stone{Value: right})
				slog.Debug("B", "stone", s, "left", left, "right", right, "result", stones[idx], "insert", stones[idx+1])
				idx += 2
				continue
			}
			stones[idx] = &Stone{Value: s.Value * 2024}
			slog.Debug("C", "stone", s, "result", stones[idx])
			idx++
		}
	}

	pre := time.Now()
	N := 25
	for i := 0; i < N; i++ {
		blink()
		slog.Debug("", "N", i, "stones", stones)
		slog.Info("", "blink", i, "stones count", len(stones), "duration", time.Since(pre))
	}
	result = len(stones)

	return result
}

type Stone struct {
	Value int
}

func (s *Stone) String() string {
	return fmt.Sprintf("%d", s.Value)
}

func ParseStone(inp string) *Stone {
	n, err := strconv.Atoi(inp)
	if err != nil {
		panic(err)
	}
	return &Stone{
		Value: n,
	}
}
