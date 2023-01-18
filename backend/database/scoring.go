package database

import (
	"math"
	"regexp"
	"strconv"

	"github.com/MAA-Contest-Tester/search/backend/scrape"
)

const (
	// contests before this year should be penalized
	THRESHOLD = 2010
)

func problemScore(p scrape.Problem) float32 {
	re := regexp.MustCompile(`\d{4}`)
	num := re.Find([]byte(p.Source))
	if num != nil {
		res, _ := strconv.ParseFloat(string(num), 64)
		return float32(1.0/(1.0 + math.Exp(THRESHOLD - res)))
	} else {
		return 0.1
	}
}
