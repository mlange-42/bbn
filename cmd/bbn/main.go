package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/internal/tui"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCommand().Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		fmt.Print("\nRun `bbn -h` for help!\n\n")
		os.Exit(1)
	}
}

// rootCommand sets up the CLI
func rootCommand() *cobra.Command {
	evidence := []string{}
	var seed int64
	var samples int

	root := cobra.Command{
		Use:           "bbn [file]",
		Short:         "Bayesian Belief Network.",
		Long:          `Bayesian Belief Network.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, ev, result, err := run(args[0], evidence, samples, seed)
			if err != nil {
				return err
			}

			for _, node := range nodes {
				fmt.Print("                              ")
				states := node.Outcomes
				for _, s := range states {
					fmt.Printf(" %10s", s)
				}
				fmt.Printf("\n%30s", node.Variable)
				probs := result[node.Variable]
				for _, p := range probs {
					fmt.Printf(" %9.3f%%", p*100)
				}
				if _, ok := ev[node.Variable]; ok {
					fmt.Print("  +")
				}
				fmt.Println()
			}

			return nil
		},
	}
	root.Flags().StringSliceVarP(&evidence, "evidence", "e", []string{}, "Evidence in the format:\n    k1=v1,k2=v2,k3=v3")
	root.Flags().Int64Var(&seed, "seed", 0, "Random seed. Seeded with time by default")
	root.Flags().IntVarP(&samples, "samples", "n", 1_000_000, "Number of samples to take")

	root.Flags().SortFlags = false

	return &root
}

func run(path string, evidence []string, samples int, seed int64) ([]*bbn.Node, map[string]string, map[string][]float64, error) {
	nodes, err := bbn.NodesFromYAML(path)
	if err != nil {
		return nil, nil, nil, err
	}

	net, err := bbn.New(nodes...)
	if err != nil {
		return nil, nil, nil, err
	}

	ev, err := tui.ParseEvidence(evidence)
	if err != nil {
		return nil, nil, nil, err
	}

	if seed <= 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(int64(seed)))
	result, err := net.Sample(ev, samples, rng)
	if err != nil {
		return nil, nil, nil, err
	}

	return nodes, ev, result, nil
}
