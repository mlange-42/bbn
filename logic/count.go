package logic

import (
	"fmt"
	"math/bits"
)

type countFactor struct {
	Value int
	Rule  func(count, value int) bool
}

func CountExactly(value int) Factor {
	return &countFactor{
		Value: value,
		Rule: func(count, value int) bool {
			return count == value
		},
	}
}

func CountLess(value int) Factor {
	return &countFactor{
		Value: value,
		Rule: func(count, value int) bool {
			return count < value
		},
	}
}

func CountGreater(value int) Factor {
	return &countFactor{
		Value: value,
		Rule: func(count, value int) bool {
			return count > value
		},
	}
}

func (f *countFactor) SetArgs(args ...int) error {
	if len(args) != 1 {
		return fmt.Errorf("logic operator expects 1 argument, got %d", len(args))
	}
	f.Value = args[0]
	return nil
}

func (f *countFactor) Table(given int) ([]float64, error) {
	rows := uint(1) << uint(given) // 2^given
	table := make([]float64, rows*2)

	var i uint
	for i = 0; i < rows; i++ {
		count := given - bits.OnesCount(i)
		if f.Rule(count, f.Value) {
			table[i*2] = 1
		} else {
			table[i*2+1] = 1
		}
	}
	return table, nil
}
