package day

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 11

func Blink(rock int) []int {
	if rock == 0 {
		return []int{1}
	}
	strRock := strconv.Itoa(rock)
	length := len(strRock)
	if length%2 == 0 {
		// even number of digits
		rock1, _ := strconv.Atoi(strRock[:length/2])
		rock2, _ := strconv.Atoi(strRock[length/2:])
		return []int{rock1, rock2}
	} else {
		// odd number of digits
		return []int{2024 * rock}
	}
}

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	fields := strings.Fields(input.Lines()[0])

	// Part 1: direct simulation for 25 iterations
	var rocks []int
	for _, f := range fields {
		n, _ := strconv.Atoi(f)
		rocks = append(rocks, n)
	}
	for i := 0; i < 25; i++ {
		newRocks := make([]int, 0, len(rocks)*2) // capacity hint
		for _, r := range rocks {
			newRocks = append(newRocks, Blink(r)...)
		}
		rocks = newRocks
	}
	answer1 := len(rocks)

	// Part 2: use a counting approach for 75 iterations
	// reinitialize from original input
	countMap := make(map[int]int)
	for _, f := range fields {
		n, _ := strconv.Atoi(f)
		countMap[n]++
	}

	for i := 0; i < 75; i++ {
		newMap := make(map[int]int)
		for number, count := range countMap {
			results := Blink(number)
			for _, r := range results {
				newMap[r] += count
			}
		}
		countMap = newMap
	}
	answer2 := 0
	for _, c := range countMap {
		answer2 += c
	}

	log.Info("answers", "part1", answer1, "part2", answer2)
	return answer2
}

// package day
//
// import (
// 	"fmt"
// 	"log/slog"
// 	"strings"
// 	"time"
//
// 	"github.com/dikkadev/aoc24/days"
// 	"github.com/dikkadev/aoc24/input"
// )
//
// const DAY = 11
//
// func init() {
// 	days.registerday(day, solve)
// }
//
// // Pre-calculated powers of 10
// var pow10 = [10]int{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000}
//
// type StoneList struct {
// 	values []int
// }
//
// func NewStoneList() *StoneList {
// 	return &StoneList{
// 		values: make([]int, 0, 1024),
// 	}
// }
//
// func (list *StoneList) AddValue(val int) {
// 	list.values = append(list.values, val)
// }
//
// // Optimized digit counting using binary search
// func countDigits(n int) int {
// 	switch {
// 	case n < pow10[5]:
// 		switch {
// 		case n < pow10[2]:
// 			if n < pow10[1] {
// 				return 1
// 			}
// 			return 2
// 		case n < pow10[3]:
// 			return 3
// 		case n < pow10[4]:
// 			return 4
// 		default:
// 			return 5
// 		}
// 	case n < pow10[7]:
// 		if n < pow10[6] {
// 			return 6
// 		}
// 		return 7
// 	case n < pow10[9]:
// 		if n < pow10[8] {
// 			return 8
// 		}
// 		return 9
// 	default:
// 		return 10
// 	}
// }
//
// const chunkSize = 1e4 // Process in chunks to reduce memory pressure
//
// func (list *StoneList) blinkChunk(values []int, results []int) int {
// 	resultIdx := 0
// 	for _, val := range values {
// 		if val == 0 {
// 			results[resultIdx] = 1
// 			resultIdx++
// 			continue
// 		}
//
// 		digits := countDigits(val)
// 		if digits%2 == 0 {
// 			x := pow10[digits/2]
// 			left := val / x
// 			right := val % x
// 			results[resultIdx] = left
// 			results[resultIdx+1] = right
// 			resultIdx += 2
// 		} else {
// 			results[resultIdx] = val * 2024
// 			resultIdx++
// 		}
// 	}
// 	return resultIdx
// }
//
// func (list *StoneList) blink(log *slog.Logger) {
// 	inputLen := len(list.values)
// 	// Pre-allocate maximum possible size
// 	newValues := make([]int, 0, inputLen*2)
//
// 	// Process in chunks to reduce memory pressure
// 	numChunks := (inputLen + chunkSize - 1) / chunkSize
// 	results := make([]int, chunkSize*2) // Reusable buffer for chunk results
//
// 	for i := 0; i < numChunks; i++ {
// 		start := i * chunkSize
// 		end := start + chunkSize
// 		if end > inputLen {
// 			end = inputLen
// 		}
//
// 		chunk := list.values[start:end]
// 		resultCount := list.blinkChunk(chunk, results)
// 		newValues = append(newValues, results[:resultCount]...)
// 	}
//
// 	list.values = newValues
// }
//
// func (list *StoneList) String() string {
// 	var sb strings.Builder
// 	sb.WriteString("[")
// 	for i, val := range list.values {
// 		if i > 0 {
// 			sb.WriteString(" ")
// 		}
// 		sb.WriteString(fmt.Sprint(val))
// 	}
// 	sb.WriteString("]")
// 	return sb.String()
// }
//
// func Solve(input *input.Input, log *slog.Logger) int {
// 	fields := strings.Fields(input.Lines()[0])
//
// 	stoneList := NewStoneList()
// 	for _, f := range fields {
// 		n := 0
// 		for _, c := range f {
// 			n = n*10 + int(c-'0')
// 		}
// 		stoneList.AddValue(n)
// 	}
//
// 	pre := time.Now()
// 	N := 75
// 	for i := 1; i <= N; i++ {
// 		stoneList.blink(log)
// 		// if i%10 == 0 {
// 		log.Info("progress", "iteration", i, "count", len(stoneList.values), "duration", time.Since(pre))
// 		pre = time.Now()
// 		// }
// 	}
//
// 	return len(stoneList.values)
// }

