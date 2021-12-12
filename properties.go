package bird_data_guessing

import (
	"regexp"
	"strconv"
)

type propertySearchers struct {
	eggColor   searcher
	food       searcher
	nestType   searcher
	habitat    searcher
	funFact    searcher
	wingspan   searcher
	clutchSize searcher
	all        searcher
}

func (s *propertySearchers) getData(englishName string) (*Data, *DebugData) {
	d := Data{}
	dd := DebugData{}

	d.EnglishName = englishName

	f := s.food.Food()
	d.WormScore = f.Worm.Strength
	d.WheatScore = f.Wheat.Strength
	d.BerryScore = f.Berry.Strength
	d.FishScore = f.Fish.Strength
	d.RatScore = f.Rat.Strength
	d.NectarScore = f.Nectar.Strength
	dd.Food = *f

	n := s.nestType.NestType()
	d.CupScore = n.Cup.Strength
	d.GroundScore = n.Ground.Strength
	d.PlatformScore = n.Platform.Strength
	d.SlotScore = n.Slot.Strength
	dd.NestType = *n

	h := s.habitat.Habitat()
	d.ForestScore = h.Forest.Strength
	d.GrassScore = h.Grass.Strength
	d.WaterScore = h.Water.Strength
	dd.Habitat = *h

	ff := s.funFact.FunFact(englishName)
	d.FunFact = ff.StringValue
	dd.FunFact = *ff

	ec := s.eggColor.EggColor()
	d.EggColor = ec.StringValue
	dd.EggColor = *ec

	w := s.wingspan.Wingspan()
	d.Wingspan = w.IntValue
	dd.Wingspan = *w

	cs := s.clutchSize.ClutchSize()
	d.ClutchSize = cs.IntValue
	dd.ClutchSize = *cs

	pr := s.all.IsPredatory()
	d.PredatorScore = pr.Strength
	dd.IsPredatory = *pr

	fl := s.all.IsFlocking()
	d.FlockingScore = fl.Strength
	dd.IsFlocking = *fl

	dd.Data = d

	return &d, &dd
}

func (s *searcher) IsFlocking() *Property {
	return s.CountMatches("flock", "large flocks", "social", "large colonies", "large colony", "gregarious", "nests per colony", "roosts", "communal roost", "hundreds.{0,10}of birds")
}

func (s *searcher) IsPredatory() *Property {
	return s.CountMatches(" hawk ", "vulture", " falcon ", " eagle ", " prey ", "predator", "carrion", "small animals", "have killed", "condor", "hunter", "hunting", "scavenge", "stalk", "skewer", "impale", "striking", "talon", "hunt", "prey", "hunts", " dive ")
}

func (s *searcher) Food() *Food {
	var f Food
	f.Worm = *s.CountMatches("invertibrates", "invertebrates", "gnat", "aphid", "fruit flies", "spider", "mosquitoes", "gleaning", "hawking", "insect", "worm", " ant ", " wasp ", "beetle", " bug ")
	f.Berry = *s.CountMatches("fruit", "berries", "berry", "cherry", "nectar", "flower", " ripe ")
	f.Wheat = *s.CountMatches("grain", "wheat", "seed", " plant ", "vegetation", "shrub", " bush ", " nut ", " nuts ")
	f.Fish = *s.CountMatches("fish", "frog", "spawn", "shrimp", "molluscs", "eels", "fish", "fish", "probing", "crustacean", "snail")
	f.Rat = *s.CountMatches("mammal", "rodent", "lizard", "mice", "vole", "squirrel")
	f.Nectar = *s.CountMatches("nectar", "flower", "pollen", "sugar", "nectar", "corollae")
	return &f
}

func (s *searcher) Habitat() *Habitat {
	var h Habitat
	h.Forest = *s.CountMatches("tree.?cover", "forest", "in forest", "in trees", "woodland", "understory", "canopy", "conifer", "evergreen", "groves")
	h.Grass = *s.CountMatches("grassland", "grass", "in grass", "prairie", "meadow", "scrub", "arid", "farmland")
	h.Water = *s.CountMatches("in water", "near water", "marsh", "water.?bird", "water.?fowl", "bog", "lake", "floodplain", "riparian", "brackish")
	return &h
}

func (r *searcher) NestType() *NestType {
	var n NestType
	n.Ground = *r.CountMatches("ground nest", "ground-dwelling", "on ground", "on ground", "scrape", "base of a", "shrub", "sagebrush", "grass", " lek ", " leks ", "shrub cover", "nesting cover")
	n.Cup = *r.CountMatches("cup nest", "bowl", "above ground", "feet above ground", "deep bowl")
	n.Slot = *r.CountMatches("cavity", "cavities", "tree-nesting", "tree cavity", "woodpecker cavit", "boxes", "box", "hollow tree", "cave", "nest hole", "nesting hole", "dead tree", "cavity nest")
	n.Platform = *r.CountMatches("platform", "build.{1,10}nest", "platform nest", "sticks", "large platform", "stick nest")
	return &n
}

func (r *searcher) FunFact(englishName string) *Property {
	return r.ExtractMatch(
		"(\\.|\\n)\\s*(((all|some|most|the|a|an) )?"+englishName+"[a-zA-Z0-9\\-,'\"’‘”“ ]{40,150}\\.)", 2)
}

func caseInsReplace(input string, lookFor string, replaceWith string) string {
	return regexp.MustCompile("(?i)"+lookFor).ReplaceAllString(input, replaceWith)
}

func atoiOrFail(s string) int {
	i, e := strconv.Atoi(s)
	if e != nil {
		panic(e)
	}
	return i
}

func floatOrFail(s string) float64 {
	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		panic(e)
	}
	return f
}
