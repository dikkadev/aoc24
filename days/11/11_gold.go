package day

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 11

func init() {
	days.RegisterDay(DAY, Solve)
}

const CHUNK_SIZE = 1024

// Pre-calculated powers of 10
// var pow10 = [10]int{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000}
var pow10 = [19]int{1, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18}

// Chunk based storage to reduce allocations
type Chunk struct {
	values [CHUNK_SIZE]int
	next   *Chunk
	used   int
}

type StoneList struct {
	head    *Chunk
	tail    *Chunk
	count   int
	cache   map[int][2]int
	freeidx []int // Free slots in chunks
}

func NewStoneList() *StoneList {
	return &StoneList{
		cache:   make(map[int][2]int, 1024),
		freeidx: make([]int, 0, 64),
	}
}

func (list *StoneList) AddValue(val int) {
	if list.head == nil {
		list.head = &Chunk{}
		list.tail = list.head
	}

	if list.tail.used >= CHUNK_SIZE {
		list.tail.next = &Chunk{}
		list.tail = list.tail.next
	}

	list.tail.values[list.tail.used] = val
	list.tail.used++
	list.count++
}

// Optimized digit counting using lookup table
func countDigits(n int) int {
	switch {
	case n < pow10[1]:
		return 1
	case n < pow10[2]:
		return 2
	case n < pow10[3]:
		return 3
	case n < pow10[4]:
		return 4
	case n < pow10[5]:
		return 5
	case n < pow10[6]:
		return 6
	case n < pow10[7]:
		return 7
	case n < pow10[8]:
		return 8
	case n < pow10[9]:
		return 9
	case n < pow10[10]:
		return 10
	case n < pow10[11]:
		return 11
	case n < pow10[12]:
		return 12
	case n < pow10[13]:
		return 13
	case n < pow10[14]:
		return 14
	case n < pow10[15]:
		return 15
	case n < pow10[16]:
		return 16
	case n < pow10[17]:
		return 17
	default:
		return 18
	}
}

func (list *StoneList) blink(log *slog.Logger) {
	newValues := make([]int, 0, list.count*2) // Pre-allocate with expected growth

	// Process all chunks
	for chunk := list.head; chunk != nil; chunk = chunk.next {
		for i := 0; i < chunk.used; i++ {
			val := chunk.values[i]

			if val == 0 {
				newValues = append(newValues, 1)
				continue
			}

			if hit, ok := list.cache[val]; ok {
				newValues = append(newValues, hit[0])
				if hit[1] != -1 {
					newValues = append(newValues, hit[1])
				}
				continue
			}

			digits := countDigits(val)
			if digits%2 == 0 {
				x := pow10[digits/2]
				left := val / x
				right := val % x
				list.cache[val] = [2]int{left, right}
				newValues = append(newValues, left, right)
			} else {
				newVal := val * 2024
				list.cache[val] = [2]int{newVal, -1}
				newValues = append(newValues, newVal)
			}
		}
	}

	// Reset list and add new values
	list.head = &Chunk{}
	list.tail = list.head
	list.count = 0

	// Refill chunks with new values
	for _, val := range newValues {
		if list.tail.used >= CHUNK_SIZE {
			list.tail.next = &Chunk{}
			list.tail = list.tail.next
		}
		list.tail.values[list.tail.used] = val
		list.tail.used++
		list.count++
	}
}

func Solve(input *input.Input, log *slog.Logger) int {
	fields := strings.Fields(input.Lines()[0])

	stoneList := NewStoneList()
	for _, f := range fields {
		n := 0
		for _, c := range f {
			n = n*10 + int(c-'0')
		}
		stoneList.AddValue(n)
	}

	N := 45
	for i := 1; i <= N; i++ {
		stoneList.blink(log)
		// if i%10 == 0 {
		log.Info("progress", "iteration", i, "count", stoneList.count)
		// }
	}

	return stoneList.count
}

func (list *StoneList) String() string {
	var sb strings.Builder
	sb.WriteString("[")

	for chunk := list.head; chunk != nil; chunk = chunk.next {
		for i := 0; i < chunk.used; i++ {
			if i > 0 || chunk != list.head {
				sb.WriteString(" ")
			}
			sb.WriteString(fmt.Sprint(chunk.values[i]))
		}
	}

	sb.WriteString("]")
	return sb.String()
}

