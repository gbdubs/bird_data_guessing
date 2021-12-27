package bird_data_guessing

import (
	"fmt"

	"github.com/gbdubs/inference"
)

func (s *searcher) Wingspan() []*inference.Float64Range {
	number := `(\d+)(\.\d+)?`
	numberRange := fmt.Sprintf(`(%s( ?(to|-|–) ?%s)?)`, number, number)
	e := `[^.\d]{0,30}`

	avgFemaleWingspan := fmt.Sprintf("(female%saverage%swing.?span|wing.?span%sfemale%saverage|average%sfemale%swing.?span)", e, e, e, e, e, e)
	avgMaleWingspan := fmt.Sprintf("(male%saverage%swing.?span|wing.?span%smale%saverage|average%smale%swing.?span)", e, e, e, e, e, e)
	avgWingspan := fmt.Sprintf("(average%swing.?span|wing.?span%saverage)", e, e)
	avgPhrases := []string{avgFemaleWingspan, avgMaleWingspan, avgWingspan}

	femaleWingspan := fmt.Sprintf("(female%swing.?span|wing.?span%sfemale)", e, e)
	maleWingspan := fmt.Sprintf("(male%swing.?span|wing.?span%smale)", e, e)
	wingspan := "(wing.?span)"
	wingspanPhrases := []string{femaleWingspan, maleWingspan, wingspan}

	units := []string{"(cm|centimeter|centimetre)", "(millimeter|millimetre|mm)", "(meter|m)", "(inches|in)"}

	findAllMatches := func(keywords []string) []*inference.Float64Range {
		results := make([]*inference.Float64Range, 0)
		for _, keyword := range keywords {
		nextKeyword:
			for _, keywordFirst := range []bool{true, false} {
				for unitIndex, unit := range units {
					scalingFactor := 1.0
					if unitIndex == 1 {
						scalingFactor = 0.1
					} else if unitIndex == 2 {
						scalingFactor = 100.0
					} else if unitIndex == 3 {
						scalingFactor = 2.54
					}

					p := fmt.Sprintf(`[^.\d](%s ?%s)\)?`, numberRange, unit)
					m := make(map[string]int)
					if keywordFirst {
						m[fmt.Sprintf("%s%s%s", keyword, e, p)] = 2
					} else {
						m[fmt.Sprintf("%s%s%s", p, e, keyword)] = 1
					}
					matches := s.ExtractAllMatches(m)
					if len(matches) > 0 {
						for _, match := range matches {
							twoParts := caseInsensitiveRegex(`(\d+)(\.\d+)? ?(to|-|–) ?(\d+)(\.\d+)? ?` + unit).FindStringSubmatch(match.Value)
							onePart := caseInsensitiveRegex(`(\d+(\.\d+)?) ?` + unit).FindStringSubmatch(match.Value)
							if twoParts == nil && onePart == nil {
								panic(fmt.Errorf("Error in regex structure for input %s.", match.Value))
							} else if twoParts == nil {
								v := floatOrFail(onePart[1]) * scalingFactor
								results = append(results, &inference.Float64Range{
									Min:    v,
									Max:    v,
									Source: match.Source,
								})
							} else {
								min := floatOrFail(twoParts[1]) * scalingFactor
								max := floatOrFail(twoParts[4]) * scalingFactor
								results = append(results, &inference.Float64Range{
									Min:    min,
									Max:    max,
									Source: match.Source,
								})
							}
						}
						break nextKeyword
					}
				}
			}
		}
		return results
	}
	return append(findAllMatches(avgPhrases), findAllMatches(wingspanPhrases)...)
}
