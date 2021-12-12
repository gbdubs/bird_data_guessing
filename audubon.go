package bird_data_guessing

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gbdubs/attributions"
)

type audubonResponse struct {
	URL         string
	HTML        string
	CollectedAt time.Time
}

func getAudubonResponse(englishName string) (*audubonResponse, error) {
	ar := audubonResponse{}
	r := &ar

	missingKey := "audubon/" + englishName
	if isKnownMissing(missingKey) {
		return r, missingError(missingKey)
	}

	memoizedFileName := "/tmp/bird_data_guessing/audubon/" + englishName + ".xml"
	fileBytes, err := ioutil.ReadFile(memoizedFileName)
	if err == nil {
		err := xml.Unmarshal(fileBytes, r)
		return r, err
	}

	nameParam := strings.ReplaceAll(strings.ReplaceAll(englishName, " ", "-"), "'", "")
	url := "https://audubon.org/field-guide/bird/" + nameParam

	r.URL = url
	page, err := getDocumentFromUrl(r.URL, missingKey)
	if err != nil {
		return r, err
	}
	r.HTML, err = page.Html()
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

func (r *audubonResponse) attribution() *attributions.Attribution {
	page := getDocumentFromString(r.HTML)
	return &attributions.Attribution{
		OriginUrl:           r.URL,
		CollectedAt:         r.CollectedAt,
		OriginalTitle:       page.Find("title").First().Text(),
		Author:              "National Audubon Society, Inc.",
		AuthorUrl:           "https://audubon.org",
		License:             "All rights reserved",
		LicenseUrl:          "https://www.audubon.org/terms-use",
		ScrapingMethodology: "github.com/gbdubs/bird_data_guessing",
		Context:             []string{"HTTP requested the bird information via URL guessing."},
	}
}

func (r *audubonResponse) propertySearchers() *propertySearchers {
	page := getDocumentFromString(r.HTML)

	dietText := page.Find("h2:contains('Diet')").First().Next().Text()
	feedingText := page.Find("h2:contains('Feeding')").First().Next().Text()
	eggsText := "EggsAudubonEggs " + page.Find("h2:contains('Eggs')").First().Next().Text()
	nestingText := page.Find("h2:contains('Nesting')").First().Next().Text()
	habitatText := page.Find("th:contains('Habitat')").First().Parent().Find("td").First().Text()
	allText := page.Find("body").Text()

	return &propertySearchers{
		food:       searchIn(dietText + feedingText),
		nestType:   searchIn(nestingText),
		habitat:    searchIn(habitatText + nestingText),
		clutchSize: searchIn(eggsText + nestingText),
		eggColor:   searchIn(eggsText),
		all:        searchIn(allText),
	}
}
