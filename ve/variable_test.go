package ve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariable(t *testing.T) {
	v1 := Variable{
		id:       0,
		index:    0,
		outcomes: 2,
		nodeType: ChanceNode,
	}

	v2 := Variable{
		id:       1,
		index:    1,
		outcomes: 2,
		nodeType: ChanceNode,
	}

	assert.Equal(t, 1, v2.Id())
	assert.Equal(t, ChanceNode, v1.NodeType())

	v3 := v1.WithNodeType(UtilityNode)
	assert.Equal(t, UtilityNode, v3.NodeType())

	assert.True(t, v1.Is(v3))
	assert.False(t, v1.Is(v2))
}
