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

// Delete this
func (s *searcher) CountMatches(searchFor ...string) *Property {
	if s.isDefault() {
		return &Property{}
	}
	p := &Property{}
	for _, r := range searchFor {
		regex := caseInsensitiveRegex(r)
		n := len(regex.FindAllStringIndex(s.text, -1))
		p.Strength += n
		p.Context += fmt.Sprintf("matches:%s=%d ", r, n)
	}
	return p
}

const defaultContextWidth = 30

func (s *searcher) ZZZCountMatches(searchFor ...string) *inference.Int {
	components := make([]*inference.Int, 0)
	for key, matches := range s.ZZZGetMatches(searchFor...) {
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

func (s *searcher) ZZZGetMatches(searchFor ...string) map[string][]*inference.String {
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

// Delete this
func (s *searcher) ExtractMatch(pattern string, captureGroup int) *Property {
	if s.isDefault() {
		return &Property{}
	}
	matches := caseInsensitiveRegex(pattern).FindAllStringSubmatch(s.text, -1)
	if len(matches) == 0 {
		return &Property{
			Context: fmt.Sprintf("Pattern %s had 0 matches", pattern),
		}
	}
	match := matches[rand.Intn(len(matches))]
	return &Property{
		StringValue: match[captureGroup],
		Context:     fmt.Sprintf("Pattern %s had %d matches: %+v, selected %s, capturing %s", pattern, len(matches), matches, match[0], match[captureGroup]),
	}
}

func (s *searcher) ZZZExtractAllMatch(pattern string, captureGroup int) []*inference.String {
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

func (s *searcher) ZZZExtractMatch(pattern string, captureGroup int) *inference.String {
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

// Delete this
func (s *searcher) ExtractAnyMatch(patternToCaptureGroup map[string]int) *Property {
	if s.isDefault() {
		return &Property{}
	}
	for p, cg := range patternToCaptureGroup {
		m := s.ExtractMatch(p, cg)
		if m.StringValue != "" {
			return m
		}
	}
	return &Property{}
}

func (s *searcher) ZZZExtractAnyMatch(patternToCaptureGroup map[string]int) *inference.String {
	for p, cg := range patternToCaptureGroup {
		m := s.ZZZExtractMatch(p, cg)
		if *m != (inference.String{}) {
			return m
		}
	}
	return &inference.String{}
}

func (s *searcher) ZZZExtractAllMatches(patternToCaptureGroup map[string]int) []*inference.String {
	result := make([]*inference.String, 0)
	for p, cg := range patternToCaptureGroup {
		result = append(result, s.ZZZExtractAllMatch(p, cg)...)
	}
	return result
}

// Delete this
/*
func (s searcher) CountProximatePairs(searchForA string, searchForB string, allowedDistance int) (int, string) {
	if s.isDefault() {
		return 0, ""
	}
	a := caseInsensitiveRegex(searchForA)
	b := caseInsensitiveRegex(searchForB)
	matchesA := a.FindAllStringIndex(s.text, -1)
	matchesB := b.FindAllStringIndex(s.text, -1)
	lenA := len(matchesA)
	lenB := len(matchesB)
	hits := 0
	i := 0
	j := 0
	for i < lenA && j < lenB {
		idxA := i
		if idxA == lenA {
			idxA = lenA - 1
		}
		idxB := j
		if idxB == lenB {
			idxB = lenB - 1
		}
		valA := matchesA[idxA][0]
		valB := matchesB[idxB][0]
		if valB > valA {
			valA += len(searchForA)
			if j < lenB {
				j++
			} else {
				i++
			}
		} else {
			valB += len(searchForB)
			if i < lenA {
				i++
			} else {
				j++
			}
		}
		if abs(valA-valB) <= allowedDistance {
			hits++
		}
	}
	return hits, fmt.Sprintf("proximate-pairs:%s:%s=%d", searchForA, searchForB, hits)
}
*/
func (s *searcher) isDefault() bool {
	return s.text == ""
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
		panic(err)
	}
	regexesLock.Lock()
	regexes[s] = r
	regexesLock.Unlock()
	return r
}
