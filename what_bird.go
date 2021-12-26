package bird_data_guessing

import (
	"fmt"
	"strings"

	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
	"github.com/gbdubs/bird"
)

type whatBirdResponse struct {
	Overview       amass.GetResponse
	Identification amass.GetResponse
	Behavior       amass.GetResponse
}

const (
	whatBirdSite                  = "what_bird"
	whatBirdIdentificationPage    = "identification"
	whatBirdBehaviorPage          = "behavior"
	whatBirdOverviewPage          = "overview"
	whatBirdMaxConcurrentRequests = 2
)

func createWhatBirdRequests(birdName bird.BirdName) []*amass.GetRequest {
	nameParam := strings.ToLower(strings.ReplaceAll(birdName.English, " ", "_"))
	whatBirdId := whatBirdIdMap[nameParam]
	if whatBirdId == 0 {
		nameParam = strings.ReplaceAll(nameParam, "'", "")
		whatBirdId = whatBirdIdMap[nameParam]
	}
	if whatBirdId == 0 {
		recordMissing(whatBirdSite, birdName)
		return []*amass.GetRequest{}
	}
	makeReq := func(page string) *amass.GetRequest {
		url := fmt.Sprintf(
			"https://identify.whatbird.com/obj/%d/%s/%s.aspx",
			whatBirdId, page, nameParam)
		req := &amass.GetRequest{
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
		req.SetRoundTripData(birdName)
		return req
	}
	return []*amass.GetRequest{
		makeReq(whatBirdOverviewPage),
		makeReq(whatBirdIdentificationPage),
		makeReq(whatBirdBehaviorPage),
	}
}

func reconstructWhatBirdsResponsesKeyedByLatinName(responses []*amass.GetResponse) map[string]*whatBirdResponse {
	m := make(map[string]map[string]*amass.GetResponse)
	m[whatBirdOverviewPage] = make(map[string]*amass.GetResponse)
	m[whatBirdIdentificationPage] = make(map[string]*amass.GetResponse)
	m[whatBirdBehaviorPage] = make(map[string]*amass.GetResponse)
	latinNames := make(map[string]bool)
	for _, response := range responses {
		if response.Site != whatBirdSite {
			continue
		}
		page := ""
		if strings.HasSuffix(response.RequestKey, whatBirdIdentificationPage) {
			page = whatBirdIdentificationPage
		} else if strings.HasSuffix(response.RequestKey, whatBirdOverviewPage) {
			page = whatBirdOverviewPage
		} else if strings.HasSuffix(response.RequestKey, whatBirdBehaviorPage) {
			page = whatBirdBehaviorPage
		} else {
			panic(fmt.Errorf("Unrecongnized response request key %s for what bird.", response.RequestKey))
		}
		birdName := &bird.BirdName{}
		response.GetRoundTripData(birdName)
		latinName := birdName.Latin
		latinNames[latinName] = true
		m[page][latinName] = response
	}
	result := make(map[string]*whatBirdResponse)
	for latinName, _ := range latinNames {
		result[latinName] = &whatBirdResponse{
			Identification: *m[whatBirdIdentificationPage][latinName],
			Behavior:       *m[whatBirdBehaviorPage][latinName],
			Overview:       *m[whatBirdOverviewPage][latinName],
		}
	}
	return result
}

func whatBirdRequestForTesting(englishName string) *whatBirdResponse {
	latin := "not really a latin name"
	bn := bird.BirdName{English: englishName, Latin: latin}
	rs := createWhatBirdRequests(bn)
	if len(rs) != 3 {
		panic(fmt.Errorf("Expected 3 what bird request, was %d, for key %s.", len(rs), englishName))
	}
	resps, err := amass.AmasserForTests().GetAll(rs)
	if err != nil {
		panic(fmt.Errorf("GetAll request failed for %s: %v", englishName, err))
	}
	m := reconstructWhatBirdsResponsesKeyedByLatinName(resps)
	return m[latin]
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