// var pow10 = [10]int{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000}
//
// // Pool for Stone objects to reduce GC pressure
// var stonePool = sync.Pool{
// 	New: func() interface{} {
// 		return &Stone{}
// 	},
// }
//
// func init() {
// 	days.RegisterDay(DAY, Solve)
// }
//
// func Solve(input *input.Input, log *slog.Logger) int {
// 	fields := strings.Fields(input.Lines()[0])
//
// 	// Pre-allocate cache with estimated size
// 	cache := make(map[int][2]int, len(fields)*2)
//
// 	// Initialize linked list
// 	stoneList := &StoneList{
// 		cache: cache,
// 	}
// 	var prev *Stone
// 	for _, f := range fields {
// 		stone := ParseStone(f)
// 		if prev == nil {
// 			stoneList.head = stone
// 		} else {
// 			prev.Next = stone
// 		}
// 		prev = stone
// 		stoneList.count++
// 	}
// 	stoneList.tail = prev
//
// 	start := time.Now()
// 	N := 40
// 	for i := 1; i <= N; i++ {
// 		pre := time.Now()
// 		stoneList.blink(log)
// 		log.Info("", "blink", i, "stones count", stoneList.Len(), "duration", time.Since(pre))
// 	}
// 	log.Info("", "total duration", time.Since(start))
//
// 	return stoneList.Len()
// }
//
// type Stone struct {
// 	Value int
// 	Next  *Stone
// }
//
// type StoneList struct {
// 	head  *Stone
// 	tail  *Stone
// 	count int
// 	cache map[int][2]int
// }
//
// func countDigits(n int) int {
// 	if n < pow10[1] {
// 		return 1
// 	} else if n < pow10[2] {
// 		return 2
// 	} else if n < pow10[3] {
// 		return 3
// 	} else if n < pow10[4] {
// 		return 4
// 	} else if n < pow10[5] {
// 		return 5
// 	} else if n < pow10[6] {
// 		return 6
// 	} else if n < pow10[7] {
// 		return 7
// 	} else if n < pow10[8] {
// 		return 8
// 	} else if n < pow10[9] {
// 		return 9
// 	}
// 	return 10
// }
//
// func (list *StoneList) blink(log *slog.Logger) {
// 	cnt := 0
// 	total := list.count
// 	current := list.head
//
// 	for current != nil {
// 		cnt++
// 		if cnt%1e8 == 0 {
// 			log.Info("blink", "done", cnt, "from", total, "%", float64(cnt)/float64(total)*100)
// 		}
//
// 		if current.Value == 0 {
// 			current.Value = 1
// 			current = current.Next
// 			continue
// 		}
//
// 		if hit, ok := list.cache[current.Value]; ok {
// 			current.Value = hit[0]
// 			if hit[1] != -1 {
// 				newStone := stonePool.Get().(*Stone)
// 				newStone.Value = hit[1]
// 				newStone.Next = current.Next
// 				current.Next = newStone
// 				list.count++
// 				if current == list.tail {
// 					list.tail = newStone
// 				}
// 				current = newStone.Next
// 				continue
// 			}
// 			current = current.Next
// 			continue
// 		}
//
// 		digits := countDigits(current.Value)
// 		if digits%2 == 0 {
// 			x := pow10[digits/2]
// 			left := current.Value / x
// 			right := current.Value % x
// 			list.cache[current.Value] = [2]int{left, right}
//
// 			current.Value = left
// 			newStone := stonePool.Get().(*Stone)
// 			newStone.Value = right
// 			newStone.Next = current.Next
// 			current.Next = newStone
// 			list.count++
// 			if current == list.tail {
// 				list.tail = newStone
// 			}
// 			current = newStone.Next
// 			continue
// 		}
//
// 		newVal := current.Value * 2024
// 		list.cache[current.Value] = [2]int{newVal, -1}
// 		current.Value = newVal
// 		current = current.Next
// 	}
// }
//
// // InsertAfter adds a new stone after the given stone
// func (list *StoneList) InsertAfter(after *Stone, newStone *Stone) {
// 	list.count++
// 	if after == nil {
// 		// Insert at head
// 		newStone.Next = list.head
// 		list.head = newStone
// 		if list.tail == nil {
// 			list.tail = newStone
// 		}
// 		return
// 	}
//
// 	newStone.Next = after.Next
// 	after.Next = newStone
// 	if after == list.tail {
// 		list.tail = newStone
// 	}
// }
//
// func (list *StoneList) Len() int {
// 	return list.count
// }
//
// func (list *StoneList) String() string {
// 	sb := strings.Builder{}
// 	sb.WriteString("[")
// 	for s := list.head; s != nil; s = s.Next {
// 		sb.WriteString(s.String())
// 		sb.WriteString(" ")
// 	}
// 	sb.WriteString("]")
// 	return sb.String()
// }
//
// func (s *Stone) String() string {
// 	return fmt.Sprintf("%d", s.Value)
// }
//
// func ParseStone(inp string) *Stone {
// 	n, err := strconv.Atoi(inp)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &Stone{Value: n}
// }

