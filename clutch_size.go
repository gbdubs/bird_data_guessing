package bird_data_guessing

import (
	"fmt"
	"regexp"

	"github.com/gbdubs/inference"
)

func (s *searcher) ClutchSize() *inference.IntRange {
	nr := `(one|two|three|four|five|six|seven|eight|nine|\d+)(\.\d+)?`
	rr := fmt.Sprintf("(%s to %s|%s ?- ?%s|%s ?– ?%s)", nr, nr, nr, nr, nr, nr)
	np := `[^.\d]{0,20}`
	match := &inference.String{}
	for _, r := range []string{rr, nr} {
		m := make(map[string]int)
		m[fmt.Sprintf("la(id|y)%s%s%segg", np, r, np)] = 2
		m[fmt.Sprintf("egg%sla(id|y)%s%s", np, np, r)] = 2
		m[fmt.Sprintf("clutch size%s%s", np, r)] = 1
		m[fmt.Sprintf("%s%seggs per year", r, np)] = 1
		m[fmt.Sprintf("clutch%s%s", np, r)] = 1
		m[fmt.Sprintf("%s%segg%sla(y|id)", r, np, np)] = 1
		m[fmt.Sprintf("la(id|y)%segg%s%s", np, np, r)] = 2
		m[fmt.Sprintf("number of eggs%s%s", np, r)] = 1
		m[fmt.Sprintf(`Eggs\s*(usually|from)?\s*%s`, r)] = 2
		match = s.ExtractAnyMatch(m)
		if *match != (inference.String{}) {
			break
		}
	}
	if *match == (inference.String{}) {
		return &inference.IntRange{}
	}
	v := match.Value
	for e, a := range englishToArabicNumerals {
		v = caseInsReplace(v, e, a)
	}
	twoParts := caseInsensitiveRegex("(\\d+) ?(to|-|–) ?(\\d+)").FindStringSubmatch(v)
	if twoParts == nil {
		return &inference.IntRange{
			Min:    atoiOrFail(v),
			Max:    atoiOrFail(v),
			Source: match.Source,
		}
	}
	return &inference.IntRange{
		Min:    atoiOrFail(twoParts[1]),
		Max:    atoiOrFail(twoParts[3]),
		Source: match.Source,
	}
}

var englishToArabicNumerals = map[string]string{
	"zero":  "0",
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
}

func caseInsReplace(input string, lookFor string, replaceWith string) string {
	return regexp.MustCompile("(?i)"+lookFor).ReplaceAllString(input, replaceWith)
}
