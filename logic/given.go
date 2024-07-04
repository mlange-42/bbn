package logic

import (
	"fmt"
)

type givenFactor struct {
	Index int
}

func Given(idx int) Factor {
	return &givenFactor{
		Index: idx,
	}
}

func (f *givenFactor) SetArgs(args ...int) error {
	if len(args) != 1 {
		return fmt.Errorf("logic operator expects 1 argument, got %d", len(args))
	}
	f.Index = args[0]
	return nil
}

func (f *givenFactor) Table(given int) ([]float64, error) {
	rows := uint(1) << uint(given) // 2^given
	table := make([]float64, rows*2)

	var i uint
	for i = 0; i < rows; i++ {
		isTrue := (i & (1 << f.Index)) == 0
		if isTrue {
			table[i*2] = 1
		} else {
			table[i*2+1] = 1
		}
	}
	return table, nil
}

type givenNotFactor struct {
	Index int
}

func GivenNot(idx int) Factor {
	return &givenNotFactor{
		Index: idx,
	}
}

func (f *givenNotFactor) SetArgs(args ...int) error {
	if len(args) != 1 {
		return fmt.Errorf("logic operator expects 1 argument, got %d", len(args))
	}
	f.Index = args[0]
	return nil
}

func (f *givenNotFactor) Table(given int) ([]float64, error) {
	rows := uint(1) << uint(given) // 2^given
	table := make([]float64, rows*2)

	var i uint
	for i = 0; i < rows; i++ {
		isFalse := (i & (1 << f.Index)) != 0
		if isFalse {
			table[i*2] = 1
		} else {
			table[i*2+1] = 1
		}
	}
	return table, nil
}

type givenExclFactor struct {
	Index int
}

func GivenExcl(idx int) Factor {
	return &givenExclFactor{
		Index: idx,
	}
}

func (f *givenExclFactor) SetArgs(args ...int) error {
	if len(args) != 1 {
		return fmt.Errorf("logic operator expects 1 argument, got %d", len(args))
	}
	f.Index = args[0]
	return nil
}

func (f *givenExclFactor) Table(given int) ([]float64, error) {
	rows := uint(1) << uint(given) // 2^given
	table := make([]float64, rows*2)

	var i uint
	for i = 0; i < rows; i++ {
		isTrue := (i & (1 << f.Index)) == 0
		if !isTrue {
			table[i*2+1] = 1
			continue
		}

		otherTrue := false
		for j := 0; j < given; j++ {
			if j == f.Index {
				continue
			}
			if (i & (1 << j)) == 0 {
				otherTrue = true
				break
			}
		}
		if !otherTrue {
			table[i*2] = 1
		} else {
			table[i*2+1] = 1
		}
	}
	return table, nil
}
