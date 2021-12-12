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

type allAboutBirdsResponse struct {
	CollectedAt     time.Time
	LifeHistoryURL  string
	LifeHistoryHTML string
	IdURL           string
	IdHTML          string
	OverviewURL     string
	OverviewHTML    string
}

func getAllAboutBirdsResponse(englishName string) (*allAboutBirdsResponse, error) {
	var aabr allAboutBirdsResponse
	r := &aabr

	missingKey := "all_about_birds/" + englishName
	if isKnownMissing(missingKey) {
		return r, missingError(missingKey)
	}

	memoizedFileName := "/tmp/bird_data_guessing/all_about_birds/" + englishName + ".xml"
	fileBytes, err := ioutil.ReadFile(memoizedFileName)
	if err == nil {
		err := xml.Unmarshal(fileBytes, r)
		return r, err
	}

	nameParam := strings.Replace(englishName, " ", "_", -1)
	url := "https://allaboutbirds.org/guide/" + nameParam

	r.IdURL = url + "/id"
	idPage, err := getDocumentFromUrl(r.IdURL, missingKey)
	if err != nil {
		return r, err
	}
	r.IdHTML, err = idPage.Html()
	if err != nil {
		return r, err
	}

	r.LifeHistoryURL = url + "/lifehistory"
	lifeHistoryPage, err := getDocumentFromUrl(r.LifeHistoryURL, missingKey)
	if err != nil {
		return r, err
	}
	r.LifeHistoryHTML, err = lifeHistoryPage.Html()
	if err != nil {
		return r, err
	}

	r.OverviewURL = url + "/overview"
	overviewPage, err := getDocumentFromUrl(r.OverviewURL, missingKey)
	if err != nil {
		return r, err
	}
	r.OverviewHTML, err = overviewPage.Html()
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

func (r *allAboutBirdsResponse) attribution() *attributions.Attribution {
	lifeHistoryPage := getDocumentFromString(r.LifeHistoryHTML)
	return &attributions.Attribution{
		OriginUrl:           r.IdURL,
		CollectedAt:         r.CollectedAt,
		OriginalTitle:       lifeHistoryPage.Find("title").First().Text(),
		Author:              "The Cornell Lab of Ornithology",
		AuthorUrl:           "https://www.birds.cornell.edu",
		LicenseUrl:          "https://www.birds.cornell.edu/home/terms-of-use",
		ScrapingMethodology: "github.com/gbdubs/bird_data_guessing",
		Context:             []string{"HTTP requested the bird information with URL guessing."},
	}
}

func (r *allAboutBirdsResponse) propertySearchers() *propertySearchers {
	overviewPage := getDocumentFromString(r.OverviewHTML)
	idPage := getDocumentFromString(r.IdHTML)
	lifeHistoryPage := getDocumentFromString(r.LifeHistoryHTML)

	idText := idPage.Find("main").First().Text()
	habitatText := lifeHistoryPage.Find("[aria-labelledby=habitat]").First().Text()
	foodText := lifeHistoryPage.Find("[aria-labelledby=food]").First().Text()
	nestingText := lifeHistoryPage.Find("[aria-labelledby=nesting]").First().Text()
	behaviorText := lifeHistoryPage.Find("[aria-labelledby=behavior]").First().Text()
	coolFactsText := overviewPage.Find("ul:contains('Cool Facts')").Text()
	return &propertySearchers{
		food:       searchIn(foodText),
		nestType:   searchIn(nestingText),
		habitat:    searchIn(habitatText),
		funFact:    searchIn(coolFactsText + behaviorText),
		wingspan:   searchIn(idText),
		clutchSize: searchIn(nestingText + behaviorText),
		eggColor:   searchIn(nestingText),
		all:        searchIn(habitatText + foodText + nestingText + behaviorText + coolFactsText + idText),
	}
}
