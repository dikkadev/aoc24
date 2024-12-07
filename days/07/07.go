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

const DAY = 7

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0
	eqs := make([]*Equation, 0)
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			continue
		}
		// slog.Debug(l.T)
		eq := ParseEquation(l.T)
		eqs = append(eqs, eq)
		slog.Debug("parsed", "eq", eq)
		if eq.IsSolvable() {
			result += eq.Value
		}
	}

	return result
}

type Equation struct {
	Value    int
	Operands []int
}

func ParseEquation(s string) *Equation {
	split := strings.Split(s, ":")
	if len(split) != 2 {
		panic("not two from :")
	}
	v, err := strconv.Atoi(split[0])
	if err != nil {
		panic(err)
	}

	opsStr := strings.Fields(split[1])
	ops := make([]int, 0)
	for _, s := range opsStr {
		x, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		ops = append(ops, x)
	}

	return &Equation{
		Value:    v,
		Operands: ops,
	}
}

type Operation int

const (
	OpAdd Operation = iota
	OpMul
)

const OPERATIONS = 2

var AllOps = []Operation{OpAdd, OpMul}

func (e *Equation) IsSolvable() bool {
	try := func(ops []Operation) bool {
		calc := e.Operands[0]
		for i, x := range e.Operands[1:] {
			switch ops[i] {
			case OpAdd:
				calc += x
			case OpMul:
				calc *= x
			}
		}

		return calc == e.Value
	}

	permutationCount := int(math.Pow(OPERATIONS, float64(len(e.Operands)-1)))
	permutations := make([][]Operation, permutationCount)
	slog.Debug("", "permutationCount", permutationCount)
	for i := 0; i < permutationCount; i++ {
		perm := make([]Operation, len(e.Operands)-1)
		n := i
		for j := len(e.Operands) - 2; j >= 0; j-- {
			perm[j] = Operation(n % OPERATIONS)
			n /= OPERATIONS
		}
		permutations[i] = perm
		slog.Debug("", "perm", perm)
	}

	for _, perm := range permutations {
		if try(perm) {
			return true
		}
	}

	return false
}
