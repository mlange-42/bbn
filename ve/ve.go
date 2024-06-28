package ve

import (
	"slices"
)

type Evidence struct {
	Variable Variable
	Value    int
}

type VE struct {
	variables *Variables
	factors   map[int]*Factor
	evidence  []Evidence
	query     []Variable
}

func New(variables *Variables, factors []Factor, evidence []Evidence, query []Variable) VE {
	fac := map[int]*Factor{}
	for _, f := range factors {
		fac[f.id] = &f
	}
	return VE{
		variables: variables,
		factors:   fac,
		evidence:  evidence,
		query:     query,
	}
}

func (ve *VE) eliminateEvidence() {
	for _, ev := range ve.evidence {
		ve.restrictEvidence(ev)
	}
}
func (ve *VE) eliminateHidden() {
	hidden := map[uint16]Variable{}
	for _, v := range ve.variables.variables {
		hidden[v.id] = v
	}
	for _, ev := range ve.evidence {
		delete(hidden, ev.Variable.id)
	}
	for _, v := range ve.query {
		delete(hidden, v.id)
	}

	for _, v := range hidden {
		ve.removeHidden(v)
	}
}

func (ve *VE) summarize() *Factor {
	result := ve.multiplyAll()
	resultCopy := *result

	return &resultCopy
}

func (ve *VE) Eliminate() *Factor {
	ve.eliminateEvidence()
	ve.eliminateHidden()
	return ve.summarize()
}

func (ve *VE) restrictEvidence(evidence Evidence) {
	indices := []int{}
	for k, v := range ve.factors {
		if slices.Contains(v.variables, evidence.Variable) {
			indices = append(indices, k)
		}
	}

	for _, idx := range indices {
		fac := ve.factors[idx]
		delete(ve.factors, idx)

		fac2 := ve.variables.Restrict(fac, evidence.Variable, evidence.Value)
		ve.factors[fac2.id] = &fac2
	}
}

func (ve *VE) removeHidden(variable Variable) {
	indices := []int{}
	factors := []*Factor{}

	for k, f := range ve.factors {
		if slices.Contains(f.variables, variable) {
			indices = append(indices, k)
			factors = append(factors, f)
		}
	}

	prod := ve.variables.Product(factors...)
	prod = ve.variables.SumOut(&prod, variable)

	for _, idx := range indices {
		delete(ve.factors, idx)
	}

	ve.factors[prod.id] = &prod
}

func (ve *VE) multiplyAll() *Factor {
	factors := []*Factor{}
	for _, f := range ve.factors {
		factors = append(factors, f)
	}

	clear(ve.factors)

	f := ve.variables.Product(factors...)
	ve.factors[f.id] = &f

	return &f
}
