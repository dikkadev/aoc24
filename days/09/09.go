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
	disk := ParseDisk(input.Lines()[0])
	// slog.Debug("", "disk", fmt.Sprintf("%#v", disk))
	slog.Debug("", "disk", disk.String())

	return result
}

type Element interface {
	String() string
}

type Free struct {
	Size int
}

func (f Free) String() string {
	s := ""
	for i := 0; i < f.Size; i++ {
		s += "."
	}

	return s
}

type File struct {
	ID   int
	Size int
}

func (f File) String() string {
	s := ""
	for i := 0; i < f.Size; i++ {
		x := strconv.Itoa(f.ID)
		s += x
	}

	return s
}

type Disk struct {
	Data []Element
}

func (d *Disk) String() string {
	sb := strings.Builder{}
	for _, e := range d.Data {
		sb.WriteString(e.String())
	}
	return sb.String()
}

func ParseDisk(inp string) Disk {
	data := make([]Element, 0)

	fileId := 0
	for i, r := range inp {
		isFile := i%2 == 0
		size, err := strconv.Atoi(string(r))
		if err != nil {
			panic(err)
		}
		if isFile {
			data = append(data, File{
				ID:   fileId,
				Size: size,
			})
			fileId++
		} else {
			data = append(data, Free{
				Size: size,
			})
		}
	}

	return Disk{
		Data: data,
	}
}
