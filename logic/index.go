package logic

import (
	"fmt"
)

type indexFactor struct {
	Index int
}

func Index(idx int) Factor {
	return &indexFactor{
		Index: idx,
	}
}

func (f *indexFactor) SetArgs(args ...int) error {
	if len(args) != 1 {
		return fmt.Errorf("logic operator expects 1 argument, got %d", len(args))
	}
	f.Index = args[0]
	return nil
}

func (f *indexFactor) Table(given int) ([]float64, error) {
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

type indexNotFactor struct {
	Index int
}

func IndexNot(idx int) Factor {
	return &indexNotFactor{
		Index: idx,
	}
}

func (f *indexNotFactor) SetArgs(args ...int) error {
	if len(args) != 1 {
		return fmt.Errorf("logic operator expects 1 argument, got %d", len(args))
	}
	f.Index = args[0]
	return nil
}

func (f *indexNotFactor) Table(given int) ([]float64, error) {
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

type indexExclFactor struct {
	Index int
}

func IndexExcl(idx int) Factor {
	return &indexExclFactor{
		Index: idx,
	}
}

func (f *indexExclFactor) SetArgs(args ...int) error {
	if len(args) != 1 {
		return fmt.Errorf("logic operator expects 1 argument, got %d", len(args))
	}
	f.Index = args[0]
	return nil
}

func (f *indexExclFactor) Table(given int) ([]float64, error) {
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
