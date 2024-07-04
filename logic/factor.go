package logic

type Factor interface {
	Table(given int) ([]float64, error)
	SetArgs(args ...int) error
}

var Factors = map[string]Factor{
	"not": &Not,

	"and":         &And,
	"not-and":     &NotAnd,
	"and-not":     &AndNot,
	"not-and-not": &NotAndNot,

	"or":         &Or,
	"not-or":     &NotOr,
	"or-not":     &OrNot,
	"not-or-not": &NotOrNot,

	"xor": &XOr,

	"cond":         &Cond,
	"not-cond":     &NotCond,
	"cond-not":     &CondNot,
	"not-cond-not": &NotCondNot,

	"bicond": &BiCond,

	"if-then":         &IfThen,
	"if-not-then":     &IfNotThen,
	"if-then-not":     &IfThenNot,
	"if-not-then-not": &IfNotThenNot,

	"equals":     &Equals,
	"equals-not": &EqualsNot,

	"count-exactly": CountExactly(0),
	"count-less":    CountLess(0),
	"count-greater": CountGreater(0),

	"given":      Given(0),
	"given-not":  GivenNot(0),
	"given-excl": GivenExcl(0),
}
