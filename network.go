package bbn

import (
	"fmt"
	"slices"

	"github.com/mlange-42/bbn/internal/ve"
)

type Variable struct {
	Name     string
	Type     ve.NodeType
	Outcomes []string
	Position [2]int
	Color    string
	Factor   *Factor // Don't set this, it is initialized when constructing the network.
}

type Factor struct {
	For      string
	Given    []string  `yaml:",omitempty"`
	Table    []float64 `yaml:",omitempty"`
	outcomes []int
	columns  int
}

func (f *Factor) Row(indices []int) ([]float64, bool) {
	idx, ok := f.RowIndex(indices)
	if !ok {
		return nil, false
	}
	return f.Table[idx : idx+f.columns], true
}

func (f *Factor) RowIndex(indices []int) (int, bool) {
	if len(indices) != len(f.outcomes) {
		panic(fmt.Sprintf("factor with %d given variables can't use %d indices", len(f.outcomes), len(indices)))
	}

	if len(indices) == 0 {
		return 0, true
	}

	curr := len(f.outcomes) - 1
	idx := indices[curr]
	if idx < 0 {
		return 0, false
	}
	stride := 1
	curr--
	for curr >= 0 {
		currIdx := indices[curr]
		if currIdx < 0 {
			return 0, false
		}
		stride *= int(f.outcomes[curr+1])
		idx += currIdx * stride
		curr--
	}
	return idx, true
}

type variable struct {
	Variable   Variable
	VeVariable ve.Variable
}

type Network struct {
	name              string
	info              string
	variables         []Variable
	factors           []Factor
	policies          map[string]ve.Factor
	ve                *ve.VE
	variableNames     map[string]*variable
	totalUtilityIndex int
}

func New(name string, info string, variables []Variable, factors []Factor) (*Network, error) {
	net := &Network{
		name:      name,
		info:      info,
		variables: variables,
		factors:   factors,
		policies:  map[string]ve.Factor{},
	}
	err := net.prepareVariables()
	if err != nil {
		return nil, err
	}
	return net, nil
}

func (n *Network) prepareVariables() error {
	varNames := map[string]*Variable{}
	outcomes := make(map[string]int, len(n.variables))
	for i := range n.variables {
		v := &n.variables[i]
		if _, ok := varNames[v.Name]; ok {
			return fmt.Errorf("duplicate variable name %s", v.Name)
		}
		varNames[v.Name] = v
		outcomes[v.Name] = len(v.Outcomes)
	}

	for i := range n.variables {
		v := &n.variables[i]
		idx := slices.IndexFunc(n.factors, func(f Factor) bool { return f.For == v.Name })
		if idx < 0 {
			continue
		}
		v.Factor = &n.factors[idx]

		v.Factor.columns = len(v.Outcomes)
		v.Factor.outcomes = make([]int, len(v.Factor.Given))
		for i, g := range v.Factor.Given {
			n, ok := outcomes[g]
			if !ok {
				return fmt.Errorf("parent variable %s of %s not found", g, v.Name)
			}
			v.Factor.outcomes[i] = n
		}
	}

	return n.prepareUtilityNodes(varNames)
}

func (n *Network) prepareUtilityNodes(varNames map[string]*Variable) error {
	n.totalUtilityIndex = -1
	for i := range n.variables {
		v := &n.variables[i]
		if v.Type != ve.UtilityNode {
			continue
		}
		hasUtilParents := false
		hasOtherParents := false
		for _, parent := range v.Factor.Given {
			p, ok := varNames[parent]
			if !ok {
				return fmt.Errorf("parent node %s for %s not found", parent, v.Name)
			}
			if p.Type == ve.UtilityNode {
				hasUtilParents = true
			} else {
				hasOtherParents = true
			}
		}
		if hasUtilParents && hasOtherParents {
			return fmt.Errorf("utility node %s has utility parents and other parents; can only have either or", v.Name)
		}
		if !hasUtilParents && !hasOtherParents {
			return fmt.Errorf("utility node %s has no parents", v.Name)
		}
		if !hasUtilParents {
			continue
		}
		if n.totalUtilityIndex >= 0 {
			return fmt.Errorf("found multiple nodes for total utility")
		}
		if len(v.Outcomes) != len(v.Factor.Given) {
			return fmt.Errorf("invalid total utility node; number of parents and number of outcomes must be the same")
		}
		n.totalUtilityIndex = i
	}
	return nil
}

