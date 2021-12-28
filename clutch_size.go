package bird_data_guessing

import (
	"fmt"
	"regexp"

	"github.com/gbdubs/inference"
)

func (s *searcher) ClutchSize() *inference.IntRange {
	number := `(one|two|three|four|five|six|seven|eight|nine|\d+)`
	other := `[^.\d]{0,30}`
	spaceOrSep := `[\s:,\)-]*`

	createPatterns := func(pattern string) map[string]int {
		return map[string]int{
			fmt.Sprintf(`clutch%s%s%segg`, other, pattern, other):                            1,
			fmt.Sprintf(`%s%seggs?%slaid`, pattern, other, other):                            1,
			fmt.Sprintf(`clutch size%s%s`, other, pattern):                                   1,
			fmt.Sprintf(`%s%seggs? per year`, pattern, other):                                1,
			fmt.Sprintf(`%s%segg%sla(y|id)`, pattern, other, other):                          1,
			fmt.Sprintf(`number of eggs%s%s`, other, pattern):                                1,
			fmt.Sprintf(`(clutch|brood|nest|lay|laid)%s%s%segg`, other, pattern, spaceOrSep): 2,
			fmt.Sprintf(`la(id|y)%s%s%segg`, other, pattern, other):                          2,
			fmt.Sprintf(`eggs%s(usually|from)?%s%s`, spaceOrSep, spaceOrSep, pattern):        2,
			//fmt.Sprintf(`egg%sla(id|y)%s%s`, other, other, pattern):                     2,
		}
	}

	numberRange := fmt.Sprintf("(%s ?(or|to|and|-|-|–) ?%s)", number, number)
	matches := s.ExtractAllMatches(createPatterns(numberRange))
	for _, match := range matches {
		if *match != (inference.String{}) {
			v := match.Value
			for e, a := range englishToArabicNumerals {
				v = caseInsReplace(v, e, a)
			}
			twoParts := caseInsensitiveRegex(numberRange).FindStringSubmatch(v)
			if twoParts == nil {
				panic(fmt.Errorf("Error in regex format %s - didn't match submatches in %s.", numberRange, v))
			}
			min := atoiOrFail(twoParts[2])
			max := atoiOrFail(twoParts[4])
			if min > max {
				min, max = max, min
			}
			result := &inference.IntRange{
				Min:    min,
				Max:    max,
				Source: match.Source,
			}
			if result.Max > 33 || result.Min <= 0 {
				continue
			}
			return result
		}
	}

	upToRange := fmt.Sprintf(`((as many as|up to|as large as|as much as)[^\d]{0,10}(%s))`, number)
	matches = s.ExtractAllMatches(createPatterns(upToRange))
	for _, match := range matches {
		v := match.Value
		for e, a := range englishToArabicNumerals {
			v = caseInsReplace(v, e, a)
		}
		submatch := caseInsensitiveRegex(upToRange).FindStringSubmatch(v)
		if submatch == nil {
			panic(fmt.Errorf("Error in regex format %s - didn't match submatches in %s.", upToRange, v))
		}
		result := &inference.IntRange{
			Min:    -1,
			Max:    atoiOrFail(submatch[3]),
			Source: match.Source,
		}
		if result.Max > 33 {
			continue
		}
		return result
	}

	matches = s.ExtractAllMatches(createPatterns(number))
	for _, match := range matches {
		v := match.Value
		for e, a := range englishToArabicNumerals {
			v = caseInsReplace(v, e, a)
		}
		result := &inference.IntRange{
			Min:    atoiOrFail(v),
			Max:    atoiOrFail(v),
			Source: match.Source,
		}
		if result.Max > 33 {
			continue
		}
		return result
	}
	return &inference.IntRange{}
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
