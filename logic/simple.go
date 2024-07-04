package logic

import "fmt"

type boolFactor []bool

func (f *boolFactor) SetArgs(args ...int) {}

func (f *boolFactor) Table(given int) ([]float64, error) {
	arr := *f
	if given != len(arr)/2 {
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

func (f *floatFactor) SetArgs(args ...int) {}

func (f *floatFactor) Table(given int) ([]float64, error) {
	arr := *f
	if given != len(arr)/4 {
		return nil, fmt.Errorf("logic operator requires %d operands, but %d were given", len(arr)/4, given)
	}
	return arr, nil
}

var Not = boolFactor{
	false, // T
	true,  // F
}

var And = boolFactor{
	true,  // T T
	false, // T F
	false, // F T
	false, // F F
}
var NotAnd = boolFactor{
	false, // F T
	false, // F F
	true,  // T T
	false, // T F
}
var AndNot = boolFactor{
	false, // T F
	true,  // T T
	false, // F F
	false, // F T
}
var NotAndNot = boolFactor{
	false, // F F
	false, // F T
	false, // T F
	true,  // T T
}

var Or = boolFactor{
	true,  // T T
	true,  // T F
	true,  // F T
	false, // F F
}
var NotOr = boolFactor{
	true,  // F T
	false, // F F
	true,  // T T
	true,  // T F
}
var OrNot = boolFactor{
	true,  // T F
	true,  // T T
	false, // F F
	true,  // F T
}
var NotOrNot = boolFactor{
	false, // F F
	true,  // F T
	true,  // T F
	true,  // T T
}

var XOr = boolFactor{
	false, // T T
	true,  // T F
	true,  // F T
	false, // F F
}

var Cond = boolFactor{
	true,  // T T
	false, // T F
	true,  // F T
	true,  // F F
}
var NotCond = boolFactor{
	true,  // F T
	true,  // F F
	true,  // T T
	false, // T F
}
var CondNot = boolFactor{
	false, // T F
	true,  // T T
	true,  // F F
	true,  // F T
}
var NotCondNot = boolFactor{
	true,  // F F
	true,  // F T
	false, // T F
	true,  // T T
}

var BiCond = boolFactor{
	true,  // T T
	false, // T F
	false, // F T
	true,  // F F
}

var IfThen = floatFactor{
	1, 0, // T
	0.5, 0.5, // F
}
var IfNotThen = floatFactor{
	0.5, 0.5, // T
	1, 0, // F
}
var IfThenNot = floatFactor{
	0, 1, // T
	0.5, 0.5, // F
}
var IfNotThenNot = floatFactor{
	0.5, 0.5, // T
	0, 1, // F
}

var Equals = boolFactor{
	true,  // T
	false, // F
}
var EqualsNot = boolFactor{
	false, // T
	true,  // F
}
