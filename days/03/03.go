// +build exclude

package day

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 3

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0
	muls := make([]Mul, 0)
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			continue
		}
		log.Debug(l.T)

		head := 0
		for head < len(l.T) {
			m, moved := ParseMul(l.T[head:])
			if moved > 0 {
				muls = append(muls, m)
				slog.Debug("", "Mul found", m, "loc", head)
				head += moved
			} else {
				head++
			}
		}
	}

	for _, m := range muls {
		result += m.Mult()
	}
	return result
}

type Mul struct {
	X, Y int
}

// mul(AAA,BBB)
// 123456789012
func ParseMul(inp string) (Mul, int) {
	result := Mul{}
	cnt := 0
	for cnt < 4 {
		r := rune(inp[cnt])
		switch cnt {
		case 0:
			if r != 'm' {
				return result, 0
			}
		case 1:
			if r != 'u' {
				return result, 0
			}
		case 2:
			if r != 'l' {
				return result, 0
			}
		case 3:
			if r != '(' {
				return result, 0
			}
		}
		cnt++
	}

	var (
		digitsStr string
		digits    strings.Builder
	)

	digitProcess := func(field *int) {
		for {
			r := rune(inp[cnt])
			if IsDigit(r) {
				digits.WriteRune(r)
				cnt++
			} else {
				break
			}
		}
		digitsStr = digits.String()
		if len(digitsStr) == 0 {
			return
		}
		x, err := strconv.Atoi(digitsStr)
		if err != nil {
			return
		}
		*field = x
	}

	digitProcess(&result.X)

	r := rune(inp[cnt])
	if r != ',' {
		return result, 0
	}
	cnt++

	digits.Reset()
	digitsStr = ""

	digitProcess(&result.Y)

	r = rune(inp[cnt])
	if r != ')' {
		return result, 0
	}

	return result, cnt
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (m Mul) Mult() int {
	return m.X * m.Y
}
