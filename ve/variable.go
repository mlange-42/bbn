package ve

// NodeType ID.
type NodeType uint8

const (
	ChanceNode NodeType = iota
	DecisionNode
	UtilityNode
)

// Variable definition for variable elimination by [VE].
type Variable struct {
	id       int
	index    uint16
	outcomes uint16
	nodeType NodeType
}

// Id of the variable.
func (v *Variable) Id() int {
	return v.id
}

// NodeType of the variable.
func (v *Variable) NodeType() NodeType {
	return v.nodeType
}

// WithNodeType returns a new variable with the same Id as the original variable,
// but with a different node type.
func (v Variable) WithNodeType(tp NodeType) Variable {
	v.nodeType = tp
	return v
}

// Is returns whether this variable is the same one as other.
//
// Compared variable Id.
func (v Variable) Is(other Variable) bool {
	return v.id == other.id
}
