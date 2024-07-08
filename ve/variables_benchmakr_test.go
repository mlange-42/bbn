package ve_test

import (
	"testing"

	"github.com/mlange-42/bbn/ve"
)

func BenchmarkVariablesRestrict(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	f := v.CreateFactor([]ve.Variable{v1, v2, v3}, make([]float64, 8))

	var f2 ve.Factor

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f2 = v.Restrict(&f, v2, 1)
	}
	b.StopTimer()

	f3 := v.SumOut(&f2, v3)
	_ = f3
}
