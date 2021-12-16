package bird_data_guessing

import (
	"fmt"
	"strings"

	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
)

type allAboutBirdsResponse struct {
	LifeHistory    amass.GetResponse
	Identification amass.GetResponse
	Overview       amass.GetResponse
}

const (
	allAboutBirdsSite                  = "all_about_birds"
	allAboutBirdsIdSuffix              = "id"
	allAboutBirdsLifeHistorySuffix     = "lifehistory"
	allAboutBirdsOverviewSuffix        = "overview"
	allAboutBirdsMaxConcurrentRequests = 2
)

func createAllAboutBirdsRequests(name BirdName) []*amass.GetRequest {
	nameParam := strings.ReplaceAll(name.EnglishName, " ", "_")
	requestKeyPrefix := strings.ToLower(nameParam)
	makeReq := func(page string) *amass.GetRequest {
		result := &amass.GetRequest{
			Site:                      allAboutBirdsSite,
			RequestKey:                requestKeyPrefix + "_" + page,
			URL:                       fmt.Sprintf("https://allaboutbirds.org/guide/%s/%s", nameParam, page),
			SiteMaxConcurrentRequests: allAboutBirdsMaxConcurrentRequests,
			Attribution: attributions.Attribution{
				Author:              "The Cornell Lab of Ornithology",
				AuthorUrl:           "https://www.birds.cornell.edu",
				LicenseUrl:          "https://www.birds.cornell.edu/home/terms-of-use",
				ScrapingMethodology: "github.com/gbdubs/bird_data_guessing/all_about_birds",
			},
		}
		result.SetRoundTripData(name)
		return result
	}
	return []*amass.GetRequest{
		makeReq(allAboutBirdsOverviewSuffix),
		makeReq(allAboutBirdsIdSuffix),
		makeReq(allAboutBirdsLifeHistorySuffix),
	}
}

func reconstructAllAboutBirdsResponsesKeyedByLatinName(responses []*amass.GetResponse) map[string]*allAboutBirdsResponse {
	m := make(map[string]map[string]*amass.GetResponse)
	m[allAboutBirdsOverviewSuffix] = make(map[string]*amass.GetResponse)
	m[allAboutBirdsIdSuffix] = make(map[string]*amass.GetResponse)
	m[allAboutBirdsLifeHistorySuffix] = make(map[string]*amass.GetResponse)
	latinNames := make(map[string]bool)
	for _, response := range responses {
		if response.Site != allAboutBirdsSite {
			continue
		}
		page := ""
		if strings.HasSuffix(response.RequestKey, allAboutBirdsIdSuffix) {
			page = allAboutBirdsIdSuffix
		} else if strings.HasSuffix(response.RequestKey, allAboutBirdsOverviewSuffix) {
			page = allAboutBirdsOverviewSuffix
		} else if strings.HasSuffix(response.RequestKey, allAboutBirdsLifeHistorySuffix) {
			page = allAboutBirdsLifeHistorySuffix
		} else {
			panic(fmt.Errorf("Unrecongnized response request key %s for all about birds.", response.RequestKey))
		}
		birdName := &BirdName{}
		response.GetRoundTripData(birdName)
		latinName := birdName.LatinName
		latinNames[latinName] = true
		m[page][latinName] = response
	}
	result := make(map[string]*allAboutBirdsResponse)
	for latinName, _ := range latinNames {
		result[latinName] = &allAboutBirdsResponse{
			Identification: *m[allAboutBirdsIdSuffix][latinName],
			LifeHistory:    *m[allAboutBirdsLifeHistorySuffix][latinName],
			Overview:       *m[allAboutBirdsOverviewSuffix][latinName],
		}
	}
	return result
}

func (r *allAboutBirdsResponse) propertySearchers() *propertySearchers {
	overviewPage := r.Overview.AsDocument()
	idPage := r.Identification.AsDocument()
	lifeHistoryPage := r.LifeHistory.AsDocument()

	wingspanText := idPage.Find("h5:contains('measurements')").Next().Text()
	idText := idPage.Find("main").First().Text()
	habitatText := lifeHistoryPage.Find("[aria-labelledby=habitat]").First().Text()
	foodText := lifeHistoryPage.Find("[aria-labelledby=food]").First().Text()
	nestingText := lifeHistoryPage.Find("[aria-labelledby=nesting]").First().Text()
	behaviorText := lifeHistoryPage.Find("[aria-labelledby=behavior]").First().Text()
	coolFactsText := overviewPage.Find("ul:contains('Cool Facts')").Text()
	return &propertySearchers{
		wingspan:   attributedSearch(&r.Identification.Attribution, wingspanText),
		clutchSize: attributedSearch(&r.LifeHistory.Attribution, nestingText),
		eggColor:   attributedSearch(&r.LifeHistory.Attribution, nestingText),
		funFact:    attributedSearch(&r.Overview.Attribution, coolFactsText+behaviorText),

		food:     attributedSearch(&r.LifeHistory.Attribution, foodText),
		nestType: attributedSearch(&r.LifeHistory.Attribution, nestingText),
		habitat:  attributedSearch(&r.LifeHistory.Attribution, habitatText),

		predator: attributedSearch(&r.Overview.Attribution, foodText+behaviorText+coolFactsText+idText),
		flocking: attributedSearch(&r.Overview.Attribution, behaviorText+coolFactsText+idText+nestingText),
	}
}