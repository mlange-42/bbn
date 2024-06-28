package ve

import (
	"fmt"
	"slices"
)

type Evidence struct {
	Variable Variable
	Value    int
}

type Dependencies struct {
	Decision Variable
	Parents  []Variable
}

type Policy struct {
	Decision Variable
	Factor   *Factor
}

type VE struct {
	variables          *Variables
	eliminated         []bool
	decisions          []Variable
	unhandledDecisions []Variable
	//policies           []Policy
	dependencies []Dependencies
	factors      map[int]*Factor
	evidence     []Evidence
	query        []Variable
}

func New(variables *Variables, factors []Factor, dependencies []Dependencies, evidence []Evidence, query []Variable) VE {
	fac := map[int]*Factor{}
	for _, f := range factors {
		fac[f.id] = &f
	}

	dec := []Variable{}
	for _, v := range variables.variables {
		if v.nodeType == DecisionNode {
			dec = append(dec, v)
		}
	}

	dec = sortTopological(dec, dependencies)

	return VE{
		variables:          variables,
		decisions:          dec,
		unhandledDecisions: append([]Variable{}, dec...),
		eliminated:         make([]bool, len(variables.variables)),
		dependencies:       dependencies,
		factors:            fac,
		evidence:           evidence,
		query:              query,
	}
}

func (ve *VE) eliminateEvidence() {
	for _, ev := range ve.evidence {
		ve.restrictEvidence(ev)
	}
}

func (ve *VE) eliminateDecisions() {
	ve.eliminateHidden()
}

func (ve *VE) eliminateHidden() {
	isDecisionParent := make([]bool, len(ve.variables.variables))
	for _, dep := range ve.dependencies {
		if ve.eliminated[dep.Decision.id] {
			continue
		}
		for _, v := range dep.Parents {
			isDecisionParent[v.id] = true
		}
	}

	hidden := map[uint16]Variable{}
	for _, v := range ve.variables.variables {
		if v.nodeType != ChanceNode || isDecisionParent[v.id] {
			continue
		}
		hidden[v.id] = v
	}
	// TODO: really exclude evidence variables?
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
	ve.eliminateDecisions()
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
	fmt.Println("Product:")
	fmt.Println(prod)

	prod = ve.variables.SumOut(&prod, variable)
	fmt.Println("SumOut:")
	fmt.Println(prod)

	for _, idx := range indices {
		delete(ve.factors, idx)
	}

	ve.eliminated[variable.id] = true
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

func sortTopological(dec []Variable, deps []Dependencies) []Variable {
	result := []Variable{}
	visited := make([]bool, len(dec))

	for i := range dec {
		result = sortTopologicalRecursive(dec, i, deps, visited, result)
	}

	return result
}

func sortTopologicalRecursive(dec []Variable, index int, deps []Dependencies, visited []bool, result []Variable) []Variable {
	if visited[index] {
		return result
	}

	visited[index] = true

	v := dec[index]
	for _, dep := range deps {
		if dep.Decision.id == v.id {
			for _, p := range dep.Parents {
				if p.nodeType != DecisionNode {
					continue
				}
				idx := slices.Index(dec, p)
				result = sortTopologicalRecursive(dec, idx, deps, visited, result)
			}
		}
	}

	return append(result, dec[index])
}
