package ve

import (
	"cmp"
	"fmt"
	"slices"
)

type Evidence struct {
	Variable Variable
	Value    int
}

type Policy struct {
	Decision Variable
	Factor   *Factor
}

type VE struct {
	Variables    *Variables
	eliminated   []bool
	dependencies map[Variable][]Variable
	factors      map[int]*Factor
	weights      []float64
}

func New(variables *Variables, factors []Factor, dependencies map[Variable][]Variable, weights []float64) *VE {
	fac := map[int]*Factor{}
	for _, f := range factors {
		fac[f.id] = &f
	}

	return &VE{
		Variables:    variables,
		eliminated:   make([]bool, len(variables.variables)),
		dependencies: dependencies,
		factors:      fac,
		weights:      weights,
	}
}

func (ve *VE) getDecisions() []Variable {
	dec := []Variable{}
	for _, v := range ve.Variables.variables {
		if v.NodeType == DecisionNode {
			dec = append(dec, v)
		}
	}
	return sortTopological(dec, ve.dependencies)
}

func (ve *VE) eliminateEvidence(evidence []Evidence) {
	for _, ev := range evidence {
		ve.restrictEvidence(ev)
	}
}

func (ve *VE) removeUtilities(except *Variable) {
	utils := []Variable{}
	for _, u := range ve.Variables.variables {
		if u.NodeType == UtilityNode && (except == nil || u.Id != except.Id) {
			utils = append(utils, u)
		}
	}

	if len(utils) == 0 {
		return
	}

	indices := []int{}

	for k, f := range ve.factors {
		for _, u := range utils {
			if slices.Contains(f.Variables, u) {
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
	for _, u := range ve.Variables.variables {
		if u.NodeType != UtilityNode {
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
			if slices.Contains(f.Variables, u) {
				scaled := f
				if ve.weights != nil {
					ff := ve.Variables.Product(scaled, &Factor{Data: []float64{ve.weights[i]}})
					scaled = &ff
				}
				indices = append(indices, k)
				factors = append(factors, scaled)
				break
			}
		}
	}

	sum := ve.Variables.Sum(factors...)
	for _, u := range utils {
		sum = ve.Variables.SumOut(&sum, u)
		ve.eliminated[u.Id] = true
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

	hidden := map[uint16]Variable{}
	for i, v := range ve.Variables.variables {
		if v.NodeType != ChanceNode || ve.eliminated[i] || isDecisionParent[v.Id] {
			continue
		}
		hidden[v.Id] = v
	}
	for _, ev := range evidence {
		delete(hidden, ev.Variable.Id)
	}
	for _, v := range query {
		delete(hidden, v.Id)
	}

	// TODO: check elimination order of hidden variables
	hiddenList := make([]variableDegree, len(hidden))
	i := 0
	newVars := []Variable{}
	for _, v := range hidden {
		for _, f := range ve.factors {
			if !slices.ContainsFunc(f.Variables, func(v2 Variable) bool { return v2.Id == v.Id }) {
				continue
			}
			for _, vv := range f.Variables {
				if !slices.ContainsFunc(newVars, func(v2 Variable) bool { return v2.Id == vv.Id }) {
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
	slices.SortFunc(hiddenList, func(a, b variableDegree) int { return cmp.Compare(a.Variable.Id, b.Variable.Id) })
	slices.SortStableFunc(hiddenList, func(a, b variableDegree) int { return cmp.Compare(a.Degree, b.Degree) })

	for _, v := range hiddenList {
		ve.removeHidden(v.Variable)
	}
}

func (ve *VE) getDecisionParents(single bool) []bool {
	isDecisionParent := make([]bool, len(ve.Variables.variables))

	decisions := ve.getDecisions()
	if len(decisions) == 0 {
		return isDecisionParent
	}

	for i := len(decisions) - 1; i >= 0; i-- {
		dec := decisions[i]
		if vars, ok := ve.dependencies[dec]; ok {
			for _, v := range vars {
				isDecisionParent[v.Id] = true
			}
		}
		if single {
			break
		}
	}

	/*
		for _, v := range ve.Variables.variables {
			if ve.eliminated[v.Id] || v.NodeType != DecisionNode {
				continue
			}
			if vars, ok := ve.dependencies[v]; ok {
				for _, v := range vars {
					isDecisionParent[v.Id] = true
				}
			}
		}*/
	return isDecisionParent
}

func (ve *VE) summarize() *Factor {
	result := ve.multiplyAll()
	resultCopy := *result

	return &resultCopy
}

func (ve *VE) SolveQuery(evidence []Evidence, query []Variable) *Factor {
	return ve.solve(evidence, query, false, nil)
}

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

		deps := ve.dependencies[dec]
		for _, f := range ve.factors {
			if !slices.Contains(f.Variables, dec) {
				continue
			}
			if len(deps) == 0 && len(f.Variables) == 1 {
				factors = append(factors, f)
				continue
			}
			hasParent := true
			//hasNonParent := false
			for _, v := range f.Variables {
				if slices.Contains(deps, v) {
					hasParent = true
				} /*else if v != dec {
					hasNonParent = true
				}*/
			}
			if !hasParent /*|| hasNonParent*/ {
				continue
			}
			factors = append(factors, f)
		}
		if len(factors) == 0 {
			for _, f := range ve.factors {
				if !slices.Contains(f.Variables, dec) {
					continue
				}
				if len(f.Variables) > 1 {
					continue
				}
				factors = append(factors, f)
			}
		}

		/*fmt.Println("Decision on", dec)
		fmt.Println("Remaining factors")
		for _, f := range ve.factors {
			fmt.Println(f)
		}*/

		if len(factors) == 0 {
			panic(fmt.Sprintf("found no factors containing variable %d and its parents", dec.Id))
		}

		// TODO: check that multiplying when multiple factors are remaining is correct!
		var fac *Factor
		if len(factors) == 1 {
			fac = factors[0]
		} else {
			f := ve.Variables.Product(factors...)
			fac = &f
		}
		/*fmt.Println("Selected factors")
		for _, f := range factors {
			fmt.Println(f)
		}
		fmt.Println("Factor product")
		fmt.Println(fac)*/

		policy := ve.Variables.Policy(fac, dec)

		policies[dec] = [2]*Factor{fac, &policy}
		ve.factors[policy.id] = &policy
		ve.Variables.variables[dec.Id].NodeType = ChanceNode

		ve.eliminateHidden(nil, nil, single)

		factors = factors[:0]

		if single {
			break
		}
	}

	return policies
}

func (ve *VE) restrictEvidence(evidence Evidence) {
	indices := []int{}
	for k, v := range ve.factors {
		if slices.Contains(v.Variables, evidence.Variable) {
			indices = append(indices, k)
		}
	}

	for _, idx := range indices {
		fac := ve.factors[idx]
		delete(ve.factors, idx)

		fac2 := ve.Variables.Restrict(fac, evidence.Variable, evidence.Value)

		ve.factors[fac2.id] = &fac2
	}
}

func (ve *VE) removeHidden(variable Variable) {
	indices := []int{}
	factors := []*Factor{}
	//utilityFactors := []*Factor{}

	for k, f := range ve.factors {
		if slices.ContainsFunc(f.Variables, func(v Variable) bool { return v.Id == variable.Id }) {
			indices = append(indices, k)

			/*hasUtility := false
			for _, v := range f.Variables {
				if v.NodeType == UtilityNode {
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
		sum := ve.Variables.Sum(utilityFactors...)
		factors = append(factors, &sum)
	} else if len(utilityFactors) == 1 {
		factors = append(factors, utilityFactors[0])
	}*/

	prod := ve.Variables.Product(factors...)
	prod = ve.Variables.SumOut(&prod, variable)

	for _, idx := range indices {
		delete(ve.factors, idx)
	}

	ve.eliminated[variable.Id] = true
	ve.factors[prod.id] = &prod
}

func (ve *VE) multiplyAll() *Factor {
	factors := []*Factor{}
	//utilityFactors := []*Factor{}
	for _, f := range ve.factors {
		/*hasUtility := false
		for _, v := range f.Variables {
			if v.NodeType == UtilityNode {
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

	f := ve.Variables.Product(factors...)
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
		if id.Id == v.Id {
			for _, p := range vars {
				if p.NodeType != DecisionNode {
					continue
				}
				idx := slices.Index(dec, p)
				result = sortTopologicalRecursive(dec, idx, deps, visited, result)
			}
		}
	}

	return append(result, dec[index])
}
