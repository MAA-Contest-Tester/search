package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use: "classify [subcommand]",
	}
	generate := &cobra.Command{
		Use:     "generate [output]",
		Aliases: []string{"g"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			PickDataset(args[0])
		},
	}
	server := &cobra.Command{
		Use:     "serve [csv] [database]",
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dist, err := cmd.Flags().GetString("dist")
			if err != nil {
				log.Fatal(err)
			}
			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				log.Fatal(err)
			}
			mux := CreateMux(args[0], args[1], dist)
			log.Println("Serving at", ":"+fmt.Sprint(port), "...")
			http.ListenAndServe(":"+fmt.Sprint(port), mux)
		},
	}
	server.Flags().StringP("dist", "D", "", "filesystem path to serve")
	server.Flags().IntP("port", "P", 7828, "filesystem path to serve")
	root.AddCommand(generate, server)
	root.Execute()
}
