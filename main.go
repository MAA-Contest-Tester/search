package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MAA-Contest-Tester/search/scrape"
	"github.com/MAA-Contest-Tester/search/database"
	"github.com/MAA-Contest-Tester/search/server"
);

func main() {
  if len(os.Args) < 2 {
    log.Fatal("Not Enough Arguments. Must be one of 'load', 'dataset', and 'server'.")
  }
  switch(os.Args[1]) {
    case "dataset":
      res := scrape.ScrapeList(scrape.ScrapeContestDefaults());
      b,_ := json.MarshalIndent(res, "", "  ");
      fmt.Println(string(b));
      break;
    case "load":
      res := scrape.ScrapeList(scrape.ScrapeContestDefaults());
      client := database.Client();
      log.Println("Inserting into Redis...");
      client.AddProblems(res);
      log.Println("Done");
      break;
    case "server":
      dir := (*string)(nil);
      if len(os.Args) >= 3 {
        dir = &os.Args[2];
      }
      mux := server.InitServer(dir)
      log.Println("Running server on port 7827...")
      http.ListenAndServe(":7827", mux);
    default:
      log.Fatal("Not a valid command. Must be one of 'load', 'dataset', and 'server'.")
  }
}