// package day
//
// import (
// 	"fmt"
// 	"log/slog"
// 	"strings"
//
// 	"github.com/dikkadev/aoc24/days"
// 	"github.com/dikkadev/aoc24/input"
// )
//
// const DAY = 11
//
// // Pre-calculated powers of 10
// var pow10 = [10]int{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000}
//
// type StoneList struct {
// 	// Use a slice to maintain order and ensure each value is processed individually
// 	values []int
// }
//
// func NewStoneList() *StoneList {
// 	return &StoneList{
// 		values: make([]int, 0, 1024),
// 	}
// }
//
// func (list *StoneList) AddValue(val int) {
// 	list.values = append(list.values, val)
// }
//
// // Optimized digit counting using lookup table
// func countDigits(n int) int {
// 	switch {
// 	case n < pow10[1]:
// 		return 1
// 	case n < pow10[2]:
// 		return 2
// 	case n < pow10[3]:
// 		return 3
// 	case n < pow10[4]:
// 		return 4
// 	case n < pow10[5]:
// 		return 5
// 	case n < pow10[6]:
// 		return 6
// 	case n < pow10[7]:
// 		return 7
// 	case n < pow10[8]:
// 		return 8
// 	case n < pow10[9]:
// 		return 9
// 	default:
// 		return 10
// 	}
// }
//
// func (list *StoneList) blink(log *slog.Logger) {
// 	// Pre-allocate with expected capacity
// 	newValues := make([]int, 0, len(list.values)*2)
//
// 	// Process each value in order
// 	for _, val := range list.values {
// 		if val == 0 {
// 			newValues = append(newValues, 1)
// 			continue
// 		}
//
// 		digits := countDigits(val)
// 		if digits%2 == 0 {
// 			x := pow10[digits/2]
// 			left := val / x
// 			right := val % x
// 			newValues = append(newValues, left, right)
// 		} else {
// 			newValues = append(newValues, val*2024)
// 		}
// 	}
//
// 	// Update slice
// 	list.values = newValues
// }
//
// func (list *StoneList) String() string {
// 	var sb strings.Builder
// 	sb.WriteString("[")
// 	for i, val := range list.values {
// 		if i > 0 {
// 			sb.WriteString(" ")
// 		}
// 		sb.WriteString(fmt.Sprint(val))
// 	}
// 	sb.WriteString("]")
// 	return sb.String()
// }
//
// func init() {
// 	days.RegisterDay(DAY, Solve)
// }
//
// func Solve(input *input.Input, log *slog.Logger) int {
// 	fields := strings.Fields(input.Lines()[0])
//
// 	stoneList := NewStoneList()
// 	for _, f := range fields {
// 		n := 0
// 		for _, c := range f {
// 			n = n*10 + int(c-'0')
// 		}
// 		stoneList.AddValue(n)
// 	}
//
// 	N := 45
// 	for i := 1; i <= N; i++ {
// 		stoneList.blink(log)
// 		// if i%10 == 0 {
// 		log.Info("progress", "iteration", i, "count", len(stoneList.values))
// 		// }
// 	}
//
// 	return len(stoneList.values)
// }