func (n *Network) Name() string {
	return n.name
}

func (n *Network) Info() string {
	return n.info
}

func (n *Network) Variables() []Variable {
	return n.variables
}

func (n *Network) TotalUtilityIndex() int {
	return n.totalUtilityIndex
}

// SolvePolicies solves and inserts policies for decisions, using Variable Elimination.
func (n *Network) SolvePolicies(stepwise bool) (map[string]Factor, error) {
	clear(n.policies)

	decisions := n.countDecisionSteps(stepwise)
	for i := 0; i < decisions; i++ {
		var err error
		n.ve, n.variableNames, err = n.toVE(nil)
		if err != nil {
			return nil, err
		}
		policies := n.ve.SolvePolicies(stepwise)
		if policies == nil {
			break
		}

		for name, v := range n.variableNames {
			if v.VeVariable.NodeType != ve.DecisionNode {
				continue
			}
			if p, ok := policies[v.VeVariable]; ok {
				n.policies[name] = *p[1]
			}
		}
	}

	result := map[string]Factor{}
	for name, f := range n.policies {
		variables := f.Variables()
		forVar := n.variableNames[name]
		newVars := make([]ve.Variable, len(variables))
		idx := slices.IndexFunc(
			variables,
			func(v ve.Variable) bool { return v.Id == forVar.VeVariable.Id },
		)

		for i := 0; i < idx; i++ {
			newVars[i] = variables[i]
		}
		for i := idx + 1; i < len(variables); i++ {
			newVars[i-1] = variables[i]
		}
		newVars[len(newVars)-1] = variables[idx]

		f := n.ve.Variables.Rearrange(&f, newVars)

		given := make([]string, len(newVars)-1)
		for i := 0; i < len(newVars)-1; i++ {
			given[i] = n.variables[newVars[i].Id].Name
		}

		ff := Factor{
			For:   name,
			Given: given,
			Table: f.Data(),
		}
		result[name] = ff
	}

	return result, nil
}

func (n *Network) countDecisionSteps(stepwise bool) int {
	if !stepwise {
		return 1
	}
	decisions := 0
	for i := range n.variables {
		if n.variables[i].Type == ve.DecisionNode {
			decisions++
		}
	}
	if decisions == 0 {
		return 1
	}
	return decisions
}

// SolveQuery solves a query, using Variable Elimination.
func (n *Network) SolveQuery(evidence map[string]string, query []string, ignorePolicies bool) (map[string][]float64, *ve.Factor, error) {
	f, err := n.solve(evidence, query, false, "", ignorePolicies)
	if err != nil {
		return nil, nil, err
	}

	result := map[string][]float64{}
	for _, q := range query {
		m := n.Marginal(f, q)
		n := n.Normalize(&m)
		result[q] = n.Data()
	}

	return result, f, nil
}

// SolveUtility solves utility, using Variable Elimination.
func (n *Network) SolveUtility(evidence map[string]string, query []string, utilityVar string, ignorePolicies bool) (*ve.Factor, error) {
	return n.solve(evidence, query, true, utilityVar, ignorePolicies)
}

// solve solves a query or utility, using Variable Elimination.
func (n *Network) solve(evidence map[string]string, query []string, utility bool, utilityVar string, ignorePolicies bool) (*ve.Factor, error) {
	var decisionEvidence map[string]string
	if ignorePolicies {
		decisionEvidence = evidence
	}

	var err error
	n.ve, n.variableNames, err = n.toVE(decisionEvidence)
	if err != nil {
		return nil, err
	}

	ev := []ve.Evidence{}
	for name, value := range evidence {
		vv, ok := n.variableNames[name]
		if !ok {
			return nil, fmt.Errorf("evidence variable %s not found", name)
		}
		idx := slices.Index(vv.Variable.Outcomes, value)
		if idx < 0 {
			return nil, fmt.Errorf("outcome %s for evidence variable %s not found", value, name)
		}
		ev = append(ev, ve.Evidence{Variable: vv.VeVariable, Value: idx})
	}

	q := make([]ve.Variable, len(query))
	for i, name := range query {
		vv, ok := n.variableNames[name]
		if !ok {
			return nil, fmt.Errorf("query variable %s not found", name)
		}
		q[i] = vv.VeVariable
	}

	if utility {
		var util *ve.Variable
		if utilityVar != "" {
			u, ok := n.variableNames[utilityVar]
			if !ok {
				return nil, fmt.Errorf("utility query variable %s not found", utilityVar)
			}
			util = &u.VeVariable
		}
		return n.ve.SolveUtility(ev, q, util), nil
	} else {
		return n.ve.SolveQuery(ev, q), nil
	}
}

