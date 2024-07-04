package logic

import (
	"fmt"
)

type outcomeIsFactor struct {
	Index  int
	Length int
}

func OutcomeIs(idx, ln int) Factor {
	return &outcomeIsFactor{
		Index:  idx,
		Length: ln,
	}
}

func (f *outcomeIsFactor) SetArgs(args ...int) error {
	if len(args) != 2 {
		return fmt.Errorf("logic operator expects 2 argument (index, length), got %d", len(args))
	}
	f.Index = args[0]
	f.Length = args[1]
	return nil
}

func (f *outcomeIsFactor) Table(given int) ([]float64, error) {
	if given != 1 {
		return nil, fmt.Errorf("logic operator requires 1 operand, but %d were given", given)
	}

	table := make([]float64, f.Length*2)
	for i := 0; i < f.Length; i++ {
		table[i*2+1] = 1
	}
	table[f.Index*2] = 1
	table[f.Index*2+1] = 0

	return table, nil
}

type outcomeIsNotFactor struct {
	Index  int
	Length int
}

func OutcomeIsNot(idx, ln int) Factor {
	return &outcomeIsNotFactor{
		Index:  idx,
		Length: ln,
	}
}

func (f *outcomeIsNotFactor) SetArgs(args ...int) error {
	if len(args) != 2 {
		return fmt.Errorf("logic operator expects 2 argument (index, length), got %d", len(args))
	}
	f.Index = args[0]
	f.Length = args[1]
	return nil
}

func (f *outcomeIsNotFactor) Table(given int) ([]float64, error) {
	if given != 1 {
		return nil, fmt.Errorf("logic operator requires 1 operand, but %d were given", given)
	}

	table := make([]float64, f.Length*2)
	for i := 0; i < f.Length; i++ {
		table[i*2] = 1
	}
	table[f.Index*2] = 0
	table[f.Index*2+1] = 1

	return table, nil
}

type OutcomeEitherFactor struct {
	Indices []int
	Length  int
}

func OutcomeEither(idx []int, ln int) Factor {
	return &OutcomeEitherFactor{
		Indices: idx,
		Length:  ln,
	}
}

func (f *OutcomeEitherFactor) SetArgs(args ...int) error {
	if len(args) < 2 {
		return fmt.Errorf("logic operator expects at least 2 argument (index..., length), got %d", len(args))
	}
	f.Indices = args[:len(args)-1]
	f.Length = args[len(args)-1]
	return nil
}

func (f *OutcomeEitherFactor) Table(given int) ([]float64, error) {
	if given != 1 {
		return nil, fmt.Errorf("logic operator requires 1 operand, but %d were given", given)
	}

	table := make([]float64, f.Length*2)
	for i := 0; i < f.Length; i++ {
		table[i*2+1] = 1
	}
	for _, idx := range f.Indices {
		table[idx*2] = 1
		table[idx*2+1] = 0
	}

	return table, nil
}
