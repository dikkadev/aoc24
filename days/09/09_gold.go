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

	// for {
	// 	lastFileBlock := -1
	// 	firstFreeBlock := -1
	// 	for i := len(d.Data) - 1; i >= 0; i-- {
	// 		if d.Data[i].T == FILE {
	// 			lastFileBlock = i
	// 			break
	// 		}
	// 	}
	// 	if lastFileBlock == -1 {
	// 		panic("got no more file blocks?")
	// 	}
	// 	for i, e := range d.Data {
	// 		if e.T == FREE {
	// 			firstFreeBlock = i
	// 			break
	// 		}
	// 	}
	// 	if firstFreeBlock == -1 {
	// 		slog.Error("No more free blocks")
	// 		break
	// 	}
	// 	if firstFreeBlock > lastFileBlock {
	// 		slog.Debug("No more free blocks before last file block")
	// 		break
	// 	}
	// 	d.Data[firstFreeBlock] = d.Data[lastFileBlock]
	// 	d.Data[lastFileBlock] = Block{
	// 		T:  FREE,
	// 		ID: -1,
	// 	}
	// 	slog.Debug("", "disk", d.String())
	// }

	files := make([][]int, d.FileCount())
	for i, b := range d.Data {
		if b.T == FILE {
			files[b.ID] = append(files[b.ID], i)
		}
	}

	for i := len(files) - 1; i >= 0; i-- {
		//check for free space that can fit the file
		fittingFreeSpaceStart := -1
		for j := 0; j < len(d.Data); j++ {
			if d.Data[j].T == FREE {
				fittingFreeSpaceStart = j
				neededEnd := fittingFreeSpaceStart + len(files[i])
				if neededEnd > len(d.Data) {
					break
				}
				fits := true
				for k := fittingFreeSpaceStart; k < neededEnd; k++ {
					if d.Data[k].T != FREE {
						fits = false
						break
					}
				}
				if fits {
					break
				} else {
					fittingFreeSpaceStart = -1
				}
			}
		}
		if fittingFreeSpaceStart == -1 || fittingFreeSpaceStart+len(files[i]) > len(d.Data) {
			// slog.Warn("No free space for file", "file", i, "size", len(files[i]))
			continue
			// break
		}

		if fittingFreeSpaceStart > files[i][0] {
			// slog.Warn("No need to move file", "file", i, "start", fittingFreeSpaceStart, "end", fittingFreeSpaceStart+len(files[i]))
			continue
		}

		for j, k := range files[i] {
			d.Data[fittingFreeSpaceStart+j] = d.Data[k]
			d.Data[k] = Block{
				T:  FREE,
				ID: -1,
			}
		}
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

func (d *Disk) FileCount() int {
	ids := make(map[int]struct{})
	for _, b := range d.Data {
		if b.T == FILE {
			ids[b.ID] = struct{}{}
		}
	}
	return len(ids)
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
