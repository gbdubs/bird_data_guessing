package bird_data_guessing

import (
	"github.com/gbdubs/inference"
)

type Input struct {
	Birds []BirdName
	Debug bool
}

type BirdName struct {
	EnglishName string
	LatinName   string
}

type Output struct {
	BirdData []BirdData
}

type BirdData struct {
	Name BirdName
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

type singleSourceData struct {
	Name BirdName
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
