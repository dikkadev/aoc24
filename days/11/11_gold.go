package day

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 11

// We'll keep the current stones in a file on disk, to avoid huge RAM usage.
// Each stone is a line in the file.
// On each blink, we read line by line, process, and write to a new file.
// This is slow but should control memory usage.

func multiplyBy2024(numStr string) string {
	const factor = 2024
	digits := make([]int, len(numStr))
	for i, c := range numStr {
		digits[i] = int(c - '0')
	}

	carry := 0
	for i := len(digits) - 1; i >= 0; i-- {
		res := digits[i]*factor + carry
		digits[i] = res % 10
		carry = res / 10
	}

	for carry > 0 {
		digits = append([]int{carry % 10}, digits...)
		carry /= 10
	}

	var sb strings.Builder
	for _, d := range digits {
		sb.WriteByte(byte('0' + d))
	}
	return sb.String()
}

func splitEvenDigits(numStr string) (string, string) {
	half := len(numStr) / 2
	left := strings.TrimLeft(numStr[:half], "0")
	right := strings.TrimLeft(numStr[half:], "0")

	if left == "" {
		left = "0"
	}
	if right == "" {
		right = "0"
	}
	return left, right
}

func blinkFile(inFile, outFile string) (int, error) {
	inF, err := os.Open(inFile)
	if err != nil {
		return 0, err
	}
	defer inF.Close()

	outF, err := os.Create(outFile)
	if err != nil {
		return 0, err
	}
	defer outF.Close()

	sc := bufio.NewScanner(inF)
	bw := bufio.NewWriter(outF)

	count := 0
	for sc.Scan() {
		stone := sc.Text()
		if stone == "0" {
			bw.WriteString("1\n")
			count++
			continue
		}

		length := len(stone)
		if length%2 == 0 {
			left, right := splitEvenDigits(stone)
			bw.WriteString(left + "\n")
			bw.WriteString(right + "\n")
			count += 2
		} else {
			newVal := multiplyBy2024(stone)
			bw.WriteString(newVal + "\n")
			count++
		}
	}

	if err := sc.Err(); err != nil {
		return 0, err
	}

	if err := bw.Flush(); err != nil {
		return 0, err
	}

	return count, nil
}

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	fields := strings.Fields(input.Lines()[0])

	// Create a temp directory for stone files
	tmpDir, err := os.MkdirTemp("", "stones")
	if err != nil {
		log.Error("failed to create temp dir", "error", err)
		return -1
	}
	defer os.RemoveAll(tmpDir) // cleanup

	currentFile := filepath.Join(tmpDir, "stones_0.txt")
	{
		f, err := os.Create(currentFile)
		if err != nil {
			log.Error("failed to create initial file", "error", err)
			return -1
		}
		for _, fld := range fields {
			f.WriteString(fld + "\n")
		}
		f.Close()
	}

	pre := time.Now()
	N := 75 // we ultimately want iteration 75
	var count int
	for i := 1; i <= N; i++ {
		nextFile := filepath.Join(tmpDir, fmt.Sprintf("stones_%d.txt", i))
		count, err = blinkFile(currentFile, nextFile)
		if err != nil {
			log.Error("error blinking", "iteration", i, "error", err)
			return -1
		}
		log.Info("progress", "iteration", i, "count", count, "duration", time.Since(pre))
		os.Remove(currentFile)
		currentFile = nextFile
	}

	return count
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
