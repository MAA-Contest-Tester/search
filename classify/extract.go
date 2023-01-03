package main

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type DataPoint struct {
	Source    string
	Statement string
	Hits      [numOfFields]int
}

func ExtractDB(g *gorm.DB) [][]string {
	entries := []ClassificationEntry{}
	g.Preload("Problem").Find(&entries);
	entrymap := map[uint]*DataPoint{}
	for _, x := range entries {
		_, exist := entrymap[x.Problem.ID]
		if x.Answer > numOfFields {
			log.Println("Found an entry with answer", x.Answer, "greater than", numOfFields, ", skipping...")
			continue
		}
		if !exist {
			entrymap[x.Problem.ID] = &DataPoint{
				Source:    x.Problem.Source,
				Statement: StripSource(x.Problem.Statement),
			}
		}
		entrymap[x.Problem.ID].Hits[x.Answer]++
	}
	csvdata := make([][]string, 0)
	header := []string{"source", "statement"}
	for i := 0; i < numOfFields; i++ {
		header = append(header, fmt.Sprintf("category%v", i))
	}
	csvdata = append(csvdata, header)
	for _, v := range entrymap {
		point := []string{v.Source, v.Statement}
		for _, x := range v.Hits {
			point = append(point, strconv.Itoa(x))
		}
		csvdata = append(csvdata, point)
	}
	return csvdata
}
