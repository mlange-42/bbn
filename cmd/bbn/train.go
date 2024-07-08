package main

import (
	"fmt"

	"github.com/mlange-42/bbn"
	"github.com/mlange-42/bbn/internal/tui"
	"github.com/spf13/cobra"
)

// trainCommand performs network training.
func trainCommand() *cobra.Command {
	var delim string
	var noData string

	root := cobra.Command{
		Use:           "train file data-file",
		Short:         "Performs network training.",
		Long:          `Performs network training.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(2),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			netFile := args[0]
			datafile := args[1]

			delimRunes := []rune(delim)
			if len(delimRunes) != 1 {
				return fmt.Errorf("argument for --delim must be a single rune; got '%s'", delim)
			}

			net, err := bbn.FromFile(netFile)
			if err != nil {
				return err
			}

			nodes := net.Variables()
			net, err = tui.TrainNetwork(net, nodes, datafile, noData, delimRunes[0])
			if err != nil {
				return err
			}

			yml, err := bbn.ToYAML(net)
			if err != nil {
				return err
			}

			fmt.Println(string(yml))
			return nil
		},
	}

	root.Flags().StringVarP(&noData, "no-data", "n", "", "Value for missing data (default \"\")")
	root.Flags().StringVarP(&delim, "delim", "d", ",", "CSV delimiter")

	root.Flags().SortFlags = false

	return &root
}
