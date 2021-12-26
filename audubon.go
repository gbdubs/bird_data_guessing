package bird_data_guessing

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
	"github.com/gbdubs/bird"
	"github.com/gbdubs/sitemaps"
)

type audubonResponse struct {
	Response amass.GetResponse
}

const (
	audubonSite                  = "audubon"
	maxAudubonConcurrentRequests = 3
)

var audubonSiteMap *sitemaps.Sitemap = nil

func audubonSitemap() *sitemaps.Sitemap {
	if audubonSiteMap == nil {
		s, err := sitemaps.GetPagedSitemapFromURL("https://www.audubon.org/sitemap.xml")
		if err != nil {
			panic(err)
		}
		audubonSiteMap = s
	}
	return audubonSiteMap
}

func createAudubonRequests(birdName bird.BirdName) []*amass.GetRequest {
	if isMissing(audubonSite, birdName) {
		return []*amass.GetRequest{}
	}
	nameParam := strings.ReplaceAll(strings.ReplaceAll(birdName.English, " ", "-"), "'", "")
	// For whatever reason, audubon's sitemap references .ngo links, but these 400 when you
	// actually request them. However, the same URL pattern DOES work when you use the .org TLD
	ngoTargetURL := "https://audubon.ngo/field-guide/bird/" + nameParam
	ngoURL, levDist := audubonSitemap().BestFuzzyMatch(ngoTargetURL)
	if levDist > 3 {
		recordMissing(rspbSite, birdName)
		return []*amass.GetRequest{}
	}
	requestKey := regexp.MustCompile("audubon.ngo/field-guide/bird/([^/]+)").FindStringSubmatch(ngoURL)[1]
	actualURL := "https://audubon.org/field-guide/bird/" + requestKey
	req := &amass.GetRequest{
		Site:                      audubonSite,
		RequestKey:                requestKey,
		URL:                       actualURL,
		SiteMaxConcurrentRequests: maxAudubonConcurrentRequests,
		Attribution: attributions.Attribution{
			Author:              "National Audubon Society, Inc.",
			AuthorUrl:           "https://audubon.org",
			License:             "All rights reserved",
			LicenseUrl:          "https://www.audubon.org/terms-use",
			ScrapingMethodology: "github.com/gbdubs/bird_data_guessing/audubon",
			CreatedAt:           audubonSitemap().LastUpdated[ngoURL],
		},
	}
	req.SetRoundTripData(birdName)
	return []*amass.GetRequest{req}
}

func reconstructAudubonResponsesKeyedByLatinName(responses []*amass.GetResponse) map[string]*audubonResponse {
	result := make(map[string]*audubonResponse)
	for _, response := range responses {
		if response.Site != audubonSite {
			continue
		}
		birdName := &bird.BirdName{}
		response.GetRoundTripData(birdName)
		if isAudubonResponseMissing(response) {
			recordMissing(audubonSite, *birdName)
			continue
		}
		result[birdName.Latin] = &audubonResponse{
			Response: *response,
		}
	}
	return result
}

func audubonRequestForTesting(englishName string) *audubonResponse {
	latin := "any old string"
	bn := bird.BirdName{Latin: latin, English: englishName}
	rs := createAudubonRequests(bn)
	if len(rs) != 1 {
		panic(fmt.Errorf("Expected 1 audubon request, was %d, for key %s.", len(rs), englishName))
	}
	resp, err := rs[0].Get()
	if err != nil {
		panic(fmt.Errorf("Get request failed for %s: %v", englishName, err))
	}
	m := reconstructAudubonResponsesKeyedByLatinName([]*amass.GetResponse{resp})
	return m[latin]
}

func isAudubonResponseMissing(r *amass.GetResponse) bool {
	return strings.Contains(r.AsDocument().Find("title").Text(), "Sorry, We Couldn't Find That Page")
}

func (r *audubonResponse) propertySearchers() *propertySearchers {
	page := r.Response.AsDocument()

	dietText := page.Find("h2:contains('Diet')").First().Next().Text()
	feedingText := page.Find("h2:contains('Feeding')").First().Next().Text()
	eggsText := "EggsAudubonEggs " + page.Find("h2:contains('Eggs')").First().Next().Text()
	nestingText := page.Find("h2:contains('Nesting')").First().Next().Text()
	habitatText := page.Find("th:contains('Habitat')").First().Parent().Find("td").First().Text()
	allText := page.Find("body").Text()

	return &propertySearchers{
		// Wingspan is omitted, it isn't consistently helpful.
		clutchSize: attributedSearch(&r.Response.Attribution, eggsText+nestingText),
		eggColor:   attributedSearch(&r.Response.Attribution, eggsText),
		// Fun fact is omitted, it's not reliably fun.

		food:     attributedSearch(&r.Response.Attribution, dietText+feedingText),
		nestType: attributedSearch(&r.Response.Attribution, nestingText),
		habitat:  attributedSearch(&r.Response.Attribution, habitatText+nestingText),

		predator: attributedSearch(&r.Response.Attribution, allText),
		flocking: attributedSearch(&r.Response.Attribution, allText),
	}
}
