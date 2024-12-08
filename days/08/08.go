// +build exclude

package day

import (
	"log/slog"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 8

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			continue
		}
		// slog.Debug(l.T)
	}

	d := Data{
		AntennaMaps:  make(map[Frequency]Map),
		AntinodeMaps: make(map[Frequency]Map),
		Antinodes:    make([]Point, 0),
	}

	for _, f := range ALLFREQS {
		if m, ok := ParseFreqMap(f, input.Lines()); ok {
			d.AntennaMaps[f] = m
			// slog.Debug("map", "freq", f, "map", m.String(f))
			var nodeCnt int
			d.AntinodeMaps[f], nodeCnt = CalcAntinodes(m)
			if nodeCnt > 0 {
				for _, p := range d.AntinodeMaps[f].Poi {
					d.Antinodes = append(d.Antinodes, p)
				}
				// slog.Debug("antinode", "freq", f, "map", d.AntinodeMaps[f].String(ANTINODEFREQ))
				slog.Debug("combined", "freq", f, "map", d.AntennaMaps[f].StringCombined(f, d.AntinodeMaps[f]))
			}
		}
	}

	combinedAntinodes := make(map[Point]bool)
	for _, p := range d.Antinodes {
		combinedAntinodes[p] = true
	}
	result = len(combinedAntinodes)

	return result
}

type Frequency rune
type Point struct {
	X, Y int
}

var (
	ALLFREQS     = []Frequency{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	ANTINODEFREQ = Frequency('#')
)

const (
	FREQCOUNT = 26 + 26 + 10
)

type Map struct {
	Data   [][]bool
	Poi    []Point
	PoiCnt int
}

const (
	ANTENNACOLOR  = "\033[1;34m"
	ANTINODECOLOR = "\033[1;32m"

	GRAY = "\033[1;30m"

	RESET = "\033[0m"
)

func (m Map) String(freq Frequency) string {
	sb := strings.Builder{}
	sb.WriteString("\n")
	for _, row := range m.Data {
		for _, cell := range row {
			if cell {
				sb.WriteString(ANTENNACOLOR)
				sb.WriteString(string(freq))
				sb.WriteString(RESET)
			} else {
				sb.WriteString(GRAY)
				sb.WriteString(".")
				sb.WriteString(RESET)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m Map) StringCombined(freq Frequency, anti Map) string {
	sb := strings.Builder{}
	sb.WriteString("\n")
	for i, row := range m.Data {
		for j, cell := range row {
			if cell {
				sb.WriteString(ANTENNACOLOR)
				sb.WriteString(string(freq))
				sb.WriteString(RESET)
			} else if anti.Data[i][j] {
				sb.WriteString(ANTINODECOLOR)
				sb.WriteString(string(ANTINODEFREQ))
				sb.WriteString(RESET)
			} else {
				sb.WriteString(GRAY)
				sb.WriteString(".")
				sb.WriteString(RESET)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type Data struct {
	AntennaMaps  map[Frequency]Map
	AntinodeMaps map[Frequency]Map
	Antinodes    []Point
}

func ParseFreqMap(freq Frequency, inp []string) (Map, bool) {
	data := make([][]bool, len(inp)-1)
	pois := make([]Point, 0)
	for i := range data {
		data[i] = make([]bool, len(inp[0]))
	}

	found := 0
	for row, l := range inp {
		for col, field := range l {
			if Frequency(field) == freq {
				data[row][col] = true
				found++
				pois = append(pois, Point{
					X: row,
					Y: col,
				})
			}
		}
	}

	return Map{
		Data:   data,
		Poi:    pois,
		PoiCnt: found,
	}, found != 0
}

func CalcAntinodes(m Map) (Map, int) {
	data := make([][]bool, len(m.Data))
	for i := range data {
		data[i] = make([]bool, len(m.Data[0]))
	}
	nodes := make([]Point, 0)

	k := 2
	for _, p := range m.Poi {
		A := p
		for _, B := range m.Poi {
			if A == B {
				continue
			}
			poi := Point{
				X: A.X + k*(B.X-A.X),
				Y: A.Y + k*(B.Y-A.Y),
			}
			if poi.X < 0 || poi.X >= len(m.Data) || poi.Y < 0 || poi.Y >= len(m.Data[0]) {
				continue
			}
			data[poi.X][poi.Y] = true
			nodes = append(nodes, poi)
		}
	}

	return Map{
		Data:   data,
		Poi:    nodes,
		PoiCnt: len(nodes),
	}, len(nodes)
}
