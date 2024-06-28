package ve

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEliminate(t *testing.T) {
	vars := NewVariables()

	rain := vars.Add(ChanceNode, 2)
	sprinkler := vars.Add(ChanceNode, 2)
	grass := vars.Add(ChanceNode, 2)

	fRain := vars.CreateFactor([]Variable{rain}, []float64{
		0.2, 0.8,
	})

	fSprinkler := vars.CreateFactor([]Variable{rain, sprinkler}, []float64{
		0.01, 0.99, // rain+
		0.2, 0.8, // rain-
	})

	fGrass := vars.CreateFactor([]Variable{rain, sprinkler, grass}, []float64{
		0.99, 0.01, // rain+ sprinkler+
		0.8, 0.2, // rain+ sprinkler-
		0.9, 0.1, // rain- sprinkler+
		0.0, 1.0, // rain- sprinkler-
	})

	ve := New(vars, []Factor{fRain, fSprinkler, fGrass})
	query := []Variable{rain}
	result := ve.Eliminate([]Evidence{{Variable: sprinkler, Value: 1}, {Variable: grass, Value: 0}}, query)

	for _, q := range query {
		fmt.Println(vars.Marginal(result, q))
	}

	assert.Equal(t, []float64{1, 0}, vars.Marginal(result, rain).data)
}
