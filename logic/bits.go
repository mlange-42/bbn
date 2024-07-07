package logic

import (
	"fmt"
)

type bitsFactor struct{}

// Bits combines parents as bits to an integer index.
func Bits() Factor {
	return &bitsFactor{}
}

func (f *bitsFactor) SetArgs(args ...int) error {
	if len(args) > 0 {
		return fmt.Errorf("logic operator expects zero arguments, got %d", len(args))
	}
	return nil
}

func (f *bitsFactor) Table(given int) ([]float64, error) {
	rows := uint(1) << uint(given) // 2^given
	cols := uint(1) << uint(given) // 2^given
	table := make([]float64, int(rows*cols))

	max := cols - 1

	var i uint
	for i = 0; i < rows; i++ {
		value := max - i
		table[i*cols+value] = 1
	}
	return table, nil
}
