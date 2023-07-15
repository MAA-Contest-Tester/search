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
	EXPONENT  = 40
)

func extractYear(p scrape.Problem) float64 {
	re := regexp.MustCompile(`\d{4}`)
	num := re.Find([]byte(p.Source))
	if num != nil {
		res, _ := strconv.ParseFloat(string(num), 64)
		return res
	} else {
		return 1900.0
	}
}

func problemScore(p scrape.Problem) float32 {
	return float32(1.0 / (1.0 + 0.1*math.Exp(THRESHOLD-extractYear(p))))
}
