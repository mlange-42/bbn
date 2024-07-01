package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCommand().Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		fmt.Print("\nRun `bbn -h` for help!\n\n")
		os.Exit(1)
	}
}

// rootCommand sets up the CLI with sub-commands.
func rootCommand() *cobra.Command {
	root := cobra.Command{
		Use:           "bbn",
		Short:         "Bayesian Belief Network.",
		Long:          `Bayesian Belief Network.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	root.AddCommand(inferCommand())
	//root.AddCommand(trainCommand())

	return &root
}
