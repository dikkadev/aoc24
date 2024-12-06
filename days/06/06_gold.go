package day

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 6

func init() {
	days.RegisterDay(DAY, Solve)
}

const HEADLESS = false

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0
	mapInp := make([]string, 0)
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			continue
		}
		mapInp = append(mapInp, l.T)
	}

	m, g := ParseMap(mapInp)
	log.Debug("", "guard", g.String())
	originalGuard := Guard{
		X:    g.X,
		Y:    g.Y,
		Head: g.Head,
	}

	if !HEADLESS {
		print("\033[?25l")
	}

	numWorkers := 100
	jobs := make(chan Job, numWorkers)
	results := make(chan bool, numWorkers)
	var wg sync.WaitGroup

	// Start worker pool
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, m, originalGuard, log, &wg)
	}

	// Dispatch jobs
	go func() {
		for y, row := range m {
			for x, s := range row {
				if s.Blocked || (g.X == x && g.Y == y) {
					continue
				}
				jobs <- Job{X: x, Y: y}
			}
		}
		close(jobs)
	}()

	// Calculate total jobs
	totalJobs := 0
	for y, row := range m {
		for x, s := range row {
			if s.Blocked || (g.X == x && g.Y == y) {
				continue
			}
			totalJobs++
		}
	}

	var mu sync.Mutex
	cnt := 0
	for i := 0; i < totalJobs; i++ {
		if <-results {
			mu.Lock()
			result++
			mu.Unlock()
		}
		mu.Lock()
		cnt++
		log.Info("status", "cnt", cnt, "total", totalJobs, "diff", totalJobs-cnt, "result", result)
		mu.Unlock()
	}

	wg.Wait()
	close(results)

	return result
}

type Job struct {
	X int
	Y int
}

func worker(id int, jobs <-chan Job, results chan<- bool, m Map, originalGuard Guard, log *slog.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		x, y := job.X, job.Y
		success := tryForLoop(x, y, m, originalGuard)
		results <- success
		// log.Info("Worker completed job", "workerID", id, "x", x, "y", y)
	}
}

// func tryForLoop(x, y int, m Map, originalGuard Guard) bool {
// 	timeout := 1 * time.Second
// 	deadline := time.Now().Add(timeout)
//
// 	augmentedMap := m.PlaceObstruction(x, y)
// 	thisGuard := Guard{
// 		X:    originalGuard.X,
// 		Y:    originalGuard.Y,
// 		Head: originalGuard.Head,
// 	}
//
// 	for time.Now().Before(deadline) {
// 		left, loops := thisGuard.Walk(augmentedMap, HEADLESS)
// 		if loops {
// 			return true
// 		}
// 		if left {
// 			break
// 		}
// 	}
// 	return false
// }

func tryForLoop(x, y int, m Map, originalGuard Guard) bool {
	augmentedMap := m.PlaceObstruction(x, y)
	thisGuard := Guard{
		X:    originalGuard.X,
		Y:    originalGuard.Y,
		Head: originalGuard.Head,
	}
	for i := 0; i < 1e6; i++ {
		left, loops := thisGuard.Walk(augmentedMap, HEADLESS)
		if loops {
			return true
		}
		if left {
			return false
		}
	}
	return true
}

type VisitState int

const (
	NOT_VISITED VisitState = iota
	TRAVEL_NORTH
	TRAVEL_SOUTH
	TRAVEL_WEST
	TRAVEL_EAST
)

type Space struct {
	Blocked   bool
	Visited   VisitState
	WasPlaced bool

	x, y int
}

type Map [][]*Space

// returns new fully independent map, that is the same as the og but has
func (m Map) PlaceObstruction(x, y int) Map {
	data := make([][]*Space, len(m), len(m))
	for i := range data {
		data[i] = make([]*Space, len(m[0]), len(m[0]))
	}

	for row, col := range m {
		for col, s := range col {
			data[row][col] = &Space{
				Blocked: s.Blocked,
				Visited: NOT_VISITED,
				x:       s.x,
				y:       s.y,
			}
		}
	}

	data[y][x].Blocked = true
	data[y][x].WasPlaced = true

	// slog.Debug("placed obstruction", "x", x, "y", y)
	// Display(data, Guard{X: x, Y: y, Head: NORTH})

	return Map(data)
}

type Heading rune

const (
	NORTH Heading = '^'
	SOUTH Heading = 'V'
	WEST  Heading = '<'
	EAST  Heading = '>'
)

func (h Heading) ToVisitState() VisitState {
	switch h {
	case NORTH:
		return TRAVEL_NORTH
	case SOUTH:
		return TRAVEL_SOUTH
	case WEST:
		return TRAVEL_WEST
	case EAST:
		return TRAVEL_EAST
	}
	return NOT_VISITED
}

// walk forward until bonks head, turns right)
type Guard struct {
	X, Y int
	Head Heading
}

func (g *Guard) Turn() {
	switch g.Head {
	case NORTH:
		g.Head = EAST
	case SOUTH:
		g.Head = WEST
	case WEST:
		g.Head = NORTH
	case EAST:
		g.Head = SOUTH
	}
}

