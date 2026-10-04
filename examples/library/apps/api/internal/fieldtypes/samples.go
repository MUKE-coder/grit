package fieldtypes

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/datatypes"
)

// Sample values for seeders. Each takes the row number and mixes it in, so a
// thousand seeded rows on a unique column do not collide, and each returns a
// value its own rule accepts. The fieldtypes tests check both.

var nonLetters = regexp.MustCompile("[^a-z0-9]+")

func slug(s string) string {
	s = strings.Trim(nonLetters.ReplaceAllString(strings.ToLower(s), ""), "-")
	if s == "" {
		return "user"
	}
	return s
}

// SampleEmail builds first.last.N@domain.
func SampleEmail(first, last, domain string, n int) string {
	d, err := Domain(domain)
	if err != nil {
		d = "example.com"
	}
	return fmt.Sprintf("%s.%s.%d@%s", slug(first), slug(last), n, d)
}

var samplePaths = []string{"pricing", "about", "blog", "docs", "careers", "contact", "products", "support"}

// SampleURL builds https://domain/page/N.
func SampleURL(domain string, n int) string {
	d, err := Domain(domain)
	if err != nil {
		d = "example.com"
	}
	return fmt.Sprintf("https://%s/%s/%d", d, samplePaths[n%len(samplePaths)], n)
}

var sampleTLDs = []string{"com", "co.ug", "io", "org", "net", "dev", "co.uk", "co.ke", "app", "africa"}

// SampleDomain builds wordN.tld from a spread of top-level domains.
func SampleDomain(word string, n int) string {
	return fmt.Sprintf("%s%d.%s", slug(word), n, sampleTLDs[n%len(sampleTLDs)])
}

// SampleCountry walks the list with a stride that is coprime with its length,
// so consecutive rows land on countries all over it rather than AD, AE, AF.
func SampleCountry(n int) string {
	if n < 0 {
		n = -n
	}
	return CountryCodes[(n*97)%len(CountryCodes)]
}

// SampleColor spreads rows over the colour space. The multiplier is odd, so
// the first 16.7 million rows get different colours.
func SampleColor(n int) string {
	m := int64(n % 0x1000000)
	if m < 0 {
		m = -m
	}
	return fmt.Sprintf("#%06x", (m*2654435761)%0x1000000)
}

// SamplePercent is a value with at most one decimal, clustered below 40 the
// way discounts, tax rates and completion figures tend to be.
func SamplePercent(n int) float64 {
	v := (n*7919 + 13) % 1001 // 0..1000
	if n%3 != 0 {
		v = v * 4 / 10
	}
	return float64(v) / 10
}

// sampleStars leans towards four and five out of five, which is what real
// review distributions look like.
var sampleStars = []int{5, 4, 5, 3, 4, 5, 4, 2, 5, 4, 1, 4, 5, 3, 4, 5}

// SampleRating is a whole number of stars from 1 to max.
func SampleRating(n, max int) int {
	if max <= 0 {
		max = 5
	}
	s := sampleStars[n%len(sampleStars)]
	return (s-1)*(max-1)/4 + 1
}

// SampleTime is a time in business hours, 08:00 to 17:45 in quarter hours.
func SampleTime(n int) string {
	slot := (n * 7) % 40
	return fmt.Sprintf("%02d:%02d", 8+slot/4, (slot%4)*15)
}

var samplePlans = []string{"free", "starter", "pro", "enterprise"}

// SampleJSON is a small object, the kind of settings blob a json column holds.
func SampleJSON(n int) datatypes.JSON {
	return datatypes.JSON(fmt.Sprintf("{\"plan\":%q,\"seats\":%d,\"trial\":%t,\"tags\":[\"row-%d\"]}",
		samplePlans[n%len(samplePlans)], 1+n%50, n%4 == 0, n))
}
