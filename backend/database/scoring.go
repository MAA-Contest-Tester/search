package database

import (
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/MAA-Contest-Tester/search/backend/scrape"
)

func problemScore(p scrape.Problem) float32 {
	re := regexp.MustCompile(`\d{4}`)
	num := re.Find([]byte(p.Source))
	if num != nil {
		res, _ := strconv.ParseFloat(string(num), 64)
		res = math.Max(0, float64(res-1900))
		year := math.Max(1, float64(time.Now().Year()-1900 + 3))
		return float32(res) / float32(year)
	} else {
		return 1.0
	}
}
