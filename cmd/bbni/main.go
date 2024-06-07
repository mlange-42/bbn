package main

import (
	"fmt"
	"os"

	"github.com/mlange-42/bbn/internal/tui"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCommand().Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		fmt.Print("\nRun `bbni -h` for help!\n\n")
		os.Exit(1)
	}
}

// rootCommand sets up the CLI
func rootCommand() *cobra.Command {
	var seed int64
	var samples int

	root := cobra.Command{
		Use:           "bbni [file]",
		Short:         "Bayesian Belief Network interactive TUI.",
		Long:          `Bayesian Belief Network interactive TUI.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			a := tui.New(args[0], samples, seed)
			return a.Run()
		},
	}
	root.Flags().Int64Var(&seed, "seed", 0, "Random seed. Seeded with time by default")
	root.Flags().IntVarP(&samples, "samples", "n", 1_000_000, "Number of samples to take")

	root.Flags().SortFlags = false

	return &root
}
