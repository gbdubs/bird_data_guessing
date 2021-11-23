package bird_data_guessing

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

// Generated from https://www.onlinetool.io/xmltogo/
type wikipediaResponse struct {
	XMLName       xml.Name `xml:"api"`
	Text          string   `xml:",chardata"`
	Batchcomplete string   `xml:"batchcomplete,attr"`
	Query         struct {
		Text  string `xml:",chardata"`
		Pages struct {
			Text string `xml:",chardata"`
			Page struct {
				Text      string `xml:",chardata"`
				Idx       string `xml:"_idx,attr"`
				Pageid    string `xml:"pageid,attr"`
				Ns        string `xml:"ns,attr"`
				Title     string `xml:"title,attr"`
				Cirrusdoc struct {
					Text string `xml:",chardata"`
					V    struct {
						Text    string `xml:",chardata"`
						Index   string `xml:"index,attr"`
						Type    string `xml:"type,attr"`
						ID      string `xml:"id,attr"`
						Version string `xml:"version"`
						Source  struct {
							Text            string `xml:",chardata"`
							ContentModel    string `xml:"content_model,attr"`
							OpeningText     string `xml:"opening_text,attr"`
							Wiki            string `xml:"wiki,attr"`
							Language        string `xml:"language,attr"`
							Title           string `xml:"title,attr"`
							AttrText        string `xml:"text,attr"`
							Defaultsort     string `xml:"defaultsort,attr"`
							Timestamp       string `xml:"timestamp,attr"`
							WikibaseItem    string `xml:"wikibase_item,attr"`
							SourceText      string `xml:"source_text,attr"`
							VersionType     string `xml:"version_type,attr"`
							Version         string `xml:"version,attr"`
							NamespaceText   string `xml:"namespace_text,attr"`
							Namespace       string `xml:"namespace,attr"`
							TextBytes       string `xml:"text_bytes,attr"`
							IncomingLinks   string `xml:"incoming_links,attr"`
							PopularityScore string `xml:"popularity_score,attr"`
							CreateTimestamp string `xml:"create_timestamp,attr"`
							Template        struct {
								Text string   `xml:",chardata"`
								V    []string `xml:"_v"`
							} `xml:"template"`
							AuxiliaryText struct {
								Text string   `xml:",chardata"`
								V    []string `xml:"_v"`
							} `xml:"auxiliary_text"`
							Redirect struct {
								Text string `xml:",chardata"`
								V    []struct {
									Text      string `xml:",chardata"`
									Namespace string `xml:"namespace,attr"`
									Title     string `xml:"title,attr"`
								} `xml:"_v"`
							} `xml:"redirect"`
							Heading struct {
								Text string   `xml:",chardata"`
								V    []string `xml:"_v"`
							} `xml:"heading"`
							Coordinates  string `xml:"coordinates"`
							ExternalLink struct {
								Text string   `xml:",chardata"`
								V    []string `xml:"_v"`
							} `xml:"external_link"`
							Category struct {
								Text string   `xml:",chardata"`
								V    []string `xml:"_v"`
							} `xml:"category"`
							OutgoingLink struct {
								Text string   `xml:",chardata"`
								V    []string `xml:"_v"`
							} `xml:"outgoing_link"`
							OresArticletopics struct {
								Text string   `xml:",chardata"`
								V    []string `xml:"_v"`
							} `xml:"ores_articletopics"`
							OresArticletopic struct {
								Text string   `xml:",chardata"`
								V    []string `xml:"_v"`
							} `xml:"ores_articletopic"`
							WeightedTags struct {
								Text string   `xml:",chardata"`
								V    []string `xml:"_v"`
							} `xml:"weighted_tags"`
						} `xml:"source"`
					} `xml:"_v"`
				} `xml:"cirrusdoc"`
			} `xml:"page"`
		} `xml:"pages"`
	} `xml:"query"`
}

func getWikipediaPage(latinName string) (*wikipediaResponse, error) {
	var wr wikipediaResponse
	r := &wr
	memoizedFileName := "/tmp/bird_property_guessing/" + latinName + ".xml"
	fileBytes, err := ioutil.ReadFile(memoizedFileName)
	if err == nil {
		err := xml.Unmarshal(fileBytes, r)
		return r, err
	}
	req, err := http.NewRequest("GET", "https://en.wikipedia.org/w/api.php", nil)
	if err != nil {
		return r, err
	}
	q := req.URL.Query()
	q.Add("action", "query")
	q.Add("prop", "cirrusdoc")
	q.Add("format", "xml")
	q.Add("titles", latinName)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return r, errors.New(fmt.Sprintf("Request failed: %d %s", resp.StatusCode, resp.Status))
	}
	asBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	err = os.MkdirAll(filepath.Dir(memoizedFileName), 0777)
	if err != nil {
		return r, err
	}
	err = ioutil.WriteFile(memoizedFileName, asBytes, 0777)
	if err != nil {
		return r, err
	}
	err = xml.Unmarshal(asBytes, &r)
	return r, err
}

func (r *wikipediaResponse) GetText() string {
	return r.Query.Pages.Page.Cirrusdoc.V.Source.AttrText
}

func (r *wikipediaResponse) setText(t string) {
	r.Query.Pages.Page.Cirrusdoc.V.Source.AttrText = t
}

func (r *wikipediaResponse) CountMatches(searchFor []string) *Property {
	p := &Property{}
	for _, s := range searchFor {
		regex := caseInsensitiveRegex(s)
		n := len(regex.FindAllStringIndex(r.GetText(), -1))
		p.Strength += n
		p.Context += fmt.Sprintf("matches:%s=%d ", s, n)
	}
	return p
}

func (r *wikipediaResponse) getMatches(pattern string, captureGroup int) []string {
	matches := caseInsensitiveRegex(pattern).FindAllStringSubmatch(r.GetText(), -1)
	var result []string
	for _, match := range matches {
		result = append(result, match[captureGroup])
	}
	return result
}

func (r *wikipediaResponse) ExtractMatch(pattern string, captureGroup int) *Property {
	matches := r.getMatches(pattern, captureGroup)
	if len(matches) == 0 {
		return &Property{
			Context: fmt.Sprintf("Pattern %s had 0 matches", pattern),
		}
	}
	match := matches[rand.Intn(len(matches))]
	return &Property{
		StringValue: match,
		Strength:    len(matches),
		Context:     fmt.Sprintf("Pattern %s had %d matches: %+v, selected %s", pattern, len(matches), matches, match),
	}
}

func (r *wikipediaResponse) ExtractAnyMatch(patternToCaptureGroup map[string]int) *Property {
	for p, cg := range patternToCaptureGroup {
		m := r.ExtractMatch(p, cg)
		if m.StringValue != "" {
			return m
		}
	}
	return &Property{}
}

func (r *wikipediaResponse) CountProximatePairs(searchForA string, searchForB string, allowedDistance int) (int, string) {
	a := caseInsensitiveRegex(searchForA)
	b := caseInsensitiveRegex(searchForB)
	matchesA := a.FindAllStringIndex(r.GetText(), -1)
	matchesB := b.FindAllStringIndex(r.GetText(), -1)
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func caseInsensitiveRegex(s string) *regexp.Regexp {
	r, err := regexp.Compile("(?i)" + s)
	if err != nil {
		panic(err)
	}
	return r
}
