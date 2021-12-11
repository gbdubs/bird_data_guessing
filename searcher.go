package bird_data_guessing

import (
	"fmt"
	"math/rand"
	"regexp"
	"sync"
)

type searcher struct {
	text string
}

func (s *searcher) CountMatches(searchFor []string) *Property {
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