// ======================== JO

// func Solve(input *input.Input, log *log.Logger) int {
// 	result := 0
//
// 	stones := make([]*Stone, 0)
//
// 	fields := strings.Fields(input.Lines()[0])
// 	for _, f := range fields {
// 		stones = append(stones, ParseStone(f))
// 	}
// 	log.Debug("", "stones", stones)
//
// 	insert := func(idx int, st *Stone) {
// 		stones = append(stones, nil)
// 		copy(stones[idx+1:], stones[idx:])
// 		stones[idx] = st
// 	}
// 	blink := func() {
// 		for idx := 0; idx < len(stones); idx++ {
// 			s := stones[idx]
// 			if s.Value == 0 {
// 				s.Value = 1
// 				continue
// 			}
// 			// digits := int(math.Floor(math.Log10(float64(s.Value)))) + 1
// 			digits := 1
// 			for n := s.Value; n >= 10; n /= 10 {
// 				digits++
// 			}
// 			if digits%2 == 0 {
// 				x := int(math.Pow10(digits / 2))
// 				left := s.Value / x
// 				right := s.Value % x
// 				s.Value = left
// 				// stones.InsertAfter(&Stone{Value: right}, s)
// 				insert(idx+1, &Stone{Value: right})
// 				log.Debug("B", "stone", s, "left", left, "right", right, "result", s, "insert", s.Next)
// 				idx++
// 				continue
// 			}
// 			s.Value = s.Value * 2024
// 			log.Debug("C", "stone", s, "result", s)
// 		}
// 	}
//
// 	start := time.Now()
// 	// N := 6
// 	N := 26
// 	for i := 0; i < N; i++ {
// 		pre := time.Now()
// 		blink()
// 		log.Debug("", "N", i, "stones", stones)
// 		log.Info("", "blink", i, "stones count", len(stones), "duration", time.Since(pre))
// 	}
// 	log.Info("", "total duration", time.Since(start))
// 	result = len(stones)
//
// 	return result
// }
//
// type Stone struct {
// 	Value int
// 	Next  *Stone
// }
//
// func (s *Stone) String() string {
// 	return fmt.Sprintf("%d", s.Value)
// }
//
// func ParseStone(inp string) *Stone {
// 	n, err := strconv.Atoi(inp)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &Stone{
// 		Value: n,
// 	}
// }

