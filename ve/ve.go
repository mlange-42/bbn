package ve

import (
	"cmp"
	"fmt"
	"slices"
)

// Evidence for one variable.
type Evidence struct {
	Variable Variable
	Value    int
}

// VE performs the variable elimination algorithm.
type VE struct {
	variables    *Variables
	eliminated   []bool
	dependencies map[Variable][]Variable
	factors      map[int]*Factor
	weights      []float64
}

// New creates a [VE] instance from the given variables, factors
// decision dependencies and utility weights.
//
// Weights should nil, be  as long as the number of utility variables.
// Order of weights is the same the order of utility variable in all variables.
func New(variables *Variables, factors []Factor, dependencies map[Variable][]Variable, weights []float64) *VE {
	fac := map[int]*Factor{}
	for _, f := range factors {
		fac[f.id] = &f
	}

	return &VE{
		variables:    variables,
		eliminated:   make([]bool, len(variables.variables)),
		dependencies: dependencies,
		factors:      fac,
		weights:      weights,
	}
}

// Variables for the VE.
func (ve *VE) Variables() *Variables {
	return ve.variables
}

// getDecisions finds all decision variables and returns them in topological order.
func (ve *VE) getDecisions() []Variable {
	dec := []Variable{}
	for _, v := range ve.variables.variables {
		if v.NodeType() == DecisionNode {
			dec = append(dec, v)
		}
	}
	return sortTopological(dec, ve.dependencies)
}

// Eliminate all evidence from all factors
func (ve *VE) eliminateEvidence(evidence []Evidence) {
	for _, ev := range evidence {
		ve.restrictEvidence(ev)
	}
}

func (ve *VE) removeUtilities(except *Variable) {
	utils := []Variable{}
	for _, u := range ve.variables.variables {
		if u.NodeType() == UtilityNode && (except == nil || u.id != except.id) {
			utils = append(utils, u)
		}
	}

	if len(utils) == 0 {
		return
	}

	indices := []int{}

	for k, f := range ve.factors {
		for _, u := range utils {
			if slices.Contains(f.variables, u) {
				indices = append(indices, k)
				break
			}
		}
	}

	for _, idx := range indices {
		delete(ve.factors, idx)
	}
}

func (ve *VE) sumUtilities() {
	utils := []Variable{}
	for _, u := range ve.variables.variables {
		if u.NodeType() != UtilityNode {
			continue
		}
		utils = append(utils, u)
	}

	if len(utils) == 0 {
		return
	}

	indices := []int{}
	factors := []*Factor{}

	for k, f := range ve.factors {
		for i, u := range utils {
			if slices.Contains(f.variables, u) {
				scaled := f
				if ve.weights != nil {
					ff := ve.variables.Product(scaled, &Factor{data: []float64{ve.weights[i]}})
					scaled = &ff
				}
				indices = append(indices, k)
				factors = append(factors, scaled)
				break
			}
		}
	}

	sum := ve.variables.Sum(factors...)
	for _, u := range utils {
		sum = ve.variables.SumOut(&sum, u)
		ve.eliminated[u.index] = true
	}

	for _, idx := range indices {
		delete(ve.factors, idx)
	}

	ve.factors[sum.id] = &sum
}

type variableDegree struct {
	Variable Variable
	Degree   int
}

func (ve *VE) eliminateHidden(evidence []Evidence, query []Variable, singleDecision bool) {
	isDecisionParent := ve.getDecisionParents(singleDecision)

	hidden := map[int]Variable{}
	for i, v := range ve.variables.variables {
		if v.NodeType() != ChanceNode || ve.eliminated[i] || isDecisionParent[v.index] {
			continue
		}
		hidden[v.id] = v
	}
	for _, ev := range evidence {
		delete(hidden, ev.Variable.id)
	}
	for _, v := range query {
		delete(hidden, v.id)
	}

	// TODO: check elimination order of hidden variables
	hiddenList := make([]variableDegree, len(hidden))
	i := 0
	newVars := make([]Variable, 0, len(ve.factors))
	for _, v := range hidden {
		for _, f := range ve.factors {
			if !slices.ContainsFunc(f.variables, func(v2 Variable) bool { return v2.id == v.id }) {
				continue
			}
			for _, vv := range f.variables {
				if !slices.ContainsFunc(newVars, func(v2 Variable) bool { return v2.id == vv.id }) {
					newVars = append(newVars, vv)
				}
			}
		}
		degree := len(newVars)
		newVars = newVars[:0]
		hiddenList[i] = variableDegree{
			Variable: v,
			Degree:   degree,
		}
		i++
	}
	slices.SortFunc(hiddenList, func(a, b variableDegree) int { return cmp.Compare(a.Variable.id, b.Variable.id) })
	slices.SortStableFunc(hiddenList, func(a, b variableDegree) int { return cmp.Compare(a.Degree, b.Degree) })

	for _, v := range hiddenList {
		ve.removeHidden(v.Variable)
	}
}

