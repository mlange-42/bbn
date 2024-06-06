package bbn

import (
	"fmt"
	"math/rand"
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

	net, err := New(&sprinkler, &grassWet, &rain)
	assert.Nil(t, err)

	assert.Equal(t, "Rain", net.nodes[0].Name)
	assert.Equal(t, "Sprinkler", net.nodes[1].Name)
	assert.Equal(t, "GrassWet", net.nodes[2].Name)

	assert.Equal(t, []int{}, net.nodes[0].Parents)
	assert.Equal(t, []int{0}, net.nodes[1].Parents)
	assert.Equal(t, []int{0, 1}, net.nodes[2].Parents)

	assert.Equal(t, []int(nil), net.nodes[0].Stride)
	assert.Equal(t, []int{1}, net.nodes[1].Stride)
	assert.Equal(t, []int{2, 1}, net.nodes[2].Stride)
}

func TestStride(t *testing.T) {
	a := Node{
		Name:   "A",
		States: []string{"a", "b", "c", "d"},
	}

	b := Node{
		Name:    "B",
		States:  []string{"a", "b", "c"},
		Parents: []string{"A"},
	}

	c := Node{
		Name:    "C",
		Parents: []string{"A", "B"},
		States:  []string{"a", "b"},
	}

	net, err := New(&a, &b, &c)
	assert.Nil(t, err)

	assert.Equal(t, []int(nil), net.nodes[0].Stride)
	assert.Equal(t, []int{1}, net.nodes[1].Stride)
	assert.Equal(t, []int{3, 1}, net.nodes[2].Stride)
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

	_, err := New(&c, &a, &b)
	assert.NotNil(t, err)
	assert.Equal(t, "graph has cycles", err.Error())
}

func TestSample(t *testing.T) {
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

	net, err := New(&sprinkler, &grassWet, &rain)
	assert.Nil(t, err)

	evidence := map[string]string{
		"Rain":     "no",
		"GrassWet": "yes",
	}
	result, err := net.Sample(evidence, 10000, rand.New(rand.NewSource(1)))
	assert.Nil(t, err)

	assert.Equal(t, []float64{0, 1}, result["Rain"])
	assert.Equal(t, []float64{1, 0}, result["GrassWet"])

	assert.Equal(t, []float64{1, 0}, result["Sprinkler"])

	evidence = map[string]string{
		"Sprinkler": "no",
		"GrassWet":  "no",
	}
	result, err = net.Sample(evidence, 10000, rand.New(rand.NewSource(1)))
	assert.Nil(t, err)

	assert.Equal(t, []float64{0, 1}, result["Sprinkler"])
	assert.Equal(t, []float64{0, 1}, result["GrassWet"])

	assert.Less(t, result["Rain"][0], 0.1)
	assert.Greater(t, result["Rain"][1], 0.9)

	fmt.Println(result)
}
