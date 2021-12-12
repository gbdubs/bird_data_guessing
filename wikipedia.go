package bird_data_guessing

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gbdubs/attributions"
)

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
type wikipediaResponse struct {
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

func getWikipediaResponse(latinName string) (*wikipediaResponse, error) {
	var wr wikipediaResponse
	r := &wr
	missingKey := "wikipedia/" + latinName
	if isKnownMissing(missingKey) {
		return r, missingError(missingKey)
	}
	memoizedFileName := "/tmp/bird_data_guessing/wikipedia/" + latinName + ".xml"
	fileBytes, err := ioutil.ReadFile(memoizedFileName)
	if err == nil {
		err := xml.Unmarshal(fileBytes, r)
		if err != nil || !r.isMissingCirrusdoc() {
			return r, err
		}
	}
	req, err := http.NewRequest("GET", "https://en.wikipedia.org/w/api.php", nil)
	if err != nil {
		return r, err
	}
	q := req.URL.Query()
	q.Add("action", "query")
	q.Add("prop", "cirrusdoc|info")
	q.Add("format", "xml")
	q.Add("inprop", "url")
	q.Add("titles", latinName)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		markMissing(missingKey)
		return r, missingError(missingKey)
	}
	if resp.StatusCode != 200 {
		return r, errors.New(fmt.Sprintf("Request failed: %d %s", resp.StatusCode, resp.Status))
	}
	asBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	err = xml.Unmarshal(asBytes, &r)
	if err != nil {
		return r, err
	}
	if r.Query.Pages.Page.Idx == "-1" {
		return r, errors.New(fmt.Sprintf("404: Bird wasn't found: %s", req.URL))
	}
	r.CollectedAt = time.Now()
	asBytes, err = xml.MarshalIndent(r, "", " ")
	if err != nil {
		return r, err
	}
	err = os.MkdirAll(filepath.Dir(memoizedFileName), 0777)
	if err != nil {
		return r, err
	}
	err = ioutil.WriteFile(memoizedFileName, asBytes, 0777)
	return r, err
}

func (r *wikipediaResponse) isMissingCirrusdoc() bool {
	return reflect.DeepEqual(r.Query.Pages.Page.Cirrusdoc, cirrusdoc{})
}

func (r *wikipediaResponse) getText() string {
	s := r.Query.Pages.Page.Cirrusdoc.V.Source
	t := s.AttrText
	for _, at := range s.AuxiliaryText.V {
		t += " " + at
	}
	// Strips out Unicode spaces (ex: NBSP)
	return strings.Join(strings.Fields(t), " ")
}

func (r *wikipediaResponse) englishName() string {
	return r.Query.Pages.Page.Cirrusdoc.V.Source.Title
}

func (r *wikipediaResponse) propertySearchers() *propertySearchers {
	searcher := searchIn(r.getText())
	return &propertySearchers{
		food:       searcher,
		nestType:   searcher,
		habitat:    searcher,
		funFact:    searcher,
		eggColor:   searcher,
		wingspan:   searcher,
		clutchSize: searcher,
		all:        searcher,
	}
}

func (r *wikipediaResponse) attribution() *attributions.Attribution {
	sourceTimestamp := r.Query.Pages.Page.Cirrusdoc.V.Source.Timestamp
	if sourceTimestamp == "" {
		sourceTimestamp = r.Query.Pages.Page.Touched
	}
	createdAt, err := time.Parse(time.RFC3339, sourceTimestamp)
	if err != nil {
		fmt.Printf("Error when looking at element: %s\n", r.englishName())
		panic(err)
	}
	return &attributions.Attribution{
		OriginUrl:           r.Query.Pages.Page.Canonicalurl,
		CollectedAt:         r.CollectedAt,
		OriginalTitle:       r.Query.Pages.Page.Title,
		Author:              "Wikimedia Foundation",
		AuthorUrl:           "https://wikimedia.org",
		License:             "Creative Commons Attribution-ShareAlike 3.0 Unported License (CC BY-SA)",
		LicenseUrl:          "https://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License",
		ScrapingMethodology: "github.com/gbdubs/bird_data_guessing",
		Context:             []string{"Called Wikipedia's API with action=query, see api.go for details."},
		CreatedAt:           createdAt,
	}
}