func (ve *VE) getDecisionParents(single bool) []bool {
	isDecisionParent := make([]bool, len(ve.variables.variables))

	decisions := ve.getDecisions()
	if len(decisions) == 0 {
		return isDecisionParent
	}

	for i := len(decisions) - 1; i >= 0; i-- {
		dec := decisions[i]
		if vars, ok := ve.dependencies[dec]; ok {
			for _, v := range vars {
				isDecisionParent[v.index] = true
			}
		}
		if single {
			break
		}
	}

	/*
		for _, v := range ve.variables.variables {
			if ve.eliminated[v.id] || v.NodeType() != DecisionNode {
				continue
			}
			if vars, ok := ve.dependencies[v]; ok {
				for _, v := range vars {
					isDecisionParent[v.id] = true
				}
			}
		}*/
	return isDecisionParent
}

func (ve *VE) summarize() *Factor {
	return ve.multiplyAll()
}

// SolveQuery solves marginal probabilities for the given query variables, and the given evidence.
func (ve *VE) SolveQuery(evidence []Evidence, query []Variable) *Factor {
	return ve.solve(evidence, query, false, nil)
}

// SolveQuery solves utilities for the given query variables, and the given evidence.
//
// Argument utilityVar can be used to solve for only this variable, dropping all other utilities.
// Solves total utility if utilityVar is nil.
func (ve *VE) SolveUtility(evidence []Evidence, query []Variable, utilityVar *Variable) *Factor {
	return ve.solve(evidence, query, true, utilityVar)
}

func (ve *VE) solve(evidence []Evidence, query []Variable, utility bool, utilityVar *Variable) *Factor {
	ve.eliminateEvidence(evidence)

	if utility {
		if utilityVar == nil {
			ve.sumUtilities()
		} else {
			ve.removeUtilities(utilityVar)
		}
	} else {
		ve.removeUtilities(nil)
	}

	ve.eliminateHidden(evidence, query, false)

	return ve.summarize()
}

// SolvePolicies solves decision policies.
//
// Solves only the last decision if single is true.
func (ve *VE) SolvePolicies(single bool) map[Variable][2]*Factor {
	decisions := ve.getDecisions()
	if len(decisions) == 0 {
		return nil
	}

	ve.sumUtilities()

	ve.eliminateHidden(nil, nil, single)

	return ve.solvePolicies(decisions, single)
}

func (ve *VE) solvePolicies(decisions []Variable, single bool) map[Variable][2]*Factor {
	policies := map[Variable][2]*Factor{}
	factors := []*Factor{}
	for i := len(decisions) - 1; i >= 0; i-- {
		dec := decisions[i]
		factors = ve.findDecisionFactors(dec, factors)

		/*fmt.Println("Decision on", dec)
		fmt.Println("Remaining factors")
		for _, f := range ve.factors {
			fmt.Println(f)
		}*/

		if len(factors) == 0 {
			panic(fmt.Sprintf("found no factors containing variable %d and its parents", dec.id))
		}

		// TODO: check that multiplying when multiple factors are remaining is correct!
		var fac *Factor
		if len(factors) == 1 {
			fac = factors[0]
		} else {
			f := ve.variables.Product(factors...)
			fac = &f
		}
		/*fmt.Println("Selected factors")
		for _, f := range factors {
			fmt.Println(f)
		}
		fmt.Println("Factor product")
		fmt.Println(fac)*/

		policy := ve.variables.Policy(fac, dec)

		policies[dec] = [2]*Factor{fac, &policy}
		ve.factors[policy.id] = &policy
		ve.variables.variables[dec.index].nodeType = ChanceNode

		ve.eliminateHidden(nil, nil, single)

		factors = factors[:0]

		if single {
			break
		}
	}

	return policies
}

