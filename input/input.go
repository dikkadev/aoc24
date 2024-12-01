package input

import (
	"fmt"
	"os"
	"strings"
)

type Input struct {
	data []byte
}

type AugmentedLine struct {
	T string
	N int
}

func (i *Input) AugmentedLineStream() <-chan AugmentedLine {
	ch := make(chan AugmentedLine)
	go func() {
		for i, line := range i.Lines() {
			ch <- AugmentedLine{
				T: line,
				N: i,
			}
		}
		close(ch)
	}()
	return ch
}

func (i *Input) LineStream() <-chan string {
	ch := make(chan string)
	go func() {
		for _, line := range i.Lines() {
			ch <- line
		}
		close(ch)
	}()
	return ch
}

func (i *Input) Lines() []string {
	return strings.Split(string(i.data), "\n")
}

func (i *Input) Words() []string {
	return strings.Fields(string(i.data))
}

func (i *Input) CharStream() <-chan rune {
	ch := make(chan rune)
	go func() {
		for _, r := range string(i.data) {
			ch <- r
		}
		close(ch)
	}()
	return ch
}

func (i *Input) Chars() []rune {
	chars := make([]rune, 0, len(i.data))
	for _, r := range string(i.data) {
		chars = append(chars, r)
	}
	return chars
}

func SmallFileName(day uint) string {
	return fmt.Sprintf("./input/%02d_small.input", day)
}

func RealFileName(day uint) string {
	return fmt.Sprintf("./input/%02d.input", day)
}

func NewInputForDay(fileName string) (*Input, error) {
	if _, err := os.Stat(fileName); err != nil {
		return nil, fmt.Errorf("file %s not found", fileName)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", fileName, err)
	}

	return &Input{data: data}, nil
}
