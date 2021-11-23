package bird_data_guessing

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