// ToEvidence converts a string variable/value pair to probabilities.
func (n *Network) ToEvidence(variable string, value string) ([]float64, error) {
	vv, ok := n.variableNames[variable]
	if !ok {
		return nil, fmt.Errorf("evidence variable %s not found", variable)
	}
	idx := slices.Index(vv.Variable.Outcomes, value)
	if idx < 0 {
		return nil, fmt.Errorf("outcome %s for evidence variable %s not found", value, variable)
	}
	probs := make([]float64, len(vv.Variable.Outcomes))
	probs[idx] = 1.0
	return probs, nil
}

// toVE creates a Variable Elimination solver from the network.
func (n *Network) toVE(evidence map[string]string) (*ve.VE, map[string]*variable, error) {
	vars := ve.NewVariables()
	varNames := map[string]*variable{}
	varIDs := make([]variable, len(n.variables))
	dependencies := map[ve.Variable][]ve.Variable{}
	totalUtilityName := ""

	// collect variables for lookup
	for i, v := range n.variables {
		// skip total utility node
		if i == n.totalUtilityIndex {
			totalUtilityName = v.Name
			continue
		}
		// treat decision variables with policy as normal change variables
		_, isEvidence := evidence[v.Name]
		if v.Type == ve.DecisionNode && !isEvidence {
			if _, ok := n.policies[v.Name]; ok {
				varIDs[i] = variable{
					Variable:   v,
					VeVariable: vars.AddVariable(i, ve.ChanceNode, uint16(len(v.Outcomes))),
				}
				varNames[v.Name] = &varIDs[i]
				continue
			}
		}

		// for all other variables
		varIDs[i] = variable{
			Variable:   v,
			VeVariable: vars.AddVariable(i, v.Type, uint16(len(v.Outcomes))),
		}
		varNames[v.Name] = &varIDs[i]
	}

	// create factors from tables
	factors := []ve.Factor{}
	for _, f := range n.factors {
		// skip factor for total utility
		if f.For == totalUtilityName {
			continue
		}
		// get primary variable
		forVar, ok := varNames[f.For]
		if !ok {
			return nil, nil, fmt.Errorf("variable %s for factor not found", f.For)
		}

		// collect conditional variables
		variables := make([]ve.Variable, len(f.Given))
		for j, v := range f.Given {
			vv, ok := varNames[v]
			if !ok {
				return nil, nil, fmt.Errorf("variable %s in factor for %s not found", v, f.For)
			}
			variables[j] = vv.VeVariable
		}

		// don't add factors for unsolved decision nodes, but add dependencies
		if forVar.VeVariable.NodeType == ve.DecisionNode {
			dependencies[forVar.VeVariable] = variables
			continue
		}
		// don't add factors for solved decision nodes, done later
		if forVar.Variable.Type == ve.DecisionNode {
			continue
		}

		// append primary variable as last variable of the factor
		variables = append(variables, forVar.VeVariable)

		factor := vars.CreateFactor(variables, f.Table)

		// normalize for primary chance variable
		if forVar.Variable.Type == ve.ChanceNode {
			factor = vars.NormalizeFor(&factor, variables[len(variables)-1])
		}

		// add to list of factors
		factors = append(factors, factor)
	}

	// add policies as factors
	factors = append(factors, n.policyFactors(vars, varIDs, evidence)...)

	weights, err := n.prepareUtilityWeights()
	if err != nil {
		return nil, nil, err
	}

	return ve.New(vars, factors, dependencies, weights), varNames, nil
}

