package bird_data_guessing

import (
	"strings"

	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
)

type audubonResponse struct {
	Response amass.GetResponse
}

const (
	audubonSite                  = "audubon"
	maxAudubonConcurrentRequests = 2
)

func createAudubonRequest(englishName string) *amass.GetRequest {
	nameParam := strings.ReplaceAll(strings.ReplaceAll(englishName, " ", "-"), "'", "")
	return &amass.GetRequest{
		Site:                      audubonSite,
		RequestKey:                nameParam,
		URL:                       "https://audubon.org/field-guide/bird/" + nameParam,
		SiteMaxConcurrentRequests: maxAudubonConcurrentRequests,
		Attribution: attributions.Attribution{
			Author:              "National Audubon Society, Inc.",
			AuthorUrl:           "https://audubon.org",
			License:             "All rights reserved",
			LicenseUrl:          "https://www.audubon.org/terms-use",
			ScrapingMethodology: "github.com/gbdubs/bird_data_guessing/audubon",
		},
	}
}

func reconstructAudubonResponses(responses []*amass.GetResponse) []*audubonResponse {
	result := make([]*audubonResponse, 0)
	for _, response := range responses {
		if response.Site != allAboutBirdsSite {
			continue
		}
		result = append(result, &audubonResponse{
			Response: *response,
		})
	}
	return result
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
