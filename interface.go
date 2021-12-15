package bird_data_guessing

import (
	"github.com/gbdubs/attributions"
	"github.com/gbdubs/inference"
)

type Input struct {
	LatinName   string
	EnglishName string // Optional.
	Debug       bool
}

type Output struct {
	ZZZData      ZZZData
	Attributions []attributions.Attribution
}

type ZZZData struct {
	EnglishName string
	LatinName   string
	// Precice Properties
	Wingspan   *inference.Float64Range
	ClutchSize *inference.IntRange
	EggColor   *inference.String
	FunFact    *inference.String
	// Scored Properties (affinity measure)
	// Food
	WheatScore  *inference.Int
	WormScore   *inference.Int
	BerryScore  *inference.Int
	MouseScore  *inference.Int
	FishScore   *inference.Int
	NectarScore *inference.Int
	// Habitat
	ForestScore *inference.Int
	GrassScore  *inference.Int
	WaterScore  *inference.Int
	// Nest Type
	CupScore      *inference.Int
	GroundScore   *inference.Int
	PlatformScore *inference.Int
	CavityScore   *inference.Int
	// Behavior
	FlockingScore *inference.Int
	PredatorScore *inference.Int
}

type zZZSingleSourceData struct {
	EnglishName string
	LatinName   string
	// Precice Properties
	Wingspan   []*inference.Float64Range
	ClutchSize *inference.IntRange
	EggColor   []*inference.String
	FunFact    []*inference.String
	// Scored Properties (affinity measure)
	// Food
	WheatScore  *inference.Int
	WormScore   *inference.Int
	BerryScore  *inference.Int
	MouseScore  *inference.Int
	FishScore   *inference.Int
	NectarScore *inference.Int
	// Habitat
	ForestScore *inference.Int
	GrassScore  *inference.Int
	WaterScore  *inference.Int
	// Nest Type
	CupScore      *inference.Int
	GroundScore   *inference.Int
	PlatformScore *inference.Int
	CavityScore   *inference.Int
	// Behavior
	FlockingScore *inference.Int
	PredatorScore *inference.Int
}

// Debugging Types below this line

type Habitat struct {
	Forest Property
	Water  Property
	Grass  Property
}

type Food struct {
	Worm   Property
	Wheat  Property
	Berry  Property
	Fish   Property
	Rat    Property
	Nectar Property
}

type NestType struct {
	Ground   Property
	Cup      Property
	Slot     Property
	Platform Property
}

type Property struct {
	Strength    int
	Context     string
	StringValue string
	IntValue    int
}
