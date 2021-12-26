package bird_data_guessing

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
	"github.com/gbdubs/bird"
)

type wikipediaResponse struct {
	Response           amass.GetResponse
	decodedApiResponse wikipediaApiResponse
	hasBeenDecoded     bool
}

const (
	wikipediaSite                  = "wikipedia"
	maxWikipediaConcurrentRequests = 2
)

func createWikipediaRequests(name bird.BirdName) []*amass.GetRequest {
	result := make([]*amass.GetRequest, 0)
	if isMissing(wikipediaSite, name) {
		return result
	}
	v := url.Values{}
	v.Add("action", "query")
	v.Add("prop", "cirrusdoc|info")
	v.Add("format", "xml")
	v.Add("inprop", "url")
	v.Add("redirects", "1")
	v.Add("titles", name.Latin)
	url := "https://en.wikipedia.org/w/api.php?" + v.Encode()
	r := &amass.GetRequest{
		Site:                      wikipediaSite,
		RequestKey:                name.Latin,
		URL:                       url,
		SiteMaxConcurrentRequests: maxWikipediaConcurrentRequests,
		Attribution: attributions.Attribution{
			Author:              "National Wikipedia Society, Inc.",
			AuthorUrl:           "https://wikipedia.org",
			License:             "Creative Commons Attribution-ShareAlike 3.0 Unported License (CC BY-SA)",
			LicenseUrl:          "https://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License",
			ScrapingMethodology: "github.com/gbdubs/bird_data_guessing/wikipedia",
		},
	}
	r.SetRoundTripData(name)
	result = append(result, r)
	return result
}

func reconstructWikipediaResponsesKeyedByLatinName(responses []*amass.GetResponse) map[string]*wikipediaResponse {
	result := make(map[string]*wikipediaResponse)
	for _, response := range responses {
		if response.Site != wikipediaSite {
			continue
		}
		birdName := &bird.BirdName{}
		response.GetRoundTripData(birdName)
		wr := &wikipediaResponse{
			Response: *response,
		}
		if wr.isWikipediaResponseMissing() {
			recordMissing(wikipediaSite, *birdName)
			continue
		}
		wr.tweakResponse()
		result[birdName.Latin] = wr
	}
	return result
}

func wikipediaRequestForTesting(latinName string) *wikipediaResponse {
	bn := bird.BirdName{Latin: latinName}
	rs := createWikipediaRequests(bn)
	if len(rs) != 1 {
		panic(fmt.Errorf("Expected 1 wikipedia request, was %d, for key %s.", len(rs), latinName))
	}
	resp, err := rs[0].Get()
	if err != nil {
		panic(fmt.Errorf("Get request failed for %s: %v", latinName, err))
	}
	m := reconstructWikipediaResponsesKeyedByLatinName([]*amass.GetResponse{resp})
	result, ok := m[latinName]
	if !ok {
		panic(fmt.Errorf("Expected key %s in map, but map was %+v.", latinName, m))
	}
	return result
}

func (r *wikipediaResponse) isWikipediaResponseMissing() bool {
	return r.Response.StatusCode == 404 || r.isMissingCirrusdoc()
}

func (r *wikipediaResponse) tweakResponse() {
	sourceTimestamp := r.resp().Query.Pages.Page.Cirrusdoc.V.Source.Timestamp
	if sourceTimestamp == "" {
		sourceTimestamp = r.resp().Query.Pages.Page.Touched
	}
	if sourceTimestamp != "" {
		createdAt, err := time.Parse(time.RFC3339, sourceTimestamp)
		if err != nil {
			panic(fmt.Errorf("Error when looking at element %s: %v", r.englishName(), err))
		}
		r.Response.Attribution.CreatedAt = createdAt
	}
	r.Response.Attribution.OriginUrl = r.resp().Query.Pages.Page.Canonicalurl
	r.Response.Attribution.OriginalTitle = r.resp().Query.Pages.Page.Title
}

func (r *wikipediaResponse) englishName() string {
	return r.resp().Query.Pages.Page.Cirrusdoc.V.Source.Title
}

func (r *wikipediaResponse) resp() wikipediaApiResponse {
	if !r.hasBeenDecoded {
		err := r.Response.AsXMLObject(&r.decodedApiResponse)
		if err != nil {
			panic(fmt.Errorf("Couldn't decode XML for wikipedia response: %v", err))
		}
		r.hasBeenDecoded = true
	}
	return r.decodedApiResponse
}

func (r *wikipediaResponse) propertySearchers() *propertySearchers {
	s := r.resp().Query.Pages.Page.Cirrusdoc.V.Source
	t := s.AttrText
	for _, at := range s.AuxiliaryText.V {
		t += " " + at
	}
	// Strips out Unicode spaces (ex: NBSP)
	t = strings.Join(strings.Fields(t), " ")
	searcher := attributedSearch(&r.Response.Attribution, t)
	return &propertySearchers{
		food:       searcher,
		nestType:   searcher,
		habitat:    searcher,
		funFact:    searcher,
		eggColor:   searcher,
		wingspan:   searcher,
		clutchSize: searcher,
		predator:   searcher,
		flocking:   searcher,
	}
}

func (r *wikipediaResponse) isMissingCirrusdoc() bool {
	return reflect.DeepEqual(r.resp().Query.Pages.Page.Cirrusdoc, cirrusdoc{})
}

type cirrusdoc struct {
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
}

// Generated from https://www.onlinetool.io/xmltogo/
type wikipediaApiResponse struct {
	XMLName       xml.Name `xml:"api"`
	Text          string   `xml:",chardata"`
	Batchcomplete string   `xml:"batchcomplete,attr"`
	Query         struct {
		Text  string `xml:",chardata"`
		Pages struct {
			Text string `xml:",chardata"`
			Page struct {
				Text                 string    `xml:",chardata"`
				Idx                  string    `xml:"_idx,attr"`
				Pageid               string    `xml:"pageid,attr"`
				Ns                   string    `xml:"ns,attr"`
				Title                string    `xml:"title,attr"`
				Contentmodel         string    `xml:"contentmodel,attr"`
				Pagelanguage         string    `xml:"pagelanguage,attr"`
				Pagelanguagehtmlcode string    `xml:"pagelanguagehtmlcode,attr"`
				Pagelanguagedir      string    `xml:"pagelanguagedir,attr"`
				Touched              string    `xml:"touched,attr"`
				Lastrevid            string    `xml:"lastrevid,attr"`
				Length               string    `xml:"length,attr"`
				Redirect             string    `xml:"redirect,attr"`
				Fullurl              string    `xml:"fullurl,attr"`
				Editurl              string    `xml:"editurl,attr"`
				Canonicalurl         string    `xml:"canonicalurl,attr"`
				Cirrusdoc            cirrusdoc `xml:"cirrusdoc"`
			} `xml:"page"`
		} `xml:"pages"`
	} `xml:"query"`
	CollectedAt time.Time
	NotFound    bool
}
