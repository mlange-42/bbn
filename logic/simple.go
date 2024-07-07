package logic

import "fmt"

type boolFactor []bool

func (f *boolFactor) SetArgs(args ...int) error {
	if len(args) > 0 {
		return fmt.Errorf("logic operator expects zero arguments, got %d", len(args))
	}
	return nil
}

func (f *boolFactor) Table(given int) ([]float64, error) {
	arr := *f
	if len(arr) != 1<<given { // 2^given
		return nil, fmt.Errorf("logic operator requires %d operands, but %d were given", len(arr)/2, given)
	}

	table := make([]float64, len(arr)*2)
	for i, v := range arr {
		if v {
			table[i*2] = 1
		} else {
			table[i*2+1] = 1
		}
	}
	return table, nil
}

type floatFactor []float64

func (f *floatFactor) SetArgs(args ...int) error {
	if len(args) > 0 {
		return fmt.Errorf("logic operator expects zero arguments, got %d", len(args))
	}
	return nil
}

func (f *floatFactor) Table(given int) ([]float64, error) {
	arr := *f
	if len(arr) != 1<<(given+1) { // 2^(given+1)
		return nil, fmt.Errorf("logic operator requires %d operands, but %d were given", len(arr)/4, given)
	}
	return arr, nil
}

var not = boolFactor{
	false, // T
	true,  // F
}

func Not() Factor {
	return &not
}

var and = boolFactor{
	true,  // T T
	false, // T F
	false, // F T
	false, // F F
}

func And() Factor {
	return &and
}

var notAnd = boolFactor{
	false, // F T
	false, // F F
	true,  // T T
	false, // T F
}

func NotAnd() Factor {
	return &notAnd
}

var andNot = boolFactor{
	false, // T F
	true,  // T T
	false, // F F
	false, // F T
}

func AndNot() Factor {
	return &andNot
}

var notAndNot = boolFactor{
	false, // F F
	false, // F T
	false, // T F
	true,  // T T
}

func NotAndNot() Factor {
	return &notAndNot
}

var or = boolFactor{
	true,  // T T
	true,  // T F
	true,  // F T
	false, // F F
}

func Or() Factor {
	return &or
}

var notOr = boolFactor{
	true,  // F T
	false, // F F
	true,  // T T
	true,  // T F
}

func NotOr() Factor {
	return &notOr
}

var orNot = boolFactor{
	true,  // T F
	true,  // T T
	false, // F F
	true,  // F T
}

func OrNot() Factor {
	return &orNot
}

var notOrNot = boolFactor{
	false, // F F
	true,  // F T
	true,  // T F
	true,  // T T
}

func NotOrNot() Factor {
	return &notOrNot
}

var xOr = boolFactor{
	false, // T T
	true,  // T F
	true,  // F T
	false, // F F
}

func XOr() Factor {
	return &xOr
}

var cond = boolFactor{
	true,  // T T
	false, // T F
	true,  // F T
	true,  // F F
}

func Cond() Factor {
	return &cond
}

var notCond = boolFactor{
	true,  // F T
	true,  // F F
	true,  // T T
	false, // T F
}

func NotCond() Factor {
	return &notCond
}

var condNot = boolFactor{
	false, // T F
	true,  // T T
	true,  // F F
	true,  // F T
}

func CondNot() Factor {
	return &condNot
}

var notCondNot = boolFactor{
	true,  // F F
	true,  // F T
	false, // T F
	true,  // T T
}

func NotCondNot() Factor {
	return &notCondNot
}

var biCond = boolFactor{
	true,  // T T
	false, // T F
	false, // F T
	true,  // F F
}

func BiCond() Factor {
	return &biCond
}

var ifThen = floatFactor{
	1, 0, // T
	0.5, 0.5, // F
}

func IfThen() Factor {
	return &ifThen
}

var ifNotThen = floatFactor{
	0.5, 0.5, // T
	1, 0, // F
}

func IfNotThen() Factor {
	return &ifNotThen
}

var ifThenNot = floatFactor{
	0, 1, // T
	0.5, 0.5, // F
}

func IfThenNot() Factor {
	return &ifThenNot
}

var ifNotThenNot = floatFactor{
	0.5, 0.5, // T
	0, 1, // F
}

func IfNotThenNot() Factor {
	return &ifNotThenNot
}

var equals = boolFactor{
	true,  // T
	false, // F
}

func Equals() Factor {
	return &equals
}

var equalsNot = boolFactor{
	false, // T
	true,  // F
}

func EqualsNot() Factor {
	return &equalsNot
}
