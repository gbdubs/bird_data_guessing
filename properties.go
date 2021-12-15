package bird_data_guessing

// Tells the searching algorithms where to look for different pieces of data.
// This CAN be a single searcher across all categories, but for most sites
// we can do some basic HTML parsing to decide which section is most likely
// to have the data that we want in it.
type propertySearchers struct {
	wingspan   searcher
	clutchSize searcher
	eggColor   searcher
	funFact    searcher
	food       searcher
	nestType   searcher
	habitat    searcher
	predator   searcher
	flocking   searcher
}

func (s *propertySearchers) getData(englishName string) *zZZSingleSourceData {
	d := zZZSingleSourceData{}

	d.EnglishName = englishName

	d.Wingspan = s.wingspan.Wingspan()
	d.ClutchSize = s.clutchSize.ClutchSize()
	d.EggColor = s.eggColor.EggColor()
	d.FunFact = s.funFact.FunFact(englishName)

	d.WormScore = s.food.ZZZCountMatches(wormKeywords...)
	d.WheatScore = s.food.ZZZCountMatches(wheatKeywords...)
	d.BerryScore = s.food.ZZZCountMatches(berryKeywords...)
	d.MouseScore = s.food.ZZZCountMatches(mouseKeywords...)
	d.FishScore = s.food.ZZZCountMatches(fishKeywords...)
	d.NectarScore = s.food.ZZZCountMatches(nectarKeywords...)

	d.CavityScore = s.nestType.ZZZCountMatches(cavityKeywords...)
	d.CupScore = s.nestType.ZZZCountMatches(cupKeywords...)
	d.GroundScore = s.nestType.ZZZCountMatches(groundKeywords...)
	d.PlatformScore = s.nestType.ZZZCountMatches(platformKeywords...)

	d.ForestScore = s.habitat.ZZZCountMatches(forestKeywords...)
	d.GrassScore = s.habitat.ZZZCountMatches(grassKeywords...)
	d.WaterScore = s.habitat.ZZZCountMatches(waterKeywords...)

	d.PredatorScore = s.predator.ZZZCountMatches(predatorKeywords...)
	d.FlockingScore = s.flocking.ZZZCountMatches(flockingKeywords...)

	return &d
}

// Food Keywords
var wormKeywords = []string{"invertibrates", "invertebrates", "gnat", "aphid", "fruit flies", "spider", "mosquitoes", "gleaning", "hawking", "insect", "worm", " ant ", " wasp ", "beetle", " bug "}
var berryKeywords = []string{"fruit", "berries", "berry", "cherry", "nectar", "flower", " ripe "}
var wheatKeywords = []string{"grain", "wheat", "seed", " plant ", "vegetation", "shrub", " bush ", " nut ", " nuts "}
var fishKeywords = []string{"fish", "frog", "spawn", "shrimp", "molluscs", "eels", "fish", "fish", "probing", "crustacean", "snail"}
var mouseKeywords = []string{"mammal", "mice", "rodent", "lizard", "mice", "vole", "squirrel"}
var nectarKeywords = []string{"nectar", "flower", "pollen", "sugar", "nectar", "corollae"}

// Nest Type Keywords
var groundKeywords = []string{"ground nest", "ground-dwelling", "on ground", "on ground", "scrape", "base of a", "shrub", "sagebrush", "grass", " lek ", " leks ", "shrub cover", "nesting cover"}
var cupKeywords = []string{"cup nest", "bowl", "above ground", "feet above ground", "deep bowl"}
var cavityKeywords = []string{"cavity", "cavities", "tree-nesting", "tree cavity", "woodpecker cavit", "boxes", "box", "hollow tree", "cave", "nest hole", "nesting hole", "dead tree", "cavity nest"}
var platformKeywords = []string{"platform", "build.{1,10}nest", "platform nest", "sticks", "large platform", "stick nest"}

// Habitat Keywords
var forestKeywords = []string{"tree.?cover", "forest", "in forest", "in trees", "woodland", "understory", "canopy", "conifer", "evergreen", "groves"}
var grassKeywords = []string{"grassland", "grass", "in grass", "prairie", "meadow", "scrub", "arid", "farmland"}
var waterKeywords = []string{"in water", "near water", "marsh", "water.?bird", "water.?fowl", "bog", "lake", "floodplain", "riparian", "brackish"}

// Behavior Keywords
var flockingKeywords = []string{"flock", "large flocks", "social", "large colonies", "large colony", "gregarious", "nests per colony", "roosts", "communal roost", "hundreds.{0,10}of birds"}
var predatorKeywords = []string{"raptor", "birds? of prey", " hawk ", "vulture", " falcon ", " eagle ", " prey ", "predator", "carrion", "small animals", "have killed", "condor", "hunter", "hunting", "scavenge", "stalk", "skewer", "impale", "striking", "talon", "hunt", "hunts", " dive "}
