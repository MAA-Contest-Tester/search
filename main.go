package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/MAA-Contest-Tester/search/database"
	"github.com/MAA-Contest-Tester/search/scrape"
	"github.com/MAA-Contest-Tester/search/server"
);

func fileExists(path string) {
  if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Fprintf(os.Stderr, "Error! %v does not exist.\n", path);
    os.Exit(1);
  }
}

func start_server() {
  dir := (*string)(nil);
  if len(os.Args) >= 3 {
    dir = &os.Args[2];
    fileExists(*dir);
  }
  mux := server.InitServer(dir)
  log.Println("Running server on port 7827...")
  http.ListenAndServe(":7827", mux);
}

func load_dataset() {
  jsonfile := (*string)(nil);
  var dataset []scrape.Problem;
  if len(os.Args) >= 3 {
    jsonfile = &os.Args[2];
    fileExists(*jsonfile);
    data, err := os.ReadFile(*jsonfile);
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error while reading %v! %v\n", *jsonfile, err);
      os.Exit(1);
    }
    err = json.Unmarshal(data, &dataset);
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error while parsing JSON at %v! %v\n", *jsonfile, err);
      os.Exit(1);
    }
    log.Printf("Loading Dataset from %v", *jsonfile)
  } else {
    dataset = scrape.ScrapeList(scrape.ScrapeContestDefaults());
  }
  client := database.Client();
  log.Printf("Inserting %v points into Redis...", len(dataset));
  client.AddProblems(dataset);
  log.Println("Done");
}

func dump_dataset() {
  var out io.Writer = os.Stdout;
  if len(os.Args) >= 3 {
    filename := os.Args[2];
    err := os.MkdirAll(path.Dir(filename), os.ModePerm);
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error while creating dirs for %v! %v\n", filename, err);
      os.Exit(1);
    }
    out_tmp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644);
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error while opening %v! %v\n", filename, err);
      os.Exit(1);
    }
    out = out_tmp;
  }
  dataset := scrape.ScrapeList(scrape.ScrapeContestDefaults());
  if out == os.Stdout {
    b,_ := json.MarshalIndent(dataset, "", "  ");
    out.Write(b);
  } else {
    b,_ := json.Marshal(dataset);
    out.Write(b);
  }
}

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintf(os.Stderr, "Not enough Arguments. The first command must be one of {dump,load,server}\n");
    os.Exit(1);
  }
  switch(os.Args[1]) {
    case "dump":
      dump_dataset();
      break;
    case "load":
      load_dataset();
      break;
    case "server":
      start_server();
    default:
      log.Fatal("Not a valid command. Must be one of 'load', 'dataset', and 'server'.")
  }
}
