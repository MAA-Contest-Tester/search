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
	EXPONENT = 40
)

func problemScore(p scrape.Problem) float32 {
	re := regexp.MustCompile(`\d{4}`)
	num := re.Find([]byte(p.Source))
	if num != nil {
		res, _ := strconv.ParseFloat(string(num), 64)
		// protect against overflow
		difference := THRESHOLD - res;
		difference = math.Max(difference, -EXPONENT);
		difference = math.Min(difference, EXPONENT);
		return float32(1.0/(1.0 + math.Exp(THRESHOLD - res)))
	} else {
		return float32(1.0/(1.0 + math.Exp(EXPONENT)))
	}
}
