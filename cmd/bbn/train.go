package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/mlange-42/bbn"
	"github.com/spf13/cobra"
)

// trainCommand performs network training.
func trainCommand() *cobra.Command {
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

			return runTrainCommand(netFile, datafile)
		},
	}

	return &root
}

func runTrainCommand(networkFile, dataFile string) error {
	net, nodes, err := bbn.FromYAMLFile(networkFile)
	if err != nil {
		return err
	}

	file, err := os.Open(dataFile)
	if err != nil {
		return err
	}

	r := csv.NewReader(file)
	r.ReuseRecord = true

	header, err := r.Read()
	if err != nil {
		return err
	}
	indices := make([]int, len(nodes))
	outcomes := make([]map[string]int, len(nodes))
	for i, node := range nodes {
		idx := slices.Index(header, node.Variable)
		if idx < 0 {
			return fmt.Errorf("no column '%s' in file '%s'", node.Variable, dataFile)
		}
		indices[i] = idx

		outcomes[i] = make(map[string]int, len(node.Outcomes))
		for j, o := range node.Outcomes {
			outcomes[i][o] = j
		}
	}

	train := bbn.NewTrainer(net)
	samples := make([]int, len(nodes))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for i, idx := range indices {
			var ok bool
			samples[i], ok = outcomes[i][record[idx]]
			if !ok {
				return fmt.Errorf("outcome '%s' not available in node '%s'", record[idx], nodes[i].Variable)
			}
			train.AddSample(samples)
		}
	}

	net, err = train.UpdateNetwork()
	if err != nil {
		return err
	}
	yml, err := bbn.ToYAML(net)
	if err != nil {
		return err
	}

	fmt.Println(string(yml))

	return nil
}
