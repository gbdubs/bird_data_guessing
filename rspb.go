package bird_data_guessing

import (
	"fmt"
	"strings"

	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
)

type rspbResponse struct {
	Response amass.GetResponse
}

const (
	rspbSite                  = "rspb"
	maxRSPBConcurrentRequests = 2
)

func createRSPBRequests(birdName BirdName) []*amass.GetRequest {
	if isMissing(rspbSite, birdName) {
		return []*amass.GetRequest{}
	}
	nameParam := strings.ReplaceAll(strings.ReplaceAll(birdName.EnglishName, " ", "-"), "'", "")
	req := &amass.GetRequest{
		Site:                      rspbSite,
		RequestKey:                nameParam,
		URL:                       "https://rspb.org.uk/birds-and-wildlife/wildlife-guides/bird-a-z/" + nameParam,
		SiteMaxConcurrentRequests: maxRSPBConcurrentRequests,
		Attribution: attributions.Attribution{
			Author:              "The Royal Society for the Protection of Birds (RSPB)",
			AuthorUrl:           "https://rspb.org.uk",
			LicenseUrl:          "https://www.rspb.org.uk/help/terms--conditions",
			ScrapingMethodology: "github.com/gbdubs/bird_data_guessing/rspb",
		},
	}
	req.SetRoundTripData(birdName)
	return []*amass.GetRequest{req}
}

func reconstructRSPBResponsesKeyedByLatinName(responses []*amass.GetResponse) map[string]*rspbResponse {
	result := make(map[string]*rspbResponse)
	for _, response := range responses {
		if response.Site != rspbSite {
			continue
		}
		birdName := &BirdName{}
		response.GetRoundTripData(birdName)
		if isRSPBResponseMissing(response) {
			recordMissing(rspbSite, *birdName)
			continue
		}
		result[birdName.LatinName] = &rspbResponse{
			Response: *response,
		}
	}
	return result
}

func rspbRequestForTesting(englishName string) *rspbResponse {
	latin := "any old string"
	bn := BirdName{LatinName: latin, EnglishName: englishName}
	rs := createRSPBRequests(bn)
	if len(rs) != 1 {
		panic(fmt.Errorf("Expected 1 rspb request, was %d, for key %s.", len(rs), englishName))
	}
	resp, err := rs[0].Get()
	if err != nil {
		panic(fmt.Errorf("Get request failed for %s: %v", englishName, err))
	}
	m := reconstructRSPBResponsesKeyedByLatinName([]*amass.GetResponse{resp})
	return m[latin]
}

func isRSPBResponseMissing(r *amass.GetResponse) bool {
	return strings.Contains(r.AsDocument().Find("title").Text(), "Page Not Found")
}

func (r *rspbResponse) propertySearchers() *propertySearchers {
	page := r.Response.AsDocument()

	wingspanText := page.Find(".species-measurements-population__measurements").First().Text()
	eatText := page.Find(".key-information__section:contains('What they eat')").First().Text()
	habitatText := page.Find(".filter-block__tags-block:contains('Natural habitats')").First().Text()

	return &propertySearchers{
		wingspan: attributedSearch(&r.Response.Attribution, wingspanText),
		food:     attributedSearch(&r.Response.Attribution, eatText),
		habitat:  attributedSearch(&r.Response.Attribution, habitatText),
		// ClutchSize, Egg Color, Fun Fact, Nest Type, Predator and Flocking are omitted,
		// they aren't consistently populated
	}
}
