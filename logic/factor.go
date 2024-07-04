package logic

type Factor interface {
	Table(given int) ([]float64, error)
}

var Factors = map[string]Factor{
	"not": Not,

	"and":         And,
	"not-and":     NotAnd,
	"and-not":     AndNot,
	"not-and-not": NotAndNot,

	"or":         Or,
	"not-or":     NotOr,
	"or-not":     OrNot,
	"not-or-not": NotOrNot,

	"xor": XOr,

	"cond":         Cond,
	"not-cond":     NotCond,
	"cond-not":     CondNot,
	"not-cond-not": NotCondNot,

	"bicond": BiCond,

	"if-then":         IfThen,
	"if-not-then":     IfNotThen,
	"if-then-not":     IfThenNot,
	"if-not-then-not": IfNotThenNot,

	"equals":     Equals,
	"equals-not": EqualsNot,
}
