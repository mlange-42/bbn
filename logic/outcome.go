package logic

import (
	"fmt"
)

type outcomeIsFactor struct {
	Index  int
	Length int
}

// OutcomeIs checks whether the parent node has the given index as outcome.
// Argument count is the total number of possible outcomes of the parent.
func OutcomeIs(idx, count int) Factor {
	return &outcomeIsFactor{
		Index:  idx,
		Length: count,
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

// OutcomeIsNot checks whether the parent node does not have the given index as outcome.
// Argument count is the total number of possible outcomes of the parent.
func OutcomeIsNot(idx, count int) Factor {
	return &outcomeIsNotFactor{
		Index:  idx,
		Length: count,
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

type outcomeEitherFactor struct {
	Indices []int
	Length  int
}

// OutcomeEither checks whether the parent node has one of the given indices as outcome.
// Argument count is the total number of possible outcomes of the parent.
func OutcomeEither(idx []int, count int) Factor {
	return &outcomeEitherFactor{
		Indices: idx,
		Length:  count,
	}
}

func (f *outcomeEitherFactor) SetArgs(args ...int) error {
	if len(args) < 2 {
		return fmt.Errorf("logic operator expects at least 2 argument (index..., length), got %d", len(args))
	}
	f.Indices = args[:len(args)-1]
	f.Length = args[len(args)-1]
	return nil
}

func (f *outcomeEitherFactor) Table(given int) ([]float64, error) {
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

type outcomeLessFactor struct {
	Index  int
	Length int
}

// OutcomeLess checks whether the parent node has an index as outcome
// that is less then the given index.
// Argument count is the total number of possible outcomes of the parent.
func OutcomeLess(idx, count int) Factor {
	return &outcomeLessFactor{
		Index:  idx,
		Length: count,
	}
}

func (f *outcomeLessFactor) SetArgs(args ...int) error {
	if len(args) != 2 {
		return fmt.Errorf("logic operator expects 2 argument (index, length), got %d", len(args))
	}
	f.Index = args[0]
	f.Length = args[1]
	return nil
}

func (f *outcomeLessFactor) Table(given int) ([]float64, error) {
	if given != 1 {
		return nil, fmt.Errorf("logic operator requires 1 operand, but %d were given", given)
	}

	table := make([]float64, f.Length*2)
	for i := 0; i < f.Index; i++ {
		table[i*2] = 1
	}
	for i := f.Index; i < f.Length; i++ {
		table[i*2+1] = 1
	}

	return table, nil
}

type outcomeGreaterFactor struct {
	Index  int
	Length int
}

// OutcomeLess checks whether the parent node has an index as outcome
// that is greater then the given index.
// Argument count is the total number of possible outcomes of the parent.
func OutcomeGreater(idx, count int) Factor {
	return &outcomeGreaterFactor{
		Index:  idx,
		Length: count,
	}
}

func (f *outcomeGreaterFactor) SetArgs(args ...int) error {
	if len(args) != 2 {
		return fmt.Errorf("logic operator expects 2 argument (index, length), got %d", len(args))
	}
	f.Index = args[0]
	f.Length = args[1]
	return nil
}

func (f *outcomeGreaterFactor) Table(given int) ([]float64, error) {
	if given != 1 {
		return nil, fmt.Errorf("logic operator requires 1 operand, but %d were given", given)
	}

	table := make([]float64, f.Length*2)
	for i := 0; i < f.Index; i++ {
		table[i*2+1] = 1
	}
	for i := f.Index; i < f.Length; i++ {
		table[i*2] = 1
	}

	return table, nil
}
