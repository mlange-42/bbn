package main

import (
	"github.com/mlange-42/bbn"
	"github.com/pkg/profile"
)

// Profiling:
// go build ./benchmark/profile
// profile
// go tool pprof -http=":8000" -nodefraction=0.001 profile cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 profile mem.pprof

func main() {
	count := 1_000_000

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(count)
	stop.Stop()

	stop = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."))
	run(count)
	stop.Stop()
}

func run(count int) {
	for i := 0; i < count; i++ {
		runOnce()
	}
}

func runOnce() {
	variables := []bbn.Variable{
		{Name: "Rain", Outcomes: []string{"yes", "no"}},
		{Name: "Sprinkler", Outcomes: []string{"yes", "no"}},
		{Name: "GrassWet", Outcomes: []string{"yes", "no"}},
	}

	factors := []bbn.Factor{
		{
			For: "Rain",
			Table: []float64{
				// rain+, rain-
				0.2, 0.8,
			},
		},
		{
			For:   "Sprinkler",
			Given: []string{"Rain"},
			Table: []float64{
				// yes   no
				0.01, 0.99, // rain+
				0.2, 0.8, // rain-
			},
		},
		{
			For:   "GrassWet",
			Given: []string{"Rain", "Sprinkler"},
			Table: []float64{
				//   yes   no
				0.99, 0.01, // rain+, sprikler+s
				0.8, 0.2, // rain+, sprikler-
				0.9, 0.1, // rain-, sprikler+
				0.0, 1.0, // rain-, sprikler-
			},
		},
	}

	net, err := bbn.New("Sprinkler", "", variables, factors)
	if err != nil {
		panic(err)
	}

	evidence := map[string]string{
		"GrassWet": "yes",
		"Rain":     "no",
	}
	query := []string{
		"Sprinkler",
	}

	result, _, err := net.SolveQuery(evidence, query, false)
	if err != nil {
		panic(err)
	}

	_ = result
}
