package bird_data_guessing

import (
	"fmt"
	"regexp"
	"strconv"
)

func (r *wikipediaResponse) IsFlocking() *Property {
	return r.CountMatches([]string{"flock", "large flocks", "social", "migratory"})
}

func (r *wikipediaResponse) IsPredatory() *Property {
	return r.CountMatches([]string{"prey", "predator", "carrion", "small animals", "have killed", "striking", "talon", "hunt", "prey"})
}

func (r *wikipediaResponse) Food() *Food {
	var f Food
	f.Worm = *r.CountMatches([]string{"invertibrates", "invertebrates", "gleaning", "hawking", "insect", "worm", "ant", "wasp"})
	f.Berry = *r.CountMatches([]string{"fruit", "berries", "berry", "cherry"})
	f.Wheat = *r.CountMatches([]string{"grain", "wheat", "seed", "plant"})
	f.Fish = *r.CountMatches([]string{"fish", "frog", "spawn", "shrimp", "molluscs", "eels", "fish", "fish", "snail"})
	f.Rat = *r.CountMatches([]string{"mammal", "rodent", "lizard", "mice", "vole", "squirrel"})
	f.Nectar = *r.CountMatches([]string{"nectar", "flower", "pollen"})
	return &f
}

func (r *wikipediaResponse) Habitat() *Habitat {
	var h Habitat
	h.Forest = *r.CountMatches([]string{"tree.?cover", "forest"})
	h.Grass = *r.CountMatches([]string{"grassland", "grass", "prairie", "meadow"})
	h.Water = *r.CountMatches([]string{"marsh", "water.?bird", "water.?fowl", "bog"})
	return &h
}

func (r *wikipediaResponse) NestType() *NestType {
	var n NestType
	n.Ground = *r.CountMatches([]string{"ground nest", "ground-dwelling", "ground"})
	n.Cup = *r.CountMatches([]string{"cup nest", "bowl"})
	n.Slot = *r.CountMatches([]string{"cavity", "cavities", "tree-nesting", "boxes", "box"})
	n.Platform = *r.CountMatches([]string{"platform", "build.{1,10}nest"})
	return &n
}

func (r *wikipediaResponse) FunFact(englishName string) *Property {
	return r.ExtractMatch(
		"(\\.|\\n)\\s*(((all|some|most|the|a|an) )?"+englishName+"[^\\.]{20,150}\\.)", 2)
}

func (r *wikipediaResponse) Wingspan() *Property {
	numberRegex := "\\D(\\d+)(\\.\\d+)?"
	inchesRegex := numberRegex + " ?in\\.?"
	cmRegex := numberRegex + " ?(cm|centimeter|centimetre)"
	np := "[^.\\d]*"
	avg := fmt.Sprintf("(female%saverage%swing.?span|wing.?span%sfemale%saverage|average%sfemale%swing.?span)", np, np, np, np, np, np)
	ws := "wing.?span"
	m := make(map[string]int)
	m[fmt.Sprintf("%s%s%s", avg, np, cmRegex)] = 2
	m[fmt.Sprintf("%s%s%s%s%s", avg, np, inchesRegex, np, cmRegex)] = 4
	p := r.ExtractAnyMatch(m)
	if p.StringValue == "" {
		m = make(map[string]int)
		m[fmt.Sprintf("%s%s%s", ws, np, cmRegex)] = 1
		m[fmt.Sprintf("%s%s%s%s%s", ws, np, inchesRegex, np, cmRegex)] = 4
		p = r.ExtractAnyMatch(m)
	}
	if p.StringValue == "" {
		m = make(map[string]int)
		m[fmt.Sprintf("%s.*%s", ws, cmRegex)] = 1
		p = r.ExtractAnyMatch(m)
	}
	if p.StringValue == "" {
		return p
	}
	p.IntValue = atoiOrFail(p.StringValue)
	return p
}

func (r *wikipediaResponse) ClutchSize() *Property {
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
	twoParts := caseInsensitiveRegex("(\\d+) ?(to|-) ?(\\d+)").FindStringSubmatch(s)
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

func (wr *wikipediaResponse) getClutchMatch(matchRange bool) *Property {
	nr := "(one|two|three|four|five|six|seven|eight|nine|\\d+)(.\\d+)?"
	rr := fmt.Sprintf("(%s to %s|%s ?- ?%s)", nr, nr, nr, nr)
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
