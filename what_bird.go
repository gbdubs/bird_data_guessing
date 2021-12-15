package bird_data_guessing

import (
	"fmt"
	"strings"

	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
)

type whatBirdResponse struct {
	Overview       amass.GetResponse
	Identification amass.GetResponse
	Behavior       amass.GetResponse
}

const (
	whatBirdSite                  = "what_bird"
	whatBirdIdPage                = "id"
	whatBirdBehaviorPage          = "lifehistory"
	whatBirdOverviewPage          = "overview"
	whatBirdMaxConcurrentRequests = 2
)

func createWhatBirdsRequests(englishName string) []*amass.GetRequest {
	nameParam := strings.ToLower(strings.ReplaceAll(englishName, " ", "_"))
	whatBirdId := whatBirdIdMap[nameParam]
	if whatBirdId == 0 {
		nameParam = strings.ReplaceAll(nameParam, "'", "")
	}
	if whatBirdId == 0 {
		// TODO log this and figure out other naming modifications
		return []*amass.GetRequest{}
	}
	makeReq := func(page string) *amass.GetRequest {
		url := fmt.Sprintf(
			"https://identify.whatbird.com/obj/%d/%s/%s.aspx",
			whatBirdId, page, nameParam)
		return &amass.GetRequest{
			Site:                      whatBirdSite,
			RequestKey:                nameParam + "_" + page,
			URL:                       url,
			SiteMaxConcurrentRequests: whatBirdMaxConcurrentRequests,
			Attribution: attributions.Attribution{
				Author:              "Mitch Waite Group",
				AuthorUrl:           "http://www.whatbird.com",
				License:             "Copyright 2002 - 2013, All Rights Reserved, Mitch Waite Group",
				ScrapingMethodology: "github.com/gbdubs/bird_data_guessing/what_bird",
			},
		}
	}
	return []*amass.GetRequest{
		makeReq(whatBirdOverviewPage),
		makeReq(whatBirdIdPage),
		makeReq(whatBirdBehaviorPage),
	}
}

func reconstructWhatBirdsResponses(responses []*amass.GetResponse) []*whatBirdResponse {
	m := make(map[string]map[string]*amass.GetResponse)
	m[whatBirdOverviewPage] = make(map[string]*amass.GetResponse)
	m[whatBirdIdPage] = make(map[string]*amass.GetResponse)
	m[whatBirdBehaviorPage] = make(map[string]*amass.GetResponse)
	birdKeys := make(map[string]bool)
	for _, response := range responses {
		if response.Site != whatBirdSite {
			continue
		}
		page := ""
		if strings.HasSuffix(response.RequestKey, whatBirdIdPage) {
			page = whatBirdIdPage
		} else if strings.HasSuffix(response.RequestKey, whatBirdOverviewPage) {
			page = whatBirdOverviewPage
		} else if strings.HasSuffix(response.RequestKey, whatBirdBehaviorPage) {
			page = whatBirdBehaviorPage
		} else {
			panic(fmt.Errorf("Unrecongnized response request key %s for what bird.", response.RequestKey))
		}
		birdKey := strings.ReplaceAll(response.RequestKey, "_"+page, "")
		birdKeys[birdKey] = true
		m[page][birdKey] = response
	}
	result := make([]*whatBirdResponse, len(birdKeys))
	i := 0
	for birdKey, _ := range birdKeys {
		result[i] = &whatBirdResponse{
			Identification: *m[whatBirdIdPage][birdKey],
			Behavior:       *m[whatBirdBehaviorPage][birdKey],
			Overview:       *m[whatBirdOverviewPage][birdKey],
		}
		i++
	}
	return result
}

func (r *whatBirdResponse) propertySearchers() *propertySearchers {
	overview := r.Overview.AsDocument()
	identification := r.Identification.AsDocument()
	behavior := r.Behavior.AsDocument()

	wingspanText := identification.Find("li:contains('Wingspan Range')").First().Text()
	behaviorHeadingSearcher := func(s string) searcher {
		return attributedSearch(
			&r.Behavior.Attribution,
			behavior.Find("h3:contains('"+s+"')").First().Next().Text())
	}
	behaviorText := behavior.Find("#behavior").Text()
	behaviorSearcher := attributedSearch(&r.Behavior.Attribution, behaviorText)
	behaviorOverviewText := behaviorText + overview.Find("#overview").Text()
	behaviorOverviewSearcher := attributedSearch(&r.Overview.Attribution, behaviorOverviewText)

	funFactText := overview.Find("h2:contains('INTERESTING FACTS')").First().Next().Text()

	return &propertySearchers{
		wingspan:   attributedSearch(&r.Identification.Attribution, wingspanText),
		clutchSize: behaviorSearcher,
		eggColor:   behaviorSearcher,
		funFact:    attributedSearch(&r.Overview.Attribution, funFactText),

		food:     behaviorHeadingSearcher("Forraging and Feeding"),
		nestType: behaviorHeadingSearcher("Nest Location"),
		habitat:  behaviorHeadingSearcher("Range and Habitat"),

		predator: behaviorOverviewSearcher,
		flocking: behaviorOverviewSearcher,
	}
}
