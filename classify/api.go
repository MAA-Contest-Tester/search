package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	numOfFields = 5
)

type ProblemsStore struct {
	Problems []StoreProblem
	DB       *gorm.DB
}

func (d *ProblemsStore) LoadQueries(csvpath string) *ProblemsStore {
	if d.DB == nil {
		log.Fatal("Database is nil.")
	}
	if len(csvpath) != 0 {
		r, err := os.Open(csvpath)
		if err != nil {
			log.Fatal(err)
		}
		records, err := csv.NewReader(r).ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		for i := 1; i < len(records); i++ {
			p := StoreProblem{Statement: records[i][1], Source:records[i][0]};
			d.Problems = append(d.Problems, p);
		}
		log.Println("Adding Points into database...");
		d.DB.CreateInBatches(&d.Problems, 100);
		log.Println("Done");
	} else {
		log.Println("Loading Points from database...");
		var present []StoreProblem;
		d.DB.Find(&present);
		for _, point := range present {
			d.Problems = append(d.Problems, StoreProblem{Statement: point.Statement, Source: point.Source});
		}
	}
	return d
}

func (d *ProblemsStore) InitDB(sqlitepath string) *ProblemsStore {
	db, err := gorm.Open(sqlite.Open(sqlitepath))
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&ClassificationEntry{}, &StoreProblem{});
	d.DB = db
	return d
}

func (d *ProblemsStore) Random() StoreProblem {
	return d.Problems[rand.Intn(len(d.Problems))]
}

type addQuery struct {
	Source    string `json:"source"`
	Statement string `json:"statement"`
	Answer    int    `json:"answer"`
}

func CreateMux(dataset string, sqlitepath string, static string) *http.ServeMux {
	store := (&ProblemsStore{}).InitDB(sqlitepath).LoadQueries(dataset);
	mux := http.NewServeMux()
	if len(static) != 0 {
		mux.Handle("/", http.FileServer(http.Dir(static)))
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			io.WriteString(w, "{}")
		})
	}
	mux.HandleFunc("/api/choose", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(store.Random())
		if err != nil {
			// impossible.
			panic(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	})
	mux.HandleFunc("/api/add", func(w http.ResponseWriter, r *http.Request) {
		invalid := func(err error) {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "application/json")
			b, e := json.Marshal(map[string]string{
				"error": err.Error(),
			})
			if e != nil {
				panic(e)
			}
			w.Write(b)
		}
		if r.Method != "POST" {
			invalid(errors.New("Not a POST Request"))
			return
		}
		body, _ := io.ReadAll(r.Body)
		query := addQuery{}
		err := json.Unmarshal(body, &query)
		if err != nil {
			invalid(err)
			return
		}
		if !(0 <= query.Answer && query.Answer < numOfFields) {
			invalid(errors.New(fmt.Sprint("Answer must be between 0 and", numOfFields)))
			return
		}
		var problem *StoreProblem;
		store.DB.Model(&StoreProblem{ Statement: query.Statement, Source: query.Source }).First(&problem);
		// handle 404
		if problem == nil {
			w.WriteHeader(http.StatusNotFound);
			w.Header().Add("Content-Type", "application/json");
			b, _ := json.Marshal(map[string]string{
				"error": "nonexistent problem",
			})
			w.Write(b);
			return
		}
		store.DB.Create(&ClassificationEntry{Answer: query.Answer, Problem: *problem});
		w.Header().Add("Content-Type", "application/json")
		b, err := json.Marshal(map[string]interface{}{
			"inserted": query,
		})
		if err != nil {
			panic(err)
		}
		w.Write(b)
	})
	return mux
}
