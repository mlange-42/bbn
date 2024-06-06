package bbn

import (
	"fmt"
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
	result, err := net.Sample(evidence, 10000)
	assert.Nil(t, err)

	assert.Equal(t, []float64{0, 1}, result["Rain"])
	assert.Equal(t, []float64{1, 0}, result["GrassWet"])

	assert.Equal(t, []float64{1, 0}, result["Sprinkler"])

	evidence = map[string]string{
		"Sprinkler": "no",
		"GrassWet":  "no",
	}
	result, err = net.Sample(evidence, 10000)
	assert.Nil(t, err)

	assert.Equal(t, []float64{0, 1}, result["Sprinkler"])
	assert.Equal(t, []float64{0, 1}, result["GrassWet"])

	assert.Less(t, result["Rain"][0], 0.1)
	assert.Greater(t, result["Rain"][1], 0.9)

	fmt.Println(result)
}

func TestMontyHallProblem(t *testing.T) {
	player := Node{
		Name:   "Player",
		States: []string{"D1", "D2", "D3"},
		CPT:    [][]float64{{1, 1, 1}},
	}

	car := Node{
		Name:   "Car",
		States: []string{"D1", "D2", "D3"},
		CPT:    [][]float64{{1, 1, 1}},
	}

	host := Node{
		Name:    "Host",
		Parents: []string{"Player", "Car"},
		States:  []string{"D1", "D2", "D3"},
		CPT: [][]float64{
			{0, 1, 1}, // P1 C1
			{0, 0, 1}, // P1 C2
			{0, 1, 0}, // P1 C3

			{0, 0, 1}, // P2 C1
			{1, 0, 1}, // P2 C2
			{1, 0, 0}, // P2 C3

			{0, 1, 0}, // P3 C1
			{1, 0, 0}, // P3 C2
			{1, 1, 0}, // P3 C3
		},
	}

	net, err := New(&player, &car, &host)
	assert.Nil(t, err)

	evidence := map[string]string{
		"Player": "D1",
		"Host":   "D2",
	}
	result, err := net.Sample(evidence, 100000)
	assert.Nil(t, err)

	fmt.Println(result)
}

func TestMontyHallProblem5Doors(t *testing.T) {
	player := Node{
		Name:   "Player",
		States: []string{"D1", "D2", "D3", "D4", "D5"},
		CPT:    [][]float64{{1, 1, 1, 1, 1}},
	}

	car := Node{
		Name:   "Car",
		States: []string{"D1", "D2", "D3", "D4", "D5"},
		CPT:    [][]float64{{1, 1, 1, 1, 1}},
	}

	hostCPT := make([][]float64, len(player.States)*len(car.States))
	idx := 0
	for p := 0; p < len(player.States); p++ {
		for c := 0; c < len(car.States); c++ {
			probs := make([]float64, 5)
			for i := range probs {
				if i != p && i != c {
					probs[i] = 1
				}
			}
			hostCPT[idx] = probs
			idx++
		}
	}
	host := Node{
		Name:    "Host",
		Parents: []string{"Player", "Car"},
		States:  []string{"D1", "D2", "D3", "D4", "D5"},
		CPT:     hostCPT,
	}

	net, err := New(&player, &car, &host)
	assert.Nil(t, err)

	evidence := map[string]string{
		"Player": "D1",
		"Host":   "D2",
	}
	result, err := net.Sample(evidence, 100000)
	assert.Nil(t, err)

	fmt.Println(result)
}
