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

var client database.SearchClient = *database.Client();

func load_dataset(jsonfile *string, forum bool) {
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
		if forum {
			dataset = scrape.ScrapeForumDefaults()
		} else {
			dataset = scrape.ScrapeWikiDefaults()
		}
	}
	log.Printf("Inserting %v points into Redis...", len(dataset))
	client.AddProblems(dataset)
	log.Println("Done")
}

func dump_dataset(filename *string, forum bool) {
	var out io.Writer = os.Stdout
	if filename != nil {
		filename := *filename;
		err := os.MkdirAll(path.Dir(filename), os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating dirs for %v! %v\n", filename, err)
			os.Exit(1)
		}
		fmt.Println("Encountered filename", filename);
		out_tmp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while opening %v! %v\n", filename, err)
			os.Exit(1)
		}
		out = out_tmp
	}
	dataset := []scrape.Problem{};
	if forum {
		dataset = scrape.ScrapeForumDefaults();
	} else {
		dataset = scrape.ScrapeWikiDefaults();
	}
	b, _ := json.MarshalIndent(dataset, "", "  ")
	out.Write(b);
}

func start_server(dir *string, port int, load []string, forum bool) {
	if dir != nil {
		fileExists(*dir)
	}
	if len(load) > 0 {
		client.Drop();
		for _, l := range load {
			load_dataset(&l, forum);
		}
	}
	mux := server.InitServer(dir)
	log.Printf("Running server on port %v...", port)
	http.ListenAndServe(":" + fmt.Sprint(port), mux)
}

func main() {
	dump := &cobra.Command{Use: "dump [file]", Args:cobra.MaximumNArgs(1), Aliases: []string{"d"}, Run: func(cmd *cobra.Command, args []string) {
		forum,err := cmd.InheritedFlags().GetBool("forum"); if err != nil {
			panic(err);
		}
		if len(args) == 1 {
			dump_dataset(&args[0], forum)
		} else {
			dump_dataset(nil, forum)
		}
	}}
	load := &cobra.Command{Use: "load [files...]", Aliases: []string{"l"}, Run: func(cmd *cobra.Command, args []string) {
		forum,err := cmd.InheritedFlags().GetBool("forum"); if err != nil {
			panic(err);
		}
		if len(args) >= 1 {
			client.Drop();
			for _, a := range args {
				load_dataset(&a, forum)
			}
		} else {
			load_dataset(nil, forum)
		}
	}}
	server := &cobra.Command{Use: "server", Aliases: []string{"s"}, Run: func(cmd *cobra.Command, args []string) {
		forum,err := cmd.InheritedFlags().GetBool("forum"); if err != nil {
			panic(err);
		}
		port, _ := cmd.Flags().GetInt("port")
		load,_ := cmd.Flags().GetStringArray("load")
		dir,_ := cmd.Flags().GetString("dir")
		if len(dir) > 0 {
			start_server(&dir, port, load, forum)
		} else {
			start_server(nil, port, load, forum)
		}
	}}
	server.Flags().IntP("port", "P", 7827, "Port to use")
	server.Flags().StringArrayP("load", "L", []string{}, "File to load once connected to Redis");
	server.Flags().StringP("dir", "D", "", "Generated directory to store");

	root := &cobra.Command{Use: "psearch", Short: "A fast search engine for browsing math problems to try"}
	root.PersistentFlags().BoolP("forum", "F", false, "Switch for dumping the AoPS forum dataset")
	root.AddCommand(dump, load, server)
	root.Execute()
}