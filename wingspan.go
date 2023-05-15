package bird_data_guessing

import (
	"fmt"

	"github.com/gbdubs/inference"
)

func (s *searcher) Wingspan() []*inference.Float64Range {
	number := `(\d+)(\.\d+)?`
	numberRange := fmt.Sprintf(`((%s) ?(to|-|–) ?(%s))`, number, number)
	e := `[^.\d]{0,30}`
	optionallyOtherMeasurementTwenty := fmt.Sprintf(`[^.;,\d]{0,5}(%s|%s)?[^.;,\d]{0,15}`, number, numberRange)

	avgFemaleWingspan := fmt.Sprintf("(female%saverage%swing.?span|wing.?span%sfemale%saverage|average%sfemale%swing.?span)", e, e, e, e, e, e)
	avgMaleWingspan := fmt.Sprintf("(male%saverage%swing.?span|wing.?span%smale%saverage|average%smale%swing.?span)", e, e, e, e, e, e)
	avgWingspan := fmt.Sprintf("(average%swing.?span|wing.?span%saverag)", e, e)
	avgPhrases := []string{avgFemaleWingspan, avgMaleWingspan, avgWingspan}

	femaleWingspan := fmt.Sprintf("(female%swing.?span|wing.?span%sfemale)", e, e)
	maleWingspan := fmt.Sprintf("(male%swing.?span|wing.?span%smale)", e, e)
	wingspan := "(wing.?span)"
	wingspanPhrases := []string{femaleWingspan, maleWingspan, wingspan}

	avgWingChord := fmt.Sprintf("(average%swing.?chord|wing.?chord%saverag)", e, e)
	wingChord := "(wing.?chord)"
	wingChordPhrases := []string{avgWingChord, wingChord}

	units := []string{"(cm|centimeter|centimetre)", "(millimeter|millimetre|mm)", "(meter|m)", "(inches|in)", "(feet|ft)"}

	findAllMatches := func(keywords []string, multipler float64) []*inference.Float64Range {
		results := make([]*inference.Float64Range, 0)
		for _, keyword := range keywords {
		nextKeyword:
			for _, keywordFirst := range []bool{true, false} {
				for unitIndex, unit := range units {
					for _, useRange := range []bool{true, false} {
						scalingFactor := multipler
						if unitIndex == 1 {
							scalingFactor = 0.1
						} else if unitIndex == 2 {
							scalingFactor = 100.0
						} else if unitIndex == 3 {
							scalingFactor = 2.54
						} else if unitIndex == 4 {
							scalingFactor = 30.48
						}
						rawPattern := number
						if useRange {
							rawPattern = numberRange
						}
						pattern := fmt.Sprintf(`[^.\d](%s ?%s)\)?`, rawPattern, unit)
						m := make(map[string]int)
						if keywordFirst {
							m[fmt.Sprintf("%s%s%s", keyword, optionallyOtherMeasurementTwenty, pattern)] = 13
						} else {
							m[fmt.Sprintf("%s%s%s", pattern, optionallyOtherMeasurementTwenty, keyword)] = 1
						}
						matches := s.ExtractAllMatches(m)
						if len(matches) > 0 {
							for _, match := range matches {
								twoParts := caseInsensitiveRegex(numberRange + " ?" + unit).FindStringSubmatch(match.Value)
								onePart := caseInsensitiveRegex(number + " ?" + unit).FindStringSubmatch(match.Value)
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
									min := floatOrFail(twoParts[2]) * scalingFactor
									max := floatOrFail(twoParts[6]) * scalingFactor
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
		}
		return results
	}

	result := append(findAllMatches(avgPhrases, 1.0), findAllMatches(wingspanPhrases, 1.0)...)
	if len(result) > 0 {
		return result
	}
	return findAllMatches(wingChordPhrases, 2.25)
}
