package ve_test

import (
	"testing"

	"github.com/mlange-42/bbn/ve"
)

func BenchmarkVariablesAddVariable(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	var v1 ve.Variable

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v1 = v.AddVariable(i, ve.ChanceNode, 2)
	}
	b.StopTimer()

	id := v1.Id()
	_ = id
}

func BenchmarkVariablesCreateFactor(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	data := make([]float64, 8)
	var f ve.Factor

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f = v.CreateFactor([]ve.Variable{v1, v2, v3}, data)

	}
	b.StopTimer()

	f2 := v.SumOut(&f, v3)
	_ = f2
}

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

func BenchmarkVariablesRearrange(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	f := v.CreateFactor([]ve.Variable{v1, v2, v3}, make([]float64, 8))

	var f2 ve.Factor
	newVars := []ve.Variable{v3, v2, v1}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f2 = v.Rearrange(&f, newVars)
	}
	b.StopTimer()

	f3 := v.SumOut(&f2, v3)
	_ = f3
}

func BenchmarkVariablesProduct(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	f1 := v.CreateFactor([]ve.Variable{v1, v2}, make([]float64, 4))
	f2 := v.CreateFactor([]ve.Variable{v2, v3}, make([]float64, 4))

	var f3 ve.Factor

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f3 = v.Product(&f1, &f2)
	}
	b.StopTimer()

	f4 := v.SumOut(&f3, v3)
	_ = f4
}

func BenchmarkVariablesSum(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	f1 := v.CreateFactor([]ve.Variable{v1, v2}, make([]float64, 4))
	f2 := v.CreateFactor([]ve.Variable{v2, v3}, make([]float64, 4))

	var f3 ve.Factor

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f3 = v.Sum(&f1, &f2)
	}
	b.StopTimer()

	f4 := v.SumOut(&f3, v3)
	_ = f4
}

func BenchmarkVariablesSumOut(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	f := v.CreateFactor([]ve.Variable{v1, v2, v3}, make([]float64, 8))

	var f2 ve.Factor

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f2 = v.SumOut(&f, v2)
	}
	b.StopTimer()

	f3 := v.SumOut(&f2, v3)
	_ = f3
}

func BenchmarkVariablesMarginal(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	f := v.CreateFactor([]ve.Variable{v1, v2, v3}, make([]float64, 8))

	var f2 ve.Factor

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f2 = v.Marginal(&f, v2)
	}
	b.StopTimer()

	f3 := v.SumOut(&f2, v2)
	_ = f3
}

func BenchmarkVariablesNormalize(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	f := v.CreateFactor([]ve.Variable{v1, v2, v3}, make([]float64, 8))

	var f2 ve.Factor

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f2 = v.Normalize(&f)
	}
	b.StopTimer()

	f3 := v.SumOut(&f2, v2)
	_ = f3
}

func BenchmarkVariablesNormalizeFor(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	f := v.CreateFactor([]ve.Variable{v1, v2, v3}, make([]float64, 8))

	var f2 ve.Factor

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f2 = v.NormalizeFor(&f, v2)
	}
	b.StopTimer()

	f3 := v.SumOut(&f2, v2)
	_ = f3
}

func BenchmarkVariablesInvert(b *testing.B) {
	b.StopTimer()

	v := ve.NewVariables()

	v1 := v.AddVariable(0, ve.ChanceNode, 2)
	v2 := v.AddVariable(1, ve.ChanceNode, 2)
	v3 := v.AddVariable(2, ve.ChanceNode, 2)

	f := v.CreateFactor([]ve.Variable{v1, v2, v3}, make([]float64, 8))

	var f2 ve.Factor

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		f2 = v.Invert(&f)
	}
	b.StopTimer()

	f3 := v.SumOut(&f2, v2)
	_ = f3
}
