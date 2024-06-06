package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	rain := Node{
		Name:   "Rain",
		States: []string{"yes", "no"},
		CPT:    [][]float64{{0.2, 0.8}},
	}

	sprinkler := Node{
		Name:    "Sprinkler",
		Parents: []string{"Rain"},
		States:  []string{"yes", "no"},
		CPT: [][]float64{
			{0.01, 0.99}, // rain yes
			{0.2, 0.8},   // rain no
		},
	}

	grassWet := Node{
		Name:    "GrassWet",
		Parents: []string{"Rain", "Sprinkler"},
		States:  []string{"yes", "no"},
		CPT: [][]float64{
			{0.99, 0.01}, // rain yes, sprikler yes
			{0.8, 0.2},   // rain yes, sprikler no
			{0.9, 0.1},   // rain no, sprikler yes
			{0.0, 1.0},   // rain no, sprikler no
		},
	}

	net, err := New([]*Node{&sprinkler, &grassWet, &rain})
	assert.Nil(t, err)

	assert.Equal(t, "Rain", net.nodes[0].Name)
	assert.Equal(t, "Sprinkler", net.nodes[1].Name)
	assert.Equal(t, "GrassWet", net.nodes[2].Name)

	assert.Equal(t, []int{}, net.nodes[0].Parents)
	assert.Equal(t, []int{0}, net.nodes[1].Parents)
	assert.Equal(t, []int{0, 1}, net.nodes[2].Parents)
}

func TestSortCycles(t *testing.T) {
	a := Node{
		Name:    "A",
		Parents: []string{"C"},
	}

	b := Node{
		Name:    "B",
		Parents: []string{"A"},
	}

	c := Node{
		Name:    "C",
		Parents: []string{"B"},
	}

	_, err := New([]*Node{&c, &a, &b})
	assert.NotNil(t, err)
	assert.Equal(t, "graph has cycles", err.Error())
}
