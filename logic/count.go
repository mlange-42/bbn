package logic

import (
	"fmt"
	"math/bits"
)

type countFactor struct {
	Value int
	Rule  func(count, value int) bool
}

// CountIs compares the number of parents with outcome True against value for equality.
func CountIs(value int) Factor {
	return &countFactor{
		Value: value,
		Rule: func(count, value int) bool {
			return count == value
		},
	}
}

// CountIs compares the number of parents with outcome True to be less than value.
func CountLess(value int) Factor {
	return &countFactor{
		Value: value,
		Rule: func(count, value int) bool {
			return count < value
		},
	}
}

// CountIs compares the number of parents with outcome True to be greater than value.
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

type countTrueFactor struct{}

// CountTrue counts the number of parents with outcome True.
func CountTrue() Factor {
	return &countTrueFactor{}
}

func (f *countTrueFactor) SetArgs(args ...int) error {
	if len(args) > 0 {
		return fmt.Errorf("logic operator expects zero arguments, got %d", len(args))
	}
	return nil
}

func (f *countTrueFactor) Table(given int) ([]float64, error) {
	rows := uint(1) << uint(given) // 2^given
	cols := given + 1
	table := make([]float64, int(rows)*cols)

	var i uint
	for i = 0; i < rows; i++ {
		count := given - bits.OnesCount(i)
		table[int(i)*cols+count] = 1
	}
	return table, nil
}

type countFalseFactor struct{}

// CountFalse counts the number of parents with outcome False.
func CountFalse() Factor {
	return &countFalseFactor{}
}

func (f *countFalseFactor) SetArgs(args ...int) error {
	if len(args) > 0 {
		return fmt.Errorf("logic operator expects zero arguments, got %d", len(args))
	}
	return nil
}

func (f *countFalseFactor) Table(given int) ([]float64, error) {
	rows := uint(1) << uint(given) // 2^given
	cols := given + 1
	table := make([]float64, int(rows)*cols)

	var i uint
	for i = 0; i < rows; i++ {
		count := bits.OnesCount(i)
		table[int(i)*cols+count] = 1
	}
	return table, nil
}