// func Solve(input *input.Input, log *slog.Logger) int {
// 	fields := strings.Fields(input.Lines()[0])
//
// 	cache := make(map[int][2]int, len(fields)*2)
// 	_ = cache
//
// 	// Initialize linked list
// 	stoneList := &StoneList{}
// 	var prev *Stone
// 	for _, f := range fields {
// 		stone := ParseStone(f)
// 		if prev == nil {
// 			stoneList.head = stone
// 		} else {
// 			prev.Next = stone
// 		}
// 		prev = stone
// 		stoneList.count++
// 	}
// 	stoneList.tail = prev
//
// 	log.Debug("go", "stones", stoneList)
//
// 	blink := func() {
// 		cnt := 0
// 		total := stoneList.Len()
// 		current := stoneList.head
// 		for current != nil {
// 			cnt++
// 			if cnt%1e8 == 0 {
// 				log.Info("blink", "done", cnt, "from", total, "%", float64(cnt)/float64(total)*100)
// 			}
// 			if current.Value == 0 {
// 				log.Debug("A", "stone", current, "result", 1)
// 				current.Value = 1
// 				prev = current
// 				current = current.Next
// 				continue
// 			}
//
// 			if hit, ok := cache[current.Value]; ok {
// 				log.Debug("H", "stone", current, "hit", hit)
// 				current.Value = hit[0]
// 				if hit[1] != -1 {
// 					stoneList.InsertAfter(current, &Stone{Value: hit[1]})
// 					prev = current.Next
// 					current = current.Next.Next
// 					continue
// 				}
// 				prev = current
// 				current = current.Next
// 				continue
// 			}
//
// 			digits := countDigits(current.Value)
// 			if digits%2 == 0 {
// 				x := int(math.Pow10(digits / 2))
// 				left := current.Value / x
// 				right := current.Value % x
// 				cache[current.Value] = [2]int{left, right}
// 				log.Debug("B", "stone", current, "left", left, "right", right)
// 				current.Value = left
// 				stoneList.InsertAfter(current, &Stone{Value: right})
// 				prev = current.Next
// 				current = current.Next.Next
// 				continue
// 			}
// 			log.Debug("C", "stone", current, "result", current.Value*2024)
// 			cache[current.Value] = [2]int{current.Value * 2024, -1}
// 			current.Value = current.Value * 2024
// 			prev = current
// 			current = current.Next
// 		}
//
// 	}
//
// 	start := time.Now()
// 	// N := 6
// 	// N := 25
// 	N := 75
// 	for i := 1; i <= N; i++ {
// 		pre := time.Now()
// 		blink()
// 		log.Debug("", "N", i, "stones", stoneList)
// 		log.Info("", "blink", i, "stones count", stoneList.Len(), "duration", time.Since(pre))
// 	}
// 	log.Info("", "total duration", time.Since(start))
//
// 	return stoneList.Len()
// }
//
// type Stone struct {
// 	Value int
// 	Next  *Stone
// }
//
// // StoneList represents a linked list of stones with both head and tail pointers
// type StoneList struct {
// 	head  *Stone
// 	tail  *Stone
// 	count int
// }
//
// // InsertAfter adds a new stone after the given stone
// func (list *StoneList) InsertAfter(after *Stone, newStone *Stone) {
// 	list.count++
// 	if after == nil {
// 		// Insert at head
// 		newStone.Next = list.head
// 		list.head = newStone
// 		if list.tail == nil {
// 			list.tail = newStone
// 		}
// 		return
// 	}
//
// 	newStone.Next = after.Next
// 	after.Next = newStone
// 	if after == list.tail {
// 		list.tail = newStone
// 	}
// }
//
// func (list *StoneList) Len() int {
// 	return list.count
// }
//
// func (list *StoneList) String() string {
// 	sb := strings.Builder{}
// 	sb.WriteString("[")
// 	for s := list.head; s != nil; s = s.Next {
// 		sb.WriteString(s.String())
// 		sb.WriteString(" ")
// 	}
// 	sb.WriteString("]")
// 	return sb.String()
// }
//
// func (s *Stone) String() string {
// 	return fmt.Sprintf("%d", s.Value)
// }
//
// func ParseStone(inp string) *Stone {
// 	n, err := strconv.Atoi(inp)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &Stone{Value: n}
// }
//
// // countDigits efficiently counts the number of digits in an integer
// func countDigits(n int) int {
// 	if n == 0 {
// 		return 1
// 	}
// 	return int(math.Log10(float64(n))) + 1
// }
//
// // func Solve(input *input.Input, log *log.Logger) int {
// // 	result := 0
// //
// // 	stones := make([]*Stone, 0)
// //
// // 	fields := strings.Fields(input.Lines()[0])
// // 	for _, f := range fields {
// // 		stones = append(stones, ParseStone(f))
// // 	}
// // 	log.Debug("", "stones", stones)
// //
// // 	insert := func(idx int, st *Stone) {
// // 		stones = append(stones, nil)
// // 		copy(stones[idx+1:], stones[idx:])
// // 		stones[idx] = st
// // 	}
// // 	blink := func() {
// // 		for idx := 0; idx < len(stones); idx++ {
// // 			s := stones[idx]
// // 			if s.Value == 0 {
// // 				s.Value = 1
// // 				continue
// // 			}
// // 			// digits := int(math.Floor(math.Log10(float64(s.Value)))) + 1
// // 			digits := 1
// // 			for n := s.Value; n >= 10; n /= 10 {
// // 				digits++
// // 			}
// // 			if digits%2 == 0 {
// // 				x := int(math.Pow10(digits / 2))
// // 				left := s.Value / x
// // 				right := s.Value % x
// // 				s.Value = left
// // 				// stones.InsertAfter(&Stone{Value: right}, s)
// // 				insert(idx+1, &Stone{Value: right})
// // 				log.Debug("B", "stone", s, "left", left, "right", right, "result", s, "insert", s.Next)
// // 				idx++
// // 				continue
// // 			}
// // 			s.Value = s.Value * 2024
// // 			log.Debug("C", "stone", s, "result", s)
// // 		}
// // 	}
// //
// // 	start := time.Now()
// // 	// N := 6
// // 	N := 26
// // 	for i := 0; i < N; i++ {
// // 		pre := time.Now()
// // 		blink()
// // 		log.Debug("", "N", i, "stones", stones)
// // 		log.Info("", "blink", i, "stones count", len(stones), "duration", time.Since(pre))
// // 	}
// // 	log.Info("", "total duration", time.Since(start))
// // 	result = len(stones)
// //
// // 	return result
// // }
// //
// // type Stone struct {
// // 	Value int
// // 	Next  *Stone
// // }
// //
// // func (s *Stone) String() string {
// // 	return fmt.Sprintf("%d", s.Value)
// // }
// //
// // func ParseStone(inp string) *Stone {
// // 	n, err := strconv.Atoi(inp)
// // 	if err != nil {
// // 		panic(err)
// // 	}
// // 	return &Stone{
// // 		Value: n,
// // 	}
// // }
