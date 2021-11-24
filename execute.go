package bird_data_guessing

import (
	"fmt"
	"time"
)

func (i *Input) Execute() (*Output, error) {
	var oo Output
	o := &oo
	r, err := getWikipediaPage(i.LatinName)
	if err != nil {
		return o, err
	}
	o.Data.EnglishName = i.EnglishName
	o.Data.LatinName = i.LatinName

	f := r.Food()
	o.Data.WormScore = f.Worm.Strength
	o.Data.WheatScore = f.Wheat.Strength
	o.Data.BerryScore = f.Berry.Strength
	o.Data.FishScore = f.Fish.Strength
	o.Data.RatScore = f.Rat.Strength
	o.Data.NectarScore = f.Nectar.Strength

	n := r.NestType()
	o.Data.CupScore = n.Cup.Strength
	o.Data.GroundScore = n.Ground.Strength
	o.Data.PlatformScore = n.Platform.Strength
	o.Data.SlotScore = n.Slot.Strength

	h := r.Habitat()
	o.Data.ForestScore = h.Forest.Strength
	o.Data.GrassScore = h.Grass.Strength
	o.Data.WaterScore = h.Water.Strength

	o.Data.FunFact = r.FunFact(i.EnglishName).StringValue
	o.Data.Wingspan = r.Wingspan().IntValue
	o.Data.ClutchSize = r.ClutchSize().IntValue
	o.Data.PredatorScore = r.IsPredatory().Strength
	o.Data.FlockingScore = r.IsFlocking().Strength

	o.Attribution.OriginUrl = r.Query.Pages.Page.Canonicalurl
	o.Attribution.CollectedAt = time.Now()
	o.Attribution.OriginalTitle = r.Query.Pages.Page.Title
	o.Attribution.Author = "Wikimedia Foundation"
	o.Attribution.AuthorUrl = "https://wikimedia.org"
	sourceTimestamp := r.Query.Pages.Page.Cirrusdoc.V.Source.Timestamp
	o.Attribution.CreatedAt, err = time.Parse(time.RFC3339, sourceTimestamp)
	if err != nil {
		fmt.Printf("Error when looking at element: %s", i.LatinName)
		panic(err)
	}
	o.Attribution.License = "Creative Commons Attribution-ShareAlike 3.0 Unported License (CC BY-SA)"
	o.Attribution.LicenseUrl = "https://en.wikipedia.org/wiki/Wikipedia:Text_of_Creative_Commons_Attribution-ShareAlike_3.0_Unported_License"
	o.Attribution.ScrapingMethodology = "github.com/gbdubs/bird_data_guessing"
	o.Attribution.Context = []string{"Called Wikipedia's API with action=query, see api.go for details."}
	//	o.setDebuggingFields(i.EnglishName, r)
	return o, err
}

/* Debug-only method */
func (o *Output) setDebuggingFields(en string, r *wikipediaResponse) {
	o.DebugData.Habitat = *r.Habitat()
	o.DebugData.IsPredatory = *r.IsPredatory()
	o.DebugData.IsFlocking = *r.IsFlocking()
	o.DebugData.NestType = *r.NestType()
	o.DebugData.FunFact = *r.FunFact(en)
	o.DebugData.ClutchSize = *r.ClutchSize()
	o.DebugData.Wingspan = *r.Wingspan()
}