func (ve *VE) findDecisionFactors(decision Variable, result []*Factor) []*Factor {
	deps := ve.dependencies[decision]
	for _, f := range ve.factors {
		if !slices.Contains(f.variables, decision) {
			continue
		}
		if len(deps) == 0 && len(f.variables) == 1 {
			result = append(result, f)
			continue
		}
		hasParent := true
		//hasNonParent := false
		for _, v := range f.variables {
			if slices.Contains(deps, v) {
				hasParent = true
			} /*else if v != dec {
				hasNonParent = true
			}*/
		}
		if !hasParent /*|| hasNonParent*/ {
			continue
		}
		result = append(result, f)
	}
	if len(result) == 0 {
		for _, f := range ve.factors {
			if !slices.Contains(f.variables, decision) {
				continue
			}
			if len(f.variables) > 1 {
				continue
			}
			result = append(result, f)
		}
	}

	return result
}

func (ve *VE) restrictEvidence(evidence Evidence) {
	indices := make([]int, 0, len(ve.factors))
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
	indices := make([]int, 0, len(ve.factors))
	factors := make([]*Factor, 0, len(ve.factors))
	//utilityFactors := []*Factor{}

	for k, f := range ve.factors {
		if slices.ContainsFunc(f.variables, func(v Variable) bool { return v.id == variable.id }) {
			indices = append(indices, k)

			/*hasUtility := false
			for _, v := range f.Variables {
				if v.NodeType() == UtilityNode {
					hasUtility = true
					break
				}
			}
			if hasUtility {
				utilityFactors = append(utilityFactors, f)
			} else {*/
			factors = append(factors, f)
			//}
		}
	}

	/*if len(utilityFactors) > 1 {
		sum := ve.variables.Sum(utilityFactors...)
		factors = append(factors, &sum)
	} else if len(utilityFactors) == 1 {
		factors = append(factors, utilityFactors[0])
	}*/

	prod := ve.variables.Product(factors...)
	prod = ve.variables.SumOut(&prod, variable)

	for _, idx := range indices {
		delete(ve.factors, idx)
	}

	ve.eliminated[variable.index] = true
	ve.factors[prod.id] = &prod
}

func (ve *VE) multiplyAll() *Factor {
	factors := make([]*Factor, 0, len(ve.factors))
	//utilityFactors := []*Factor{}
	for _, f := range ve.factors {
		/*hasUtility := false
		for _, v := range f.Variables {
			if v.NodeType() == UtilityNode {
				hasUtility = true
				break
			}
		}
		if hasUtility {
			utilityFactors = append(utilityFactors, f)
		} else {*/
		factors = append(factors, f)
		//}
	}

	clear(ve.factors)

	/*if len(utilityFactors) > 1 {
		sum := ve.Variables.Sum(utilityFactors...)
		factors = append(factors, &sum)
	} else if len(utilityFactors) == 1 {
		factors = append(factors, utilityFactors[0])
	}*/

	f := ve.variables.Product(factors...)
	ve.factors[f.id] = &f

	return &f
}

func sortTopological(dec []Variable, deps map[Variable][]Variable) []Variable {
	result := []Variable{}
	visited := make([]bool, len(dec))

	for i := range dec {
		result = sortTopologicalRecursive(dec, i, deps, visited, result)
	}

	return result
}

func sortTopologicalRecursive(dec []Variable, index int, deps map[Variable][]Variable, visited []bool, result []Variable) []Variable {
	if visited[index] {
		return result
	}

	visited[index] = true

	v := dec[index]
	for id, vars := range deps {
		if id.id == v.id {
			for _, p := range vars {
				if p.NodeType() != DecisionNode {
					continue
				}
				idx := slices.Index(dec, p)
				result = sortTopologicalRecursive(dec, idx, deps, visited, result)
			}
		}
	}

	return append(result, dec[index])
}
