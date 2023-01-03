package main

import (
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use: "classify [subcommand]",
	}
	generate := &cobra.Command{
		Use:     "generate [output file]",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			PickDataset(args[0])
		},
	}
	root.AddCommand(generate)
	root.Execute()
}
