package bird_data_guessing

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gbdubs/attributions"
)

type whatBirdResponse struct {
	CollectedAt        time.Time
	OverviewURL        string
	OverviewHTML       string
	IdentificationURL  string
	IdentificationHTML string
	BehaviorURL        string
	BehaviorHTML       string
}

func getWhatBirdResponse(englishName string) (*whatBirdResponse, error) {
	wbr := whatBirdResponse{}
	r := &wbr

	missingKey := "what_bird/" + englishName
	if isKnownMissing(missingKey) {
		return r, missingError(missingKey)
	}

	memoizedFileName := "/tmp/bird_data_guessing/what_bird/" + englishName + ".xml"
	fileBytes, err := ioutil.ReadFile(memoizedFileName)
	if err == nil {
		err := xml.Unmarshal(fileBytes, r)
		return r, err
	}

	nameParam := strings.ToLower(
		strings.ReplaceAll(englishName, " ", "_"))
	id := whatBirdIdMap[nameParam]
	if id == 0 {
		nameParam = strings.ReplaceAll(nameParam, "'", "")
	}
	if id == 0 {
		markMissing(missingKey)
		return r, fmt.Errorf("404: bird name not in map %s", nameParam)
	}
	url := func(s string) string {
		return fmt.Sprintf("https://identify.whatbird.com/obj/%d/%s/%s.aspx", id, s, nameParam)
	}

	r.OverviewURL = url("overview")
	overviewPage, err := getDocumentFromUrl(r.OverviewURL, missingKey)
	if err != nil {
		return r, err
	}
	r.OverviewHTML, err = overviewPage.Html()
	if err != nil {
		return r, err
	}

	r.IdentificationURL = url("identification")
	identificationPage, err := getDocumentFromUrl(r.IdentificationURL, missingKey)
	if err != nil {
		return r, err
	}
	r.IdentificationHTML, err = identificationPage.Html()
	if err != nil {
		return r, err
	}

	r.BehaviorURL = url("behavior")
	behaviorPage, err := getDocumentFromUrl(r.BehaviorURL, missingKey)
	if err != nil {
		return r, err
	}
	r.BehaviorHTML, err = behaviorPage.Html()
	if err != nil {
		return r, err
	}

	r.CollectedAt = time.Now()

	asBytes, err := xml.MarshalIndent(r, "", "  ")
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

func (r *whatBirdResponse) attribution() *attributions.Attribution {
	overviewPage := getDocumentFromString(r.OverviewHTML)
	return &attributions.Attribution{
		OriginUrl:           r.OverviewURL,
		CollectedAt:         r.CollectedAt,
		OriginalTitle:       overviewPage.Find("title").First().Text(),
		Author:              "Mitch Waite Group",
		AuthorUrl:           "http://www.whatbird.com",
		License:             "Copyright 2002 - 2013, All Rights Reserved, Mitch Waite Group",
		ScrapingMethodology: "github.com/gbdubs/bird_data_guessing",
	}
}

func (r *whatBirdResponse) propertySearchers() *propertySearchers {
	overview := getDocumentFromString(r.OverviewHTML)
	identification := getDocumentFromString(r.IdentificationHTML)
	behavior := getDocumentFromString(r.BehaviorHTML)

	behaviorSearcher := func(s string) searcher {
		return searchIn(behavior.Find("h3:contains('" + s + "')").First().Next().Text())
	}

	return &propertySearchers{
		food:       behaviorSearcher("Forraging and Feeding"),
		nestType:   behaviorSearcher("Nest Location"),
		habitat:    behaviorSearcher("Range and Habitat"),
		clutchSize: searchIn(behavior.Find("#behavior").Text()),
		all: searchIn(
			behavior.Find("#behavior").Text() +
				identification.Find("#identification").Text() +
				overview.Find("#overview").Text()),
		funFact:  searchIn(overview.Find("h2:contains('INTERESTING FACTS')").First().Next().Text()),
		wingspan: searchIn(identification.Find("li:contains('Wingspan Range')").First().Text()),
		eggColor: searchIn(behavior.Find("#behavior").Text()),
	}
}
