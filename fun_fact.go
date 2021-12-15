package bird_data_guessing

import "github.com/gbdubs/inference"

func (r *searcher) FunFact(englishName string) []*inference.String {
	pattern := "(\\.|\\n)\\s*(((all|some|most|the|a|an) )?" + englishName + "[a-zA-Z0-9\\-,'\"’‘”“ ]{40,150}\\.)"
	return r.ZZZExtractAllMatch(pattern, 2)
}
