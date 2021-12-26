package bird_data_guessing

import (
	"github.com/gbdubs/bird"
	"github.com/gbdubs/inference"
)

type Input struct {
	Names []bird.BirdName
}

type Output struct {
	BirdData []BirdData
}

type BirdData struct {
	Name bird.BirdName
	// Precice Properties
	Wingspan   *inference.Float64Range
	ClutchSize *inference.IntRange
	EggColor   *inference.String
	FunFact    *inference.String
	// Scored Properties (affinity measure)
	// Food
	WheatScore  *inference.Float64
	WormScore   *inference.Float64
	BerryScore  *inference.Float64
	MouseScore  *inference.Float64
	FishScore   *inference.Float64
	NectarScore *inference.Float64
	// Habitat
	ForestScore *inference.Float64
	GrassScore  *inference.Float64
	WaterScore  *inference.Float64
	// Nest Type
	CupScore      *inference.Float64
	GroundScore   *inference.Float64
	PlatformScore *inference.Float64
	CavityScore   *inference.Float64
	// Behavior
	FlockingScore *inference.Float64
	PredatorScore *inference.Float64
}

type singleSourceData struct {
	Name bird.BirdName
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
