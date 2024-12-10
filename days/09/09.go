// +build exclude

package day

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/dikkadev/aoc24/days"
	"github.com/dikkadev/aoc24/input"
)

const DAY = 9

func init() {
	days.RegisterDay(DAY, Solve)
}

func Solve(input *input.Input, log *slog.Logger) int {
	result := 0
	for l := range input.AugmentedLineStream() {
		if len(l.T) == 0 {
			continue
		}
		slog.Debug(l.T)
	}
	d := ParseDisk(input.Lines()[0])
	slog.Debug("", "disk", d.String())

	for {
		lastFileBlock := d.LastFile()
		if lastFileBlock == -1 {
			panic("got no more file blocks?")
		}
		firstFreeBlock := d.FirstFree()
		if firstFreeBlock == -1 {
			break
		}
		if firstFreeBlock > lastFileBlock {
			break
		}
		d.MoveBlock(lastFileBlock, firstFreeBlock)
		slog.Debug("", "disk", d.String())
	}

	for i, b := range d.Data {
		if b.T == FREE {
			continue
		}
		result += i * b.ID
	}

	return result
}

type BlockType int

const (
	FREE BlockType = iota
	FILE
)

type Block struct {
	T  BlockType
	ID int
}

func (b *Block) String() string {
	switch b.T {
	case FREE:
		return "."
	case FILE:
		return strconv.Itoa(b.ID)
	default:
		return " "
	}
}

type Disk struct {
	Data []Block
}

func (d *Disk) String() string {
	sb := strings.Builder{}
	for _, e := range d.Data {
		sb.WriteString(e.String())
	}
	return sb.String()
}

func ParseDisk(inp string) Disk {
	inp = strings.TrimSpace(inp)
	data := make([]Block, 0)

	fileId := 0
	for i, r := range inp {
		isFile := i%2 == 0
		size, err := strconv.Atoi(string(r))
		if err != nil {
			panic(err)
		}
		if isFile {
			for j := 0; j < size; j++ {
				data = append(data, Block{
					T:  FILE,
					ID: fileId,
				})
			}
			fileId++
		} else {
			for j := 0; j < size; j++ {
				data = append(data, Block{
					T: FREE,
				})
			}
		}
	}

	slog.Info("Parsed full disk")

	return Disk{
		Data: data,
	}
}

func (d *Disk) FirstFree() int {
	for i, e := range d.Data {
		if e.T == FREE {
			return i
		}
	}
	return -1
}

func (d *Disk) LastFile() int {
	for i := len(d.Data) - 1; i >= 0; i-- {
		if d.Data[i].T == FILE {
			return i
		}
	}
	return -1
}

func (d *Disk) MoveBlock(from, to int) {
	d.Data[to] = d.Data[from]
	d.Data[from] = Block{
		T: FREE,
	}
}
