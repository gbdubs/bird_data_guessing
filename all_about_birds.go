package bird_data_guessing

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/gbdubs/bird"
	"github.com/gbdubs/sitemaps"

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
	allAboutBirdsMaxConcurrentRequests = 3
)

var allAboutBirdsSitemap *sitemaps.Sitemap = nil
var allAboutBirdsSitemapLock = sync.RWMutex{}

func aabSitemap() *sitemaps.Sitemap {
	if allAboutBirdsSitemap == nil {
		allAboutBirdsSitemapLock.Lock()
		if allAboutBirdsSitemap == nil {
			s, err := sitemaps.GetSitemapFromURL("https://www.allaboutbirds.org/guide/sitemap.xml")
			if err != nil {
				panic(err)
			}
			allAboutBirdsSitemap = s
		}
		allAboutBirdsSitemapLock.Unlock()
	}
	return allAboutBirdsSitemap
}

func createAllAboutBirdsRequests(name bird.BirdName) []*amass.GetRequest {
	sitemap := aabSitemap()
	nameParam := strings.ReplaceAll(name.English, " ", "_")
	requests := make([]*amass.GetRequest, 0)

	for _, page := range []string{allAboutBirdsOverviewSuffix, allAboutBirdsIdSuffix, allAboutBirdsLifeHistorySuffix} {
		url := fmt.Sprintf("https://allaboutbirds.org/guide/%s/%s", nameParam, page)
		actualUrl, levDist := sitemap.BestFuzzyMatch(url)
		if levDist > 2 {
			recordMissing(allAboutBirdsSite, name)
			continue
		}
		requestKeyPrefix := regexp.MustCompile("allaboutbirds.org/guide/(.+)/.+/$").FindStringSubmatch(actualUrl)[1]
		request := &amass.GetRequest{
			Site:                      allAboutBirdsSite,
			RequestKey:                requestKeyPrefix + "_" + page,
			URL:                       actualUrl,
			SiteMaxConcurrentRequests: allAboutBirdsMaxConcurrentRequests,
			Attribution: attributions.Attribution{
				CreatedAt:           sitemap.LastUpdated[actualUrl],
				Author:              "The Cornell Lab of Ornithology",
				AuthorUrl:           "https://www.birds.cornell.edu",
				LicenseUrl:          "https://www.birds.cornell.edu/home/terms-of-use",
				ScrapingMethodology: "github.com/gbdubs/bird_data_guessing/all_about_birds",
			},
		}
		request.SetRoundTripData(name)
		requests = append(requests, request)
	}
	return requests
}

// TODO Consider how missing works here with search result pages!!!

func reconstructAllAboutBirdsResponsesKeyedByEnglishName(responses []*amass.GetResponse) map[string]*allAboutBirdsResponse {
	m := make(map[string]map[string]*amass.GetResponse)
	m[allAboutBirdsOverviewSuffix] = make(map[string]*amass.GetResponse)
	m[allAboutBirdsIdSuffix] = make(map[string]*amass.GetResponse)
	m[allAboutBirdsLifeHistorySuffix] = make(map[string]*amass.GetResponse)
	englishNames := make(map[string]bool)
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
		birdName := &bird.BirdName{}
		response.GetRoundTripData(birdName)
		englishName := birdName.English
		englishNames[englishName] = true
		m[page][englishName] = response
	}
	result := make(map[string]*allAboutBirdsResponse)
	for englishName, _ := range englishNames {
		result[englishName] = &allAboutBirdsResponse{
			Identification: *m[allAboutBirdsIdSuffix][englishName],
			LifeHistory:    *m[allAboutBirdsLifeHistorySuffix][englishName],
			Overview:       *m[allAboutBirdsOverviewSuffix][englishName],
		}
	}
	return result
}

func allAboutBirdsRequestForTesting(englishName string) *allAboutBirdsResponse {
	bn := bird.BirdName{English: englishName}
	rs := createAllAboutBirdsRequests(bn)
	if len(rs) != 3 {
		panic(fmt.Errorf("Expected 3 all about birds requests, was %d, for key %s.", len(rs), englishName))
	}
	resps, err := amass.AmasserForTests().GetAll(rs)
	if err != nil {
		panic(fmt.Errorf("GetAll request failed for %s: %v", englishName, err))
	}
	m := reconstructAllAboutBirdsResponsesKeyedByEnglishName(resps)
	result, ok := m[englishName]
	if !ok {
		panic(fmt.Errorf("Expected key %s in map, but map was %#v.", englishName, m))
	}
	return result
}

func (r *allAboutBirdsResponse) propertySearchers() *propertySearchers {
	overviewPage := r.Overview.AsDocument()
	idPage := r.Identification.AsDocument()
	lifeHistoryPage := r.LifeHistory.AsDocument()

	wingspanText := idPage.Find("h5:contains('measurements')").Next().Text()
	if !strings.Contains(wingspanText, "Wingspan") {
		wingspanText = ""
	}
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
		funFact:    attributedSearch(&r.Overview.Attribution, coolFactsText, behaviorText),

		food:     attributedSearch(&r.LifeHistory.Attribution, foodText),
		nestType: attributedSearch(&r.LifeHistory.Attribution, nestingText),
		habitat:  attributedSearch(&r.LifeHistory.Attribution, habitatText),

		predator: attributedSearch(&r.Overview.Attribution, foodText, behaviorText, coolFactsText, idText),
		flocking: attributedSearch(&r.Overview.Attribution, behaviorText, coolFactsText, idText, nestingText),
	}
}
