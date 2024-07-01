package ve

import (
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
}

func New(variables *Variables, factors []Factor, dependencies map[Variable][]Variable) *VE {
	fac := map[int]*Factor{}
	for _, f := range factors {
		fac[f.id] = &f
	}

	return &VE{
		Variables:    variables,
		eliminated:   make([]bool, len(variables.variables)),
		dependencies: dependencies,
		factors:      fac,
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

func (ve *VE) removeUtilities() {
	utils := []Variable{}
	for _, u := range ve.Variables.variables {
		if u.NodeType == UtilityNode {
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
		if u.NodeType == UtilityNode {
			utils = append(utils, u)
		}
	}

	if len(utils) == 0 {
		return
	}

	indices := []int{}
	factors := []*Factor{}

	for k, f := range ve.factors {
		for _, u := range utils {
			if slices.Contains(f.Variables, u) {
				indices = append(indices, k)
				factors = append(factors, f)
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

func (ve *VE) eliminateHidden(evidence []Evidence, query []Variable, verbose bool) {
	isDecisionParent := ve.getDecisionParents()

	hidden := map[uint16]Variable{}
	for i, v := range ve.Variables.variables {
		if v.NodeType != ChanceNode || ve.eliminated[i] || isDecisionParent[v.Id] {
			continue
		}
		hidden[v.Id] = v
	}
	// TODO: really exclude evidence variables?
	for _, ev := range evidence {
		delete(hidden, ev.Variable.Id)
	}
	for _, v := range query {
		delete(hidden, v.Id)
	}

	if verbose {
		fmt.Println("Hidden variables: ", hidden)
	}
	for _, v := range hidden {
		ve.removeHidden(v)
		if verbose {
			fmt.Println("Eliminate", v)
			ve.printFactors()
		}
	}
}

func (ve *VE) getDecisionParents() []bool {
	isDecisionParent := make([]bool, len(ve.Variables.variables))
	for _, v := range ve.Variables.variables {
		if ve.eliminated[v.Id] || v.NodeType != DecisionNode {
			continue
		}
		if vars, ok := ve.dependencies[v]; ok {
			for _, v := range vars {
				isDecisionParent[v.Id] = true
			}
		}
	}
	return isDecisionParent
}

func (ve *VE) summarize() *Factor {
	result := ve.multiplyAll()
	resultCopy := *result

	return &resultCopy
}

func (ve *VE) SolveQuery(evidence []Evidence, query []Variable, verbose bool) *Factor {
	return ve.solve(evidence, query, false, verbose)
}

func (ve *VE) SolveUtility(evidence []Evidence, query []Variable, verbose bool) *Factor {
	return ve.solve(evidence, query, true, verbose)
}

func (ve *VE) solve(evidence []Evidence, query []Variable, utility bool, verbose bool) *Factor {
	if verbose {
		ve.printFactors()
		fmt.Println("Eliminate evidence")
	}

	ve.eliminateEvidence(evidence)

	if utility {
		if verbose {
			ve.printFactors()
			fmt.Println("Sum utilities")
		}
		ve.sumUtilities()
	} else {
		if verbose {
			ve.printFactors()
			fmt.Println("Remove utilities")
		}
		ve.removeUtilities()
	}

	if verbose {
		ve.printFactors()
		fmt.Println("Eliminate hidden")
	}

	ve.eliminateHidden(evidence, query, verbose)

	if verbose {
		ve.printFactors()
	}

	return ve.summarize()
}

func (ve *VE) SolvePolicies(verbose bool) map[Variable][2]*Factor {
	if verbose {
		fmt.Println("Sum utilities")
	}

	ve.sumUtilities()

	if verbose {
		ve.printFactors()
		fmt.Println("Eliminate hidden")
	}

	ve.eliminateHidden(nil, nil, verbose)

	if verbose {
		ve.printFactors()
		fmt.Println("Policies")
	}

	decisions := ve.getDecisions()
	if len(decisions) == 0 {
		return nil
	}

	if verbose {
		fmt.Println("Collecting decisions")
		fmt.Println(decisions)
	}

	return ve.solvePolicies(decisions, verbose)
}

func (ve *VE) solvePolicies(decisions []Variable, verbose bool) map[Variable][2]*Factor {
	policies := map[Variable][2]*Factor{}
	for i := len(decisions) - 1; i >= 0; i-- {
		dec := decisions[i]
		if verbose {
			fmt.Println("Solving decision", dec)
		}

		deps := ve.dependencies[dec]
		factorIdx := -1
		for i, f := range ve.factors {
			if !slices.Contains(f.Variables, dec) {
				continue
			}
			if len(deps) == 0 {
				if factorIdx >= 0 {
					panic(fmt.Sprintf("found multiple factors containing variable %d and its parents", dec.Id))
				}
				factorIdx = i
				continue
			}
			for _, v := range deps {
				if slices.Contains(f.Variables, v) {
					if factorIdx >= 0 {
						panic(fmt.Sprintf("found multiple factors containing variable %d and its parents", dec.Id))
					}
					factorIdx = i
					break
				}
			}
		}

		if factorIdx < 0 {
			panic(fmt.Sprintf("found no factors containing variable %d and its parents", dec.Id))
		}

		policy := ve.Variables.Policy(ve.factors[factorIdx], dec)
		if verbose {
			fmt.Println("Utility")
			fmt.Println(ve.factors[factorIdx])
			fmt.Println("Policy")
			fmt.Println(policy)
		}

		policies[dec] = [2]*Factor{ve.factors[factorIdx], &policy}
		ve.factors[policy.id] = &policy
		ve.Variables.variables[dec.Id].NodeType = ChanceNode

		if verbose {
			fmt.Println("Added policy", dec)
			ve.printFactors()
			fmt.Println("Eliminate hidden", dec)
		}

		ve.eliminateHidden(nil, nil, verbose)

		if verbose {
			ve.printFactors()
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

		// TODO: do we need these two steps?
		/*if len(fac2.Variables) == 0 {
			continue
		}
		if len(fac2.Variables) == 1 {
			fac2 = ve.Variables.Normalize(&fac2)
		}*/
		// end TODO

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
	for _, f := range ve.factors {
		factors = append(factors, f)
	}

	clear(ve.factors)

	f := ve.Variables.Product(factors...)
	ve.factors[f.id] = &f

	return &f
}

func (ve *VE) printFactors() {
	for k, v := range ve.factors {
		fmt.Printf("%d %v\n", k, v)
	}
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
