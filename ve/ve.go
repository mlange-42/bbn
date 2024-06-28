package ve

type VE struct {
	variables []*Variables
	factors   []Factor
}

func New(variables []*Variables, factors []Factor) VE {
	return VE{
		variables: variables,
		factors:   factors,
	}
}

func (ve *VE) Eliminate(evidence [][2]int, query []int) {

}
