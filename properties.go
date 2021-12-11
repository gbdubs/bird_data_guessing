package bird_data_guessing

import (
	"fmt"
	"regexp"
	"strconv"
)

type propertySearchers struct {
	food       *searcher
	nestType   *searcher
	habitat    *searcher
	funFact    *searcher
	wingspan   *searcher
	clutchSize *searcher
	all        *searcher
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
	return s.CountMatches([]string{"flock", "large flocks", "social", "large colonies", "large colony", "gregarious", "nests per colony", "roosts", "communal roost", "hundreds.{0,10}of birds"})
}

func (s *searcher) IsPredatory() *Property {
	return s.CountMatches([]string{" hawk ", "vulture", " falcon ", " eagle ", " prey ", "predator", "carrion", "small animals", "have killed", "condor", "hunter", "hunting", "scavenge", "stalk", "skewer", "impale", "striking", "talon", "hunt", "prey", "hunts", " dive "})
}

func (s *searcher) Food() *Food {
	var f Food
	f.Worm = *s.CountMatches([]string{"invertibrates", "invertebrates", "gnat", "aphid", "fruit flies", "spider", "mosquitoes", "gleaning", "hawking", "insect", "worm", " ant ", " wasp ", "beetle", " bug "})
	f.Berry = *s.CountMatches([]string{"fruit", "berries", "berry", "cherry", "nectar", "flower", " ripe "})
	f.Wheat = *s.CountMatches([]string{"grain", "wheat", "seed", " plant ", "vegetation", "shrub", " bush ", " nut ", " nuts "})
	f.Fish = *s.CountMatches([]string{"fish", "frog", "spawn", "shrimp", "molluscs", "eels", "fish", "fish", "probing", "crustacean", "snail"})
	f.Rat = *s.CountMatches([]string{"mammal", "rodent", "lizard", "mice", "vole", "squirrel"})
	f.Nectar = *s.CountMatches([]string{"nectar", "flower", "pollen", "sugar", "nectar", "corollae"})
	return &f
}

func (s *searcher) Habitat() *Habitat {
	var h Habitat
	h.Forest = *s.CountMatches([]string{"tree.?cover", "forest", "in forest", "in trees"})
	h.Grass = *s.CountMatches([]string{"grassland", "grass", "in grass", "prairie", "meadow"})
	h.Water = *s.CountMatches([]string{"in water", "near water", "marsh", "water.?bird", "water.?fowl", "bog"})
	return &h
}

func (r *searcher) NestType() *NestType {
	var n NestType
	n.Ground = *r.CountMatches([]string{"ground nest", "ground-dwelling", "scrape", "base of a", "shrub", "sagebrush", "grass", " lek ", " leks ", "shrub cover", "nesting cover", "ground"})
	n.Cup = *r.CountMatches([]string{"cup nest", "bowl", "deep bowl"})
	n.Slot = *r.CountMatches([]string{"cavity", "cavities", "tree-nesting", "boxes", "box", "hollow tree", "cave", "nest hole", "nesting hole", "dead tree", "cavity nest"})
	n.Platform = *r.CountMatches([]string{"platform", "build.{1,10}nest", "platform nest", "sticks", "large platform", "stick nest"})
	return &n
}

func (r *searcher) FunFact(englishName string) *Property {
	return r.ExtractMatch(
		"(\\.|\\n)\\s*(((all|some|most|the|a|an) )?"+englishName+"[a-zA-Z0-9\\-,'\"’‘”“ ]{40,150}\\.)", 2)
}

func (r *searcher) Wingspan() *Property {
	e := ".{0,50}"
	avgFemale := fmt.Sprintf("(female%saverage%swing.?span|wing.?span%sfemale%saverage|average%sfemale%swing.?span)", e, e, e, e, e, e)
	avg := fmt.Sprintf("(average%swing.?span|wing.?span%saverage)", e, e)
	args := []string{avgFemale, avg, "(wing.?span)", "(wing)", "(length|long|measure)"}
	units := []string{"cm|centimeter|centimetre", "millimeter|millimetre|mm", "meter|m", "inches|in"}
	for _, order := range []bool{true, false} {
		for _, arg := range args {
			for i, unit := range units {
				p := "[^.\\d](\\d+)(\\.\\d+)? ?(" + unit + ")"
				m := make(map[string]int)
				if order {
					m[fmt.Sprintf("%s%s%s", arg, e, p)] = 2
				} else {
					m[fmt.Sprintf("%s%s%s", p, e, arg)] = 1
				}
				a := r.ExtractAnyMatch(m)
				if a.StringValue != "" {
					switch i {
					case 0:
						a.IntValue = atoiOrFail(a.StringValue)
						break
					case 1:
						a.IntValue = atoiOrFail(a.StringValue) / 10
						break
					case 2:
						a.IntValue = atoiOrFail(a.StringValue) * 100
						break
					case 3:
						a.IntValue = int(float32(atoiOrFail(a.StringValue)) * 2.54)
					}
					if a.IntValue >= 5 /* Bee Hummingbird */ && a.IntValue <= 375 /* Wandering albatross */ {
						return a
					}
				}
			}
		}
	}
	return &Property{}
}

func (r *searcher) ClutchSize() *Property {
	p := r.getClutchMatch(true)
	if p.StringValue == "" {
		p = r.getClutchMatch(false)
	}
	if p.StringValue == "" {
		return p
	}
	s := p.StringValue
	s = caseInsReplace(s, "one", "1")
	s = caseInsReplace(s, "two", "2")
	s = caseInsReplace(s, "three", "3")
	s = caseInsReplace(s, "four", "4")
	s = caseInsReplace(s, "five", "5")
	s = caseInsReplace(s, "six", "6")
	s = caseInsReplace(s, "seven", "7")
	s = caseInsReplace(s, "eight", "8")
	s = caseInsReplace(s, "nine", "9")
	twoParts := caseInsensitiveRegex("(\\d+) ?(to|-|–) ?(\\d+)").FindStringSubmatch(s)
	if twoParts == nil {
		p.IntValue = atoiOrFail(s)
	} else {
		low := atoiOrFail(twoParts[1])
		high := atoiOrFail(twoParts[3])
		p.IntValue = (high + low) / 2
	}
	return p
}

func caseInsReplace(input string, lookFor string, replaceWith string) string {
	return regexp.MustCompile("(?i)"+lookFor).ReplaceAllString(input, replaceWith)
}

func (wr *searcher) getClutchMatch(matchRange bool) *Property {
	nr := "(one|two|three|four|five|six|seven|eight|nine|\\d+)(\\.\\d+)?"
	rr := fmt.Sprintf("(%s to %s|%s ?- ?%s|%s ?– ?%s)", nr, nr, nr, nr, nr, nr)
	np := "[^.\\d]*"
	r := nr
	if matchRange {
		r = rr
	}
	m := make(map[string]int)
	m[fmt.Sprintf("la(id|y)%s%s%segg", np, r, np)] = 2
	m[fmt.Sprintf("egg%sla(id|y)%s%s", np, np, r)] = 2
	m[fmt.Sprintf("clutch size%s%s", np, r)] = 1
	m[fmt.Sprintf("%s%seggs per year", r, np)] = 1
	m[fmt.Sprintf("clutch%s%s", np, r)] = 1
	m[fmt.Sprintf("%s%segg%sla(y|id)", r, np, np)] = 1
	m[fmt.Sprintf("la(id|y)%segg%s%s", np, np, r)] = 2
	return wr.ExtractAnyMatch(m)
}

func atoiOrFail(s string) int {
	i, e := strconv.Atoi(s)
	if e != nil {
		panic(e)
	}
	return i
}
