package bird_data_guessing

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/gbdubs/amass"
	"github.com/gbdubs/attributions"
	"github.com/gbdubs/bird"
	"github.com/gbdubs/sitemaps"
)

type rspbResponse struct {
	Response amass.GetResponse
}

const (
	rspbSite                  = "rspb"
	maxRSPBConcurrentRequests = 2
)

var rspbSiteMap *sitemaps.Sitemap = nil
var rspbSiteMapLock = sync.RWMutex{}

func rspbSitemap() *sitemaps.Sitemap {
	if rspbSiteMap == nil {
		rspbSiteMapLock.Lock()
		if rspbSiteMap == nil {
			s, err := sitemaps.GetSitemapFromURL("https://www.rspb.org.uk/birds-and-wildlife/sitemap.xml")
			if err != nil {
				panic(err)
			}
			rspbSiteMap = s
		}
		rspbSiteMapLock.Unlock()
	}
	return rspbSiteMap
}

func createRSPBRequests(birdName bird.BirdName) []*amass.GetRequest {
	sitemap := rspbSitemap()
	if isMissing(rspbSite, birdName) {
		return []*amass.GetRequest{}
	}
	nameParam := strings.ReplaceAll(strings.ReplaceAll(birdName.English, " ", "-"), "'", "")
	url := "https://rspb.org.uk/birds-and-wildlife/wildlife-guides/bird-a-z/" + nameParam
	actualUrl, levDist := sitemap.BestFuzzyMatch(url)
	if levDist > 2 {
		recordMissing(rspbSite, birdName)
		return []*amass.GetRequest{}
	}
	requestKey := regexp.MustCompile("bird-a-z/(.+)/$").FindStringSubmatch(actualUrl)[1]
	req := &amass.GetRequest{
		Site:                      rspbSite,
		RequestKey:                requestKey,
		URL:                       actualUrl,
		SiteMaxConcurrentRequests: maxRSPBConcurrentRequests,
		Attribution: attributions.Attribution{
			CreatedAt:           sitemap.LastUpdated[actualUrl],
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
		birdName := &bird.BirdName{}
		response.GetRoundTripData(birdName)
		if isRSPBResponseMissing(response) {
			recordMissing(rspbSite, *birdName)
			continue
		}
		result[birdName.Latin] = &rspbResponse{
			Response: *response,
		}
	}
	return result
}

func rspbRequestForTesting(englishName string) *rspbResponse {
	latin := "any old string"
	bn := bird.BirdName{Latin: latin, English: englishName}
	rs := createRSPBRequests(bn)
	if len(rs) != 1 {
		panic(fmt.Errorf("Expected 1 rspb request, was %d, for key %s.", len(rs), englishName))
	}
	resp, err := rs[0].Get()
	if err != nil {
		panic(fmt.Errorf("Get request failed for %s: %v", englishName, err))
	}
	m := reconstructRSPBResponsesKeyedByLatinName([]*amass.GetResponse{resp})
	result, ok := m[latin]
	if !ok {
		result = &rspbResponse{}
	}
	return result
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
