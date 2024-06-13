package bbn_test

import (
	"fmt"
	"math/rand"

	"github.com/mlange-42/bbn"
)

func Example_sprinkler() {
	rain := bbn.Node{
		Variable: "Rain",
		Outcomes: []string{"yes", "no"},
		Table:    [][]float64{{0.2, 0.8}},
	}

	sprinkler := bbn.Node{
		Variable: "Sprinkler",
		Given:    []string{"Rain"},
		Outcomes: []string{"yes", "no"},
		Table: [][]float64{
			{0.01, 0.99}, // rain yes
			{0.2, 0.8},   // rain no
		},
	}

	grassWet := bbn.Node{
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

	net, err := bbn.New("Sprinkler", &sprinkler, &grassWet, &rain)
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
		Variable: "Player",
		Outcomes: []string{"D1", "D2", "D3"},
		Table:    [][]float64{{1, 1, 1}},
	}

	car := bbn.Node{
		Variable: "Car",
		Outcomes: []string{"D1", "D2", "D3"},
		Table:    [][]float64{{1, 1, 1}},
	}

	host := bbn.Node{
		Variable: "Host",
		Given:    []string{"Player", "Car"},
		Outcomes: []string{"D1", "D2", "D3"},
		Table: [][]float64{
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

	net, err := bbn.New("Monty-Hall", &player, &car, &host)
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
	// map[Car:[0.3341143251930641 0 0.6658856748071387] Host:[0 1 0] Player:[1 0 0]]
}
