package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/MAA-Contest-Tester/search/backend/database"
	"github.com/MAA-Contest-Tester/search/backend/scrape"
	"github.com/MAA-Contest-Tester/search/backend/server"
	"github.com/spf13/cobra"
)

func fileExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error! %v does not exist.\n", path)
		os.Exit(1)
	}
}

func start_server(dir *string) {
	if dir != nil {
		fileExists(*dir)
	}
	mux := server.InitServer(dir)
	log.Println("Running server on port 7827...")
	http.ListenAndServe(":7827", mux)
}

func load_dataset(jsonfile *string) {
	var dataset []scrape.Problem
	if jsonfile != nil {
		fileExists(*jsonfile)
		data, err := os.ReadFile(*jsonfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading %v! %v\n", *jsonfile, err)
			os.Exit(1)
		}
		err = json.Unmarshal(data, &dataset)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while parsing JSON at %v! %v\n", *jsonfile, err)
			os.Exit(1)
		}
		log.Printf("Loading Dataset from %v", *jsonfile)
	} else {
		dataset = scrape.ScrapeList(scrape.ScrapeContestDefaults())
	}
	client := database.Client()
	log.Printf("Inserting %v points into Redis...", len(dataset))
	client.AddProblems(dataset)
	log.Println("Done")
}

func dump_dataset(filename *string) {
	var out io.Writer = os.Stdout
	if filename != nil {
		filename := os.Args[2]
		err := os.MkdirAll(path.Dir(filename), os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating dirs for %v! %v\n", filename, err)
			os.Exit(1)
		}
		out_tmp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while opening %v! %v\n", filename, err)
			os.Exit(1)
		}
		out = out_tmp
	}
	dataset := scrape.ScrapeList(scrape.ScrapeContestDefaults())
	if out == os.Stdout {
		b, _ := json.MarshalIndent(dataset, "", "  ")
		out.Write(b)
	} else {
		b, _ := json.Marshal(dataset)
		out.Write(b)
	}
}

func main() {
	dump := &cobra.Command{Use: "dump [file]", Aliases: []string{"d"}, Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			dump_dataset(&args[0])
		} else {
			dump_dataset(nil)
		}
	}}
	load := &cobra.Command{Use: "load [file]", Aliases: []string{"l"}, Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			load_dataset(&args[0])
		} else {
			load_dataset(nil)
		}
	}}
	server := &cobra.Command{Use: "server [dir]", Aliases: []string{"s"}, Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			start_server(&args[0])
		} else {
			start_server(nil)
		}
	}}
	root := &cobra.Command{Use: "psearch", Short: "A fast search engine for browsing math problems to try"}
	root.AddCommand(dump, load, server)
	root.Execute()
}
