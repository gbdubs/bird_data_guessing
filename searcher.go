package bird_data_guessing

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"sync"

	"github.com/gbdubs/attributions"
	"github.com/gbdubs/inference"
)

type searcher struct {
	text        string
	attribution *attributions.Attribution
}

func attributedSearch(a *attributions.Attribution, s string) searcher {
	// Strips out Unicode Spaces
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.ReplaceAll(s, "\\s+", " ")
	return searcher{
		text:        s,
		attribution: a,
	}
}

const defaultContextWidth = 30

func (s *searcher) CountMatches(searchFor ...string) *inference.Int {
	if s.text == "" {
		return &inference.Int{}
	}
	components := make([]*inference.Int, 0)
	for key, matches := range s.getMatches(searchFor...) {
		count := len(matches)
		source := inference.CombineSources(
			fmt.Sprintf("count of %s", key),
			count,
			inference.AsSourceables(matches)...)
		components = append(components, &inference.Int{
			Value:  count,
			Source: source,
		})
	}
	return inference.SumInt(components...)
}

func (s *searcher) getMatches(searchFor ...string) map[string][]*inference.String {
	m := make(map[string][]*inference.String)
	for _, r := range searchFor {
		regex := caseInsensitiveRegex(r)
		matchIndices := regex.FindAllStringIndex(s.text, -1)
		m[r] = make([]*inference.String, len(matchIndices))
		for i, matchIndices := range matchIndices {
			match := s.text[matchIndices[0]:matchIndices[1]]
			start := matchIndices[0] - defaultContextWidth/2
			end := matchIndices[1] + defaultContextWidth/2
			m[r][i] = inference.NewString(match, s.boundedSubstring(start, end), s.attribution)
		}
	}
	return m
}

func (s *searcher) boundedSubstring(start int, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s.text) {
		end = len(s.text)
	}
	return s.text[start:end]
}

func (s *searcher) ExtractAllMatch(pattern string, captureGroup int) []*inference.String {
	if s.text == "" {
		return []*inference.String{}
	}
	matches := caseInsensitiveRegex(pattern).FindAllStringSubmatchIndex(s.text, -1)
	if len(matches) == 0 {
		return []*inference.String{}
	}
	result := make([]*inference.String, len(matches))
	for i, indexes := range matches {
		captureMatch := s.text[indexes[2*captureGroup]:indexes[2*captureGroup+1]]
		context := s.boundedSubstring(indexes[0]-defaultContextWidth/2, indexes[1]+defaultContextWidth/2)
		result[i] = inference.NewString(
			captureMatch,
			context,
			s.attribution)
	}
	return result
}

func (s *searcher) ExtractMatch(pattern string, captureGroup int) *inference.String {
	if s.text == "" {
		return &inference.String{}
	}
	matches := caseInsensitiveRegex(pattern).FindAllStringSubmatchIndex(s.text, -1)
	if len(matches) == 0 {
		return &inference.String{}
	}
	indexes := matches[rand.Intn(len(matches))]
	captureMatch := s.text[indexes[2*captureGroup]:indexes[2*captureGroup+1]]
	context := s.boundedSubstring(indexes[0]-defaultContextWidth/2, indexes[1]+defaultContextWidth/2)
	return inference.NewString(
		captureMatch,
		context,
		s.attribution)
}

func (s *searcher) ExtractAnyMatch(patternToCaptureGroup map[string]int) *inference.String {
	for p, cg := range patternToCaptureGroup {
		m := s.ExtractMatch(p, cg)
		if *m != (inference.String{}) {
			return m
		}
	}
	return &inference.String{}
}

func (s *searcher) ExtractAllMatches(patternToCaptureGroup map[string]int) []*inference.String {
	result := make([]*inference.String, 0)
	for p, cg := range patternToCaptureGroup {
		result = append(result, s.ExtractAllMatch(p, cg)...)
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var regexesLock = sync.RWMutex{}
var regexes = make(map[string]*regexp.Regexp)

func caseInsensitiveRegex(s string) *regexp.Regexp {
	s = "(?i)" + s
	regexesLock.RLock()
	r, ok := regexes[s]
	regexesLock.RUnlock()
	if ok {
		return r
	}
	r, err := regexp.Compile("(?i)" + s)
	if err != nil {
		panic(fmt.Errorf("While compiling regex %s: %v", s, err))
	}
	regexesLock.Lock()
	regexes[s] = r
	regexesLock.Unlock()
	return r
}
