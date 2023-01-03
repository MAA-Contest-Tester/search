package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	extract := &cobra.Command{
		Use:  "extract [database] [csv]",
		Aliases: []string{"e"},
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			db, err := gorm.Open(sqlite.Open(args[0]))
			if err != nil {
				log.Fatal(err)
			}
			w, err := os.OpenFile(args[1], os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644); if err != nil {
				log.Fatal(err);
			}
			result := ExtractDB(db);
			csv.NewWriter(w).WriteAll(result);
		},
	}
	server := &cobra.Command{
		Use:     "serve [database]",
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			load, err := cmd.Flags().GetString("load")
			if err != nil {
				log.Fatal(err)
			}
			dist, err := cmd.Flags().GetString("dist")
			if err != nil {
				log.Fatal(err)
			}
			port, err := cmd.Flags().GetInt("port")
			if err != nil {
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
	root.AddCommand(generate, server, extract)
	root.Execute()
}
