package bird_data_guessing

import (
	"fmt"
	"math"
)

func (r *searcher) Wingspan() *Property {
	number := `(\d+)(\.\d+)?`
	numberRange := fmt.Sprintf(`(%s( ?(to|-|–) ?%s)?)`, number, number)
	e := `[^.\d]{0,30}`
	avgFemale := fmt.Sprintf("(female%saverage%swing.?span|wing.?span%sfemale%saverage|average%sfemale%swing.?span)", e, e, e, e, e, e)
	avg := fmt.Sprintf("(average%swing.?span|wing.?span%saverage)", e, e)
	wingspanPhrases := []string{avgFemale, avg, "(wing.?span)", "(wing)", "(length|long|measure)"}
	units := []string{"(cm|centimeter|centimetre)", "(millimeter|millimetre|mm)", "(meter|m)", "(inches|in)"}
	for _, wingspanPhrase := range wingspanPhrases {
		for _, wingspanPhraseFirst := range []bool{true, false} {
			for i, unit := range units {
				p := fmt.Sprintf(`[^.\d]%s ?%s\)?`, numberRange, unit)
				m := make(map[string]int)
				if wingspanPhraseFirst {
					m[fmt.Sprintf("%s%s%s", wingspanPhrase, e, p)] = 2
				} else {
					m[fmt.Sprintf("%s%s%s", p, e, wingspanPhrase)] = 1
				}
				a := r.ExtractAnyMatch(m)
				if a.StringValue != "" {
					twoParts := caseInsensitiveRegex(`(\d+)(\.\d+)? ?(to|-|–) ?(\d+)`).FindStringSubmatch(a.StringValue)
					var fVal float64
					if twoParts == nil {
						fVal = floatOrFail(a.StringValue)
					} else {
						low := floatOrFail(twoParts[1])
						high := floatOrFail(twoParts[4])
						fVal = (high + low) / 2
					}

					switch i {
					case 0:
						break
					case 1:
						fVal /= 10.0
						break
					case 2:
						fVal *= 100
						break
					case 3:
						fVal *= 2.54
						break
					}
					a.IntValue = int(math.Round(fVal))
					if a.IntValue >= 5 /* Bee Hummingbird */ && a.IntValue <= 375 /* Wandering albatross */ {
						return a
					}
				}
			}
		}
	}
	return &Property{}
}
