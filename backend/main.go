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

var client database.SearchClient = *database.Client()

func loadDataset(jsonfile *string) {
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
		dataset = scrape.ScrapeForumDefaults()
	}
	log.Printf("Inserting %v points into Redis...", len(dataset))
	client.AddProblems(dataset)
	log.Println("Done")
}

func dumpDataset(filename *string, contests string) {
	var out io.Writer = os.Stdout
	if filename != nil {
		filename := *filename
		err := os.MkdirAll(path.Dir(filename), os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating dirs for %v! %v\n", filename, err)
			os.Exit(1)
		}
		out_tmp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while opening %v! %v\n", filename, err)
			os.Exit(1)
		}
		out = out_tmp
	}
	var dataset []scrape.Problem
	if len(contests) == 0 {
		dataset = scrape.ScrapeForumDefaults()
	} else {
		var categories []int;
		fileExists(contests);
		b, err := os.ReadFile(contests); if err != nil {
			log.Fatal(err);
			os.Exit(1);
		}
		if json.Unmarshal(b, &categories) != nil {
			log.Fatal(err);
			os.Exit(1);
		}
		dataset = scrape.ScrapeForumCategories(categories);
	}
	b, _ := json.MarshalIndent(dataset, "", "  ")
	out.Write(b)
}

func startServer(dir *string, port int, load []string) {
	if dir != nil {
		fileExists(*dir)
	}
	if len(load) > 0 {
		client.Drop()
		for _, l := range load {
			loadDataset(&l)
		}
	}
	mux := server.InitServer(dir)
	log.Printf("Running server on port %v...", port)
	http.ListenAndServe(":"+fmt.Sprint(port), mux)
}

func main() {
	dump := &cobra.Command{Use: "dump [file]", Args: cobra.MaximumNArgs(1), Aliases: []string{"d"}, Run: func(cmd *cobra.Command, args []string) {
		contests, _ := cmd.Flags().GetString("contests");
		if len(args) == 1 {
			dumpDataset(&args[0], contests)
		} else {
			dumpDataset(nil, contests)
		}
	}}
	dump.Flags().StringP("contests", "C", "", "list of contests to parse (everything if this file is not specified)")
	load := &cobra.Command{Use: "load [files...]", Aliases: []string{"l"}, Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 {
			client.Drop()
			for _, a := range args {
				loadDataset(&a)
			}
		} else {
			loadDataset(nil)
		}
	}}
	server := &cobra.Command{Use: "server", Aliases: []string{"s"}, Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		load, _ := cmd.Flags().GetStringArray("load")
		dir, _ := cmd.Flags().GetString("dir")
		if len(dir) > 0 {
			startServer(&dir, port, load)
		} else {
			startServer(nil, port, load)
		}
	}}
	server.Flags().IntP("port", "P", 7827, "Port to use")
	server.Flags().StringArrayP("load", "L", []string{}, "File to load once connected to Redis")
	server.Flags().StringP("dir", "D", "", "Generated directory to store")

	root := &cobra.Command{Use: "psearch", Short: "A fast search engine for browsing math problems to try"}
	root.AddCommand(dump, load, server)
	root.Execute()
}
