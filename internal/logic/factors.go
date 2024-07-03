package logic

type Factor []bool

var Not = Factor{
	false, // T
	true,  // F
}

var And = Factor{
	true,  // T T
	false, // T F
	false, // F T
	false, // F F
}

var Or = Factor{
	true,  // T T
	true,  // T F
	true,  // F T
	false, // F F
}

var XOr = Factor{
	false, // T T
	true,  // T F
	true,  // F T
	false, // F F
}

var Cond = Factor{
	true,  // T T
	false, // T F
	true,  // F T
	true,  // F F
}

var BiCond = Factor{
	true,  // T T
	false, // T F
	false, // F T
	true,  // F F
}

var Factors = map[string]Factor{
	"not":    Not,
	"and":    And,
	"or":     Or,
	"xor":    XOr,
	"cond":   Cond,
	"bicond": BiCond,
}

func (f Factor) Operands() int {
	return len(f) / 2
}

func (f Factor) Table() []float64 {
	table := make([]float64, len(f)*2)
	for i, v := range f {
		if v {
			table[i*2] = 1
		} else {
			table[i*2+1] = 1
		}
	}
	return table
}
