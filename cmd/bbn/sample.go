package main

import (
	"fmt"

	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/internal/tui"
	"github.com/mlange-42/bbn/internal/ve"
	"github.com/spf13/cobra"
)

// inferCommand performs rejection sampling.
func inferCommand() *cobra.Command {
	evidence := []string{}

	root := cobra.Command{
		Use:           "inference file",
		Short:         "Performs inference by variable elimination.",
		Long:          `Performs inference by variable elimination.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			nodes, ev, result, err := runInferenceCommand(args[0], evidence)
			if err != nil {
				return err
			}

			for _, node := range nodes {
				fmt.Print("                              ")
				states := node.Outcomes
				for _, s := range states {
					fmt.Printf(" %10s", s)
				}
				fmt.Printf("\n%30s", node.Name)
				probs := result[node.Name]
				for _, p := range probs {
					if node.Type == ve.UtilityNode {
						fmt.Printf(" %10.3f", p)
					} else {
						fmt.Printf(" %9.3f%%", p*100)
					}
				}
				if _, ok := ev[node.Name]; ok {
					fmt.Print("  +")
				}
				fmt.Println()
			}

			return nil
		},
	}
	root.Flags().StringSliceVarP(&evidence, "evidence", "e", []string{}, "Evidence in the format:\n    k1=v1,k2=v2,k3=v3")

	root.Flags().SortFlags = false

	return &root
}

func runInferenceCommand(path string, evidence []string) ([]bbn.Variable, map[string]string, map[string][]float64, error) {
	net, nodes, err := bbn.FromFile(path)
	if err != nil {
		return nil, nil, nil, err
	}

	ev, err := tui.ParseEvidence(evidence)
	if err != nil {
		return nil, nil, nil, err
	}

	tuiNodes := make([]tui.Node, len(nodes))
	for i, n := range nodes {
		tuiNodes[i] = tui.NewNode(n)
	}

	_, err = net.SolvePolicies(false)
	if err != nil {
		return nil, nil, nil, err
	}

	result, err := tui.Solve(net, ev, tuiNodes)
	if err != nil {
		return nil, nil, nil, err
	}

	return nodes, ev, result, nil
}
