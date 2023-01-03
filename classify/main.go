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
			getQueries(args[0])
		},
	}
	server := &cobra.Command{
		Use:     "serve [database]",
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			load, err := cmd.Flags().GetString("load"); if err != nil {
				log.Fatal(err);
			}
			dist, err := cmd.Flags().GetString("dist"); if err != nil {
				log.Fatal(err)
			}
			port, err := cmd.Flags().GetInt("port"); if err != nil {
				log.Fatal(err)
			}
			mux := CreateMux(load, args[0], dist)
			log.Println("Serving at", ":"+fmt.Sprint(port), "...")
			http.ListenAndServe(":"+fmt.Sprint(port), mux)
		},
	}
	server.Flags().StringP("dist", "D", "", "filesystem path to serve")
	server.Flags().StringP("load", "L", "", "CSV File to Initialize DB With")
	server.Flags().IntP("port", "P", 7828, "Port to serve on")
	root.AddCommand(generate, server)
	root.Execute()
}
