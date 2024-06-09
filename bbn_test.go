package bbn

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	rain := Node{
		Variable: "Rain",
		Outcomes: []string{"yes", "no"},
		Table:    [][]float64{{0.2, 0.8}},
	}

	sprinkler := Node{
		Variable: "Sprinkler",
		Given:    []string{"Rain"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.01, 0.99}, // rain yes
			{0.2, 0.8},   // rain no
		},
	}

	grassWet := Node{
		Variable: "GrassWet",
		Given:    []string{"Rain", "Sprinkler"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.99, 0.01}, // rain yes, sprikler yes
			{0.8, 0.2},   // rain yes, sprikler no
			{0.9, 0.1},   // rain no, sprikler yes
			{0.0, 1.0},   // rain no, sprikler no
		},
	}

	net, err := New("Sprinkler", &sprinkler, &grassWet, &rain)
	assert.Nil(t, err)

	assert.Equal(t, "Rain", net.nodes[0].Variable)
	assert.Equal(t, "Sprinkler", net.nodes[1].Variable)
	assert.Equal(t, "GrassWet", net.nodes[2].Variable)

	assert.Equal(t, []int{}, net.nodes[0].Given)
	assert.Equal(t, []int{0}, net.nodes[1].Given)
	assert.Equal(t, []int{0, 1}, net.nodes[2].Given)

	assert.Equal(t, []int(nil), net.nodes[0].Stride)
	assert.Equal(t, []int{1}, net.nodes[1].Stride)
	assert.Equal(t, []int{2, 1}, net.nodes[2].Stride)
}

func TestStride(t *testing.T) {
	a := Node{
		Variable: "A",
		Outcomes: []string{"a", "b", "c", "d"},
		Table:    [][]float64{{0.25, 0.25, 0.25, 0.25}},
	}

	b := Node{
		Variable: "B",
		Outcomes: []string{"a", "b", "c"},
		Given:    []string{"A"},
		Table: [][]float64{
			{0.333, 0.333, 0.333},
			{0.333, 0.333, 0.333},
			{0.333, 0.333, 0.333},
			{0.333, 0.333, 0.333},
		},
	}

	c := Node{
		Variable: "C",
		Given:    []string{"A", "B"},
		Outcomes: []string{"a", "b"},
		Table: [][]float64{
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
			{0.5, 0.5},
		},
	}

	net, err := New("Test", &a, &b, &c)
	assert.Nil(t, err)

	assert.Equal(t, []int(nil), net.nodes[0].Stride)
	assert.Equal(t, []int{1}, net.nodes[1].Stride)
	assert.Equal(t, []int{3, 1}, net.nodes[2].Stride)
}

func TestSortCycles(t *testing.T) {
	a := Node{
		Variable: "A",
		Given:    []string{"C"},
		Outcomes: []string{"a", "b"},
		Table: [][]float64{
			{0.5, 0.5},
			{0.5, 0.5},
		},
	}

	b := Node{
		Variable: "B",
		Given:    []string{"A"},
		Outcomes: []string{"a", "b"},
		Table: [][]float64{
			{0.5, 0.5},
			{0.5, 0.5},
		},
	}

	c := Node{
		Variable: "C",
		Given:    []string{"B"},
		Outcomes: []string{"a", "b"},
		Table: [][]float64{
			{0.5, 0.5},
			{0.5, 0.5},
		},
	}

	_, err := New("Test", &c, &a, &b)
	assert.NotNil(t, err)
	assert.Equal(t, "graph has cycles", err.Error())
}

func TestSample(t *testing.T) {
	rain := Node{
		Variable: "Rain",
		Outcomes: []string{"yes", "no"},
		Table:    [][]float64{{0.2, 0.8}},
	}

	sprinkler := Node{
		Variable: "Sprinkler",
		Given:    []string{"Rain"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.01, 0.99}, // rain yes
			{0.2, 0.8},   // rain no
		},
	}

	grassWet := Node{
		Variable: "GrassWet",
		Given:    []string{"Rain", "Sprinkler"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.99, 0.01}, // rain yes, sprikler yes
			{0.8, 0.2},   // rain yes, sprikler no
			{0.9, 0.1},   // rain no, sprikler yes
			{0.0, 1.0},   // rain no, sprikler no
		},
	}

	net, err := New("Test", &sprinkler, &grassWet, &rain)
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
}

func BenchmarkSampleSprinkler_1000(b *testing.B) {
	b.StopTimer()

	rain := Node{
		Variable: "Rain",
		Outcomes: []string{"yes", "no"},
		Table:    [][]float64{{0.2, 0.8}},
	}

	sprinkler := Node{
		Variable: "Sprinkler",
		Given:    []string{"Rain"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.01, 0.99}, // rain yes
			{0.2, 0.8},   // rain no
		},
	}

	grassWet := Node{
		Variable: "GrassWet",
		Given:    []string{"Rain", "Sprinkler"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.99, 0.01}, // rain yes, sprikler yes
			{0.8, 0.2},   // rain yes, sprikler no
			{0.9, 0.1},   // rain no, sprikler yes
			{0.0, 1.0},   // rain no, sprikler no
		},
	}

	net, err := New("Test", &sprinkler, &grassWet, &rain)
	assert.Nil(b, err)

	rng := rand.New(rand.NewSource(1))
	var result map[string][]float64

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		result, err = net.Sample(nil, 1000, rng)
	}
	b.StopTimer()

	assert.Nil(b, err)
	_ = result
}
