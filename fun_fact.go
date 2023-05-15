package bird_data_guessing

import (
	"fmt"
	"strings"

	"github.com/gbdubs/inference"
)

const funFactPattern = `(\.|\n|^|Cool Facts)\s*([^\(\)."]{0,150}%s[^\(\)."]{0,150}\.)`
const prohibitedFunFactPattern = `(copyright|2021)`

func (r *searcher) FunFact(englishName string) []*inference.String {
	result := make([]*inference.String, 0)
	prohibited := caseInsensitiveRegex(prohibitedFunFactPattern)
	for _, name := range namesToSearchFunFactsFor(englishName) {
		p := fmt.Sprintf(funFactPattern, name)
		for _, match := range r.ExtractAllMatch(p, 2) {
			if len(match.Value) > 150 {
				continue
			}
			if prohibited.MatchString(match.Value) {
				continue
			}
			result = append(result, match)
		}
		if len(result) > 0 {
			return result
		}
	}
	return result
}

func namesToSearchFunFactsFor(english string) []string {
	name := strings.ReplaceAll(english, "'", ".?")
	name = strings.ReplaceAll(name, "-", ".?")
	name = strings.ReplaceAll(name, "-", ".?")

	result := []string{name}

	parts := strings.Split(name, " ")
	if len(parts) == 2 {
		result = append(result, parts[1])
	}
	if len(parts) == 3 {
		result = append(result, parts[0]+" "+parts[1])
		result = append(result, parts[1]+" "+parts[2])
	}
	if len(parts) == 4 {
		result = append(result, parts[0]+" "+parts[1])
		result = append(result, parts[1]+" "+parts[2])
		result = append(result, parts[2]+" "+parts[3])
	}
	return result
}
