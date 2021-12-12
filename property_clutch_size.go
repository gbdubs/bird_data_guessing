package bird_data_guessing

import (
	"fmt"
)

func (s *searcher) ClutchSize() *Property {
	nr := `(one|two|three|four|five|six|seven|eight|nine|\d+)(\.\d+)?`
	rr := fmt.Sprintf("(%s to %s|%s ?- ?%s|%s ?– ?%s)", nr, nr, nr, nr, nr, nr)
	np := `[^.\d]{0,20}`
	p := &Property{}
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
		p = s.ExtractAnyMatch(m)
		if p.StringValue != "" {
			break
		}
	}
	if p.StringValue == "" {
		return p
	}
	v := p.StringValue
	for e, a := range englishToArabicNumerals {
		v = caseInsReplace(v, e, a)
	}
	twoParts := caseInsensitiveRegex("(\\d+) ?(to|-|–) ?(\\d+)").FindStringSubmatch(v)
	if twoParts == nil {
		p.IntValue = atoiOrFail(v)
	} else {
		low := atoiOrFail(twoParts[1])
		high := atoiOrFail(twoParts[3])
		p.IntValue = (high + low) / 2
	}
	return p
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