func (n *Network) prepareUtilityWeights() ([]float64, error) {
	utilityNodes := []*Variable{}

	// collect variables for lookup
	for i, v := range n.variables {
		// skip total utility node
		if i == n.totalUtilityIndex {
			continue
		}
		// count utility nodes
		if v.Type == ve.UtilityNode {
			utilityNodes = append(utilityNodes, &v)
		}
	}

	var weights []float64
	if n.totalUtilityIndex < 0 {
		return nil, nil
	}
	node := &n.variables[n.totalUtilityIndex]
	table := node.Factor.Table
	parents := node.Factor.Given
	weights = make([]float64, len(utilityNodes))
	for i := range utilityNodes {
		idx := slices.Index(parents, utilityNodes[i].Name)
		if idx < 0 {
			return nil, fmt.Errorf("utility node %s not included in total utility", utilityNodes[i].Name)
		}
		weights[i] = table[idx]
	}

	return weights, nil
}

// policyFactors collects policies as factors.
func (n *Network) policyFactors(vars *ve.Variables, varIDs []variable, evidence map[string]string) []ve.Factor {
	factors := []ve.Factor{}
	for name, f := range n.policies {
		// if decision variable has evidence (and policies are ignored), don't add a factor
		if _, isEvidence := evidence[name]; isEvidence {
			continue
		}
		// collect variables
		variables := make([]ve.Variable, len(f.Variables()))
		for i, v := range f.Variables() {
			// treat solved decision nodes as chance nodes
			if v.NodeType == ve.DecisionNode {
				vv := varIDs[v.Id]
				if _, ok := n.policies[vv.Variable.Name]; ok {
					v.NodeType = ve.ChanceNode
				}
			}
			// add to list of variables
			variables[i] = v
		}
		// add to list of factors
		factors = append(factors, vars.CreateFactor(variables, f.Data()))
	}

	return factors
}

// Normalize a factor.
func (n *Network) Normalize(f *ve.Factor) ve.Factor {
	return n.ve.Variables.Normalize(f)
}

// NormalizeUtility normalizes utility factor by dividing it by a probability factor.
func (n *Network) NormalizeUtility(utility *ve.Factor, probs *ve.Factor) ve.Factor {
	inv := n.ve.Variables.Invert(probs)
	return n.ve.Variables.Product(utility, &inv)
}

func (n *Network) Marginal(f *ve.Factor, v string) ve.Factor {
	vv, ok := n.variableNames[v]
	if !ok {
		panic(fmt.Sprintf("marginal: variable %s not found", v))
	}
	return n.ve.Variables.Marginal(f, vv.VeVariable)
}

func (n *Network) Rearrange(f *ve.Factor, variables []string) ve.Factor {
	vars := n.rearrangeVariables(f, variables)
	return n.ve.Variables.Rearrange(f, vars)
}

func (n *Network) rearrangeVariables(f *ve.Factor, variables []string) []ve.Variable {
	fVariables := f.Variables()
	vars := make([]ve.Variable, 0, len(fVariables))
	done := make([]bool, len(fVariables))
	for i := 0; i < len(variables)-1; i++ {
		idx, ok := n.variableIndex(f, variables[i])
		if !ok {
			panic(fmt.Sprintf("variable %s to rearrange not in factor", variables[i]))
		}
		vars = append(vars, fVariables[idx])
		done[idx] = true
	}

	idx, ok := n.variableIndex(f, variables[len(variables)-1])
	if !ok {
		panic(fmt.Sprintf("variable %s to rearrange not in factor", variables[len(variables)-1]))
	}
	last := fVariables[idx]

	for i, d := range done {
		if d || fVariables[i].Id == last.Id {
			continue
		}
		vars = append(vars, fVariables[i])
	}

	vars = append(vars, last)

	return vars
}

func (n *Network) variableIndex(f *ve.Factor, v string) (int, bool) {
	variable, ok := n.variableNames[v]
	if !ok {
		panic(fmt.Sprintf("variable %s not found in network", v))
	}
	idx := slices.IndexFunc(f.Variables(), func(vv ve.Variable) bool { return vv.Id == variable.VeVariable.Id })
	if idx < 0 {
		return -1, false
	}
	return idx, true
}