func (g *Guard) WouldBonk(m Map) bool {
	//check if all indexes are safe
	if g.Y-1 < 0 || g.Y+1 >= len(m) || g.X-1 < 0 || g.X+1 >= len(m[0]) {
		return false
	}

	switch g.Head {
	case NORTH:
		return m[g.Y-1][g.X].Blocked
	case SOUTH:
		return m[g.Y+1][g.X].Blocked
	case WEST:
		return m[g.Y][g.X-1].Blocked
	case EAST:
		return m[g.Y][g.X+1].Blocked
	}
	return false
}

// returns (wouldLeave, doesLoop)
func (g *Guard) Walk(m Map, headless bool) (bool, bool) {
	var wouldLeave bool
	for !g.WouldBonk(m) {
		m[g.Y][g.X].Visited = g.Head.ToVisitState()
		switch g.Head {
		case NORTH:
			g.Y--
		case SOUTH:
			g.Y++
		case WEST:
			g.X--
		case EAST:
			g.X++
		}

		wouldLeave = g.X-1 < 0 || g.X+1 >= len(m[0]) || g.Y-1 < 0 || g.Y+1 >= len(m)
		if wouldLeave {
			// m[g.Y][g.X].Visited = true
			slog.Debug("would leave", "x", g.X, "y", g.Y)
			m[g.Y][g.X].Visited = g.Head.ToVisitState()
			break
		}

		existingState := m[g.Y][g.X].Visited
		if existingState != NOT_VISITED {
			if existingState == g.Head.ToVisitState() {
				return false, true
			}
		}
		if !headless {
			// Display(m, *g)
			// moveCursorBack(m)
			// time.Sleep(5 * time.Millisecond)
		}
		// slog.Debug("walking", "x", g.X, "y", g.Y, "head", string(g.Head), "wouldLeave", wouldLeave)
	}

	if !wouldLeave {
		g.Turn()
	}
	return wouldLeave, false
}

func (g *Guard) String() string {
	return fmt.Sprintf("Guard{X: %d, Y: %d, Head: %s}", g.X, g.Y, string(g.Head))
}

func ParseMap(inp []string) (Map, Guard) {
	data := make([][]*Space, len(inp), len(inp))
	for i := range data {
		data[i] = make([]*Space, len(inp[0]), len(inp[0]))
	}

	var g Guard

	for row, l := range inp {
		for col, r := range l {
			var space *Space
			switch r {
			case '#':
				// space = NewSpace(col, row, true, false)
				space = &Space{
					Blocked: true,
					Visited: NOT_VISITED,
					x:       col,
					y:       row,
				}
			case 'V', '^', '>', '<':
				// space = NewSpace(col, row, false, false)
				space = &Space{
					Blocked: false,
					Visited: Heading(r).ToVisitState(),
					x:       col,
					y:       row,
				}
				slog.Debug("found guard", "row", row, "col", col, "r", r)
				g = Guard{
					X:    col,
					Y:    row,
					Head: Heading(r),
				}
			case 'O':
				space = &Space{
					Blocked:   false,
					Visited:   NOT_VISITED,
					WasPlaced: true,
					x:         col,
					y:         row,
				}
				slog.Debug("found obstruction", "row", row, "col", col, "r", r)
			default:
				// space = NewSpace(col, row, false, false)
				space = &Space{
					Blocked: false,
					Visited: NOT_VISITED,
					x:       col,
					y:       row,
				}

			}
			data[row][col] = space
		}
	}

	return Map(data), g
}

const (
	GRAY   = "\033[1;30m"
	PRUPLE = "\033[1;35m"
	TEAL   = "\033[1;36m"
	RED    = "\033[1;31m"

	RESET = "\033[0m"
	CLEAR = "\033[2J"
)

func Display(m Map, g Guard) {
	// print(CLEAR)
	for y, col := range m {
		for x, s := range col {
			if g.X == x && g.Y == y {
				print(PRUPLE)
				print(string(g.Head))
				print(RESET)
				continue
			}
			// if s.Visited {
			// 	print(TEAL)
			// 	print("X")
			// 	print(RESET)
			// 	continue
			// }
			if s.Visited != NOT_VISITED {
				print(TEAL)
				switch s.Visited {
				case TRAVEL_NORTH:
					// print("W")
					print("|")
				case TRAVEL_SOUTH:
					// print("S")
					print("|")
				case TRAVEL_WEST:
					// print("A")
					print("-")
				case TRAVEL_EAST:
					// print("D")
					print("-")
				}
				print(RESET)
				continue
			}
			if s.Blocked {
				print(RESET)
				print("#")
				continue
			}
			if s.WasPlaced {
				print(RED)
				print("O")
				print(RESET)
				continue
			}
			print(GRAY)
			print(".")
			print(RESET)
		}
		print("\n")
	}
	print(RESET)
	print("\n")
}

func moveCursorBack(m Map) {
	print("\033[", len(m)+1, "A")
	print("\033[0G")
}
