package logic

// Factor is an interface for logic factors.
type Factor interface {
	// Table returns the CPT for the given number of parent variables.
	Table(given int) ([]float64, error)
	// SetArgs sets factor arguments.
	SetArgs(args ...int) error
}

// Get logic factors by their name.
//
// Primarily used for deserialization.
func Get(name string) (Factor, bool) {
	f, ok := factors[name]
	return f, ok
}

var factors = map[string]Factor{
	"not": Not(),

	"and":         And(),
	"not-and":     NotAnd(),
	"and-not":     AndNot(),
	"not-and-not": NotAndNot(),

	"or":         Or(),
	"not-or":     NotOr(),
	"or-not":     OrNot(),
	"not-or-not": NotOrNot(),

	"xor": XOr(),

	"cond":         Cond(),
	"not-cond":     NotCond(),
	"cond-not":     CondNot(),
	"not-cond-not": NotCondNot(),

	"bicond": BiCond(),

	"if-then":         IfThen(),
	"if-not-then":     IfNotThen(),
	"if-then-not":     IfThenNot(),
	"if-not-then-not": IfNotThenNot(),

	"equals":     Equals(),
	"equals-not": EqualsNot(),

	"count-true":  CountTrue(),
	"count-false": CountFalse(),

	"count-is":      CountIs(0),
	"count-less":    CountLess(0),
	"count-greater": CountGreater(0),

	"given":      Given(0),
	"given-not":  GivenNot(0),
	"given-excl": GivenExcl(0),

	"outcome-is":      OutcomeIs(0, 0),
	"outcome-is-not":  OutcomeIsNot(0, 0),
	"outcome-either":  OutcomeEither(nil, 0),
	"outcome-less":    OutcomeLess(0, 0),
	"outcome-greater": OutcomeGreater(0, 0),

	"bits": Bits(),
}
