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

// rootCommand sets up the CLI for the TUI.
func rootCommand() *cobra.Command {
	evidence := []string{}
	var training string
	var noData string
	var delim string

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
			ev, err := tui.ParseEvidence(evidence)
			if err != nil {
				return err
			}

			delimRunes := []rune(delim)
			if len(delimRunes) != 1 {
				return fmt.Errorf("argument for --delim must be a single rune; got '%s'", delim)
			}

			a := tui.New(args[0], ev, training, noData, delimRunes[0])
			return a.Run()
		},
	}
	root.Flags().StringSliceVarP(&evidence, "evidence", "e", []string{}, "Evidence in the format:\n    k1=v1,k2=v2,k3=v3")
	root.Flags().StringVarP(&training, "train", "t", "", "train the network from the given file")
	root.Flags().StringVarP(&noData, "no-data", "n", "", "Value for missing data (default \"\")")
	root.Flags().StringVarP(&delim, "delim", "d", ",", "CSV delimiter for training file")

	root.Flags().SortFlags = false

	return &root
}
