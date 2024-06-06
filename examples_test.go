package bbn_test

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/bbn"
)

func Example_sprinkler() {
	rain := bbn.Node{
		Name:   "Rain",
		States: []string{"yes", "no"},
		CPT:    [][]float64{{0.2, 0.8}},
	}

	sprinkler := bbn.Node{
		Name:    "Sprinkler",
		Parents: []string{"Rain"},
		States:  []string{"yes", "no"},
		CPT: [][]float64{
			{0.01, 0.99}, // rain yes
			{0.2, 0.8},   // rain no
		},
	}

	grassWet := bbn.Node{
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

	net, err := bbn.New(&sprinkler, &grassWet, &rain)
	if err != nil {
		panic(err)
	}

	evidence := map[string]string{
		"Rain":     "no",
		"GrassWet": "yes",
	}

	result, err := net.Sample(evidence, 10000, rand.New(rand.NewSource(1)))
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	// Output:
	//map[GrassWet:[1 0] Rain:[0 1] Sprinkler:[1 0]]
}

func Example_montyHall() {
	player := bbn.Node{
		Name:   "Player",
		States: []string{"D1", "D2", "D3"},
		CPT:    [][]float64{{1, 1, 1}},
	}

	car := bbn.Node{
		Name:   "Car",
		States: []string{"D1", "D2", "D3"},
		CPT:    [][]float64{{1, 1, 1}},
	}

	host := bbn.Node{
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

	net, err := bbn.New(&player, &car, &host)
	if err != nil {
		panic(err)
	}

	evidence := map[string]string{
		"Player": "D1",
		"Host":   "D2",
	}

	result, err := net.Sample(evidence, 100000, rand.New(rand.NewSource(1)))
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	// Output:
	// map[Car:[0.3379409820301701 0 0.6620590179698299] Host:[0 1 0] Player:[1 0 0]]
}
