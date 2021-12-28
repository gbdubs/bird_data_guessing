package bird_data_guessing

import (
	"fmt"
	"math"

	"github.com/davecgh/go-spew/spew"
	"github.com/gbdubs/attributions"
	"github.com/gbdubs/inference"
)

func mergeWingspan(inputs []*inference.Float64Range) (*inference.Float64Range, bool) {
	if len(inputs) == 0 {
		return &inference.Float64Range{}, false
	}
	minMin := 10000.0
	maxMin := -1.0
	minMax := 10000.0
	maxMax := -1.0
	minTotal := 0.0
	maxTotal := 0.0
	for _, input := range inputs {
		if minMin > input.Min {
			minMin = input.Min
		}
		if maxMin < input.Min {
			maxMin = input.Min
		}
		if minMax > input.Max {
			minMax = input.Max
		}
		if maxMax < input.Max {
			maxMax = input.Max
		}
		minTotal += input.Min
		maxTotal += input.Max
	}
	highConfidence := true
	if maxMin*.8 > minMax*1.2 || minMin < 0 || maxMax > 380 /* Wandering Albatros 375cm */ {
		fmt.Printf("\n\n\nLOW CONFIDENCE WINGSPAN [%f, %f] to [%f, %f]\n", minMin, maxMin, minMax, maxMax)
		for _, input := range inputs {
			fmt.Printf("%f to %f source %#v \n", input.Min, input.Max, input.Source.Origin.Attribution.OriginUrl)
		}
		fmt.Printf("\n\n")
		highConfidence = false
	}
	min := minTotal / float64(len(inputs))
	max := maxTotal / float64(len(inputs))
	result := &inference.Float64Range{
		Min: min,
		Max: max,
		Source: inference.CombineSources(
			"averaged ranges",
			fmt.Sprintf("%f - %f", min, max),
			inference.AsSourceables(inputs)...,
		),
	}
	return result, highConfidence
}

func mergeClutchSize(inputs []*inference.IntRange) (*inference.IntRange, bool) {
	if len(inputs) == 0 {
		return &inference.IntRange{}, false
	}
	minMin := math.MaxInt
	maxMin := math.MinInt
	minMax := math.MaxInt
	maxMax := math.MinInt
	minTotal := 0
	minCount := 0
	maxTotal := 0
	maxCount := 0
	for _, input := range inputs {
		if input.Min == 0 && input.Max == 0 {
			continue
		}
		// Exclude open ranges + single points from minimum calculations.
		if input.Min != -1 && input.Min != input.Max {
			if minMin > input.Min {
				minMin = input.Min
			}
			if maxMin < input.Min {
				maxMin = input.Min
			}
			minTotal += input.Min
			minCount++
		}
		if minMax > input.Max {
			minMax = input.Max
		}
		if maxMax < input.Max {
			maxMax = input.Max
		}
		maxTotal += input.Max
		maxCount++
	}
	highConfidence := true
	if (minCount > 0 && (minMin < 0 || maxMin > minMax+1)) || minMin < 0 || maxMax > 33 /* Maximum avian clutch size is 33 */ {
		fmt.Printf("\n\n\nLOW CONFIDENCE CLUTCH SIZE[%d, %d] to [%d, %d]\n", minMin, maxMin, minMax, maxMax)
		for _, input := range inputs {
			fmt.Printf("%d to %d source %s \n", input.Min, input.Max, spew.Sdump(input))
		}
		fmt.Printf("\n\n")
		highConfidence = false
	}
	avgMin := 1
	if minCount > 0 {
		avgMin = int(math.Floor(float64(minTotal) / float64(minCount)))
	}
	avgMax := int(math.Ceil(float64(maxTotal) / float64(maxCount)))
	result := &inference.IntRange{
		Min: avgMin,
		Max: avgMax,
		Source: inference.CombineSources(
			"took mean'ed min/max of input ranges, with tollerance for open ranges + point estimates",
			fmt.Sprintf("%d to %d", minMin, maxMax),
			inference.AsSourceables(inputs)...,
		),
	}
	return result, highConfidence
}

func mergeEggColor(inputs []*inference.String) (*inference.String, bool) {
	if len(inputs) == 0 {
		return inference.NewString("", "No egg color found", &attributions.Attribution{}), false
	}
	return inference.RandomChoiceString(inputs), true
}

func mergeFunFact(inputs []*inference.String) (*inference.String, bool) {
	if len(inputs) == 0 {
		return inference.NewString("", "No fun fact found", &attributions.Attribution{}), false
	}
	return inference.RandomChoiceString(inputs), true
}

func mergeScores(inputs []*inference.Int) *inference.Float64 {
	return inference.ZeroTolerantGeomMeanInt(inputs...)
}

func (input *Input) mergeSources(inputs []*singleSourceData) (data *BirdData, highConfidence bool) {
	if len(inputs) == 0 {
		panic(fmt.Errorf("expected one or more inputs to merge, but had none."))
	}
	wingspans := make([]*inference.Float64Range, 0)
	clutchSizes := make([]*inference.IntRange, 0)
	eggColors := make([]*inference.String, 0)
	funFacts := make([]*inference.String, 0)

	wheatScores := make([]*inference.Int, 0)
	wormScores := make([]*inference.Int, 0)
	berryScores := make([]*inference.Int, 0)
	mouseScores := make([]*inference.Int, 0)
	fishScores := make([]*inference.Int, 0)
	nectarScores := make([]*inference.Int, 0)

	cavityScores := make([]*inference.Int, 0)
	cupScores := make([]*inference.Int, 0)
	platformScores := make([]*inference.Int, 0)
	groundScores := make([]*inference.Int, 0)

	forestScores := make([]*inference.Int, 0)
	grassScores := make([]*inference.Int, 0)
	waterScores := make([]*inference.Int, 0)

	predatorScores := make([]*inference.Int, 0)
	flockingScores := make([]*inference.Int, 0)

	for _, input := range inputs {
		wingspans = append(wingspans, input.Wingspan...)
		clutchSizes = append(clutchSizes, input.ClutchSize)
		eggColors = append(eggColors, input.EggColor...)
		funFacts = append(funFacts, input.FunFact...)

		wheatScores = append(wheatScores, input.WheatScore)
		wormScores = append(wormScores, input.WormScore)
		berryScores = append(berryScores, input.BerryScore)
		mouseScores = append(mouseScores, input.MouseScore)
		fishScores = append(fishScores, input.FishScore)
		nectarScores = append(nectarScores, input.NectarScore)

		cavityScores = append(cavityScores, input.CavityScore)
		cupScores = append(cupScores, input.CupScore)
		platformScores = append(platformScores, input.PlatformScore)
		groundScores = append(groundScores, input.GroundScore)

		forestScores = append(forestScores, input.ForestScore)
		grassScores = append(grassScores, input.GrassScore)
		waterScores = append(waterScores, input.WaterScore)

		predatorScores = append(predatorScores, input.PredatorScore)
		flockingScores = append(flockingScores, input.FlockingScore)
	}

	data = &BirdData{}
	highConfidence = true

	if ws, hc := mergeWingspan(wingspans); hc {
		data.Wingspan = ws
	} else {
		input.VLog(" - LC - Wingspan - ")
		highConfidence = false
	}
	if cs, hc := mergeClutchSize(clutchSizes); hc {
		data.ClutchSize = cs
	} else {
		input.VLog(" - LC - ClutchSize -")
		highConfidence = false
	}
	if ec, hc := mergeEggColor(eggColors); hc {
		data.EggColor = ec
	} else {
		// Note: we're OK with egg color being missing, unlike the other 3 high-specificity properties, because we have a reasonable fallback.
	}
	if ff, hc := mergeFunFact(funFacts); hc {
		data.FunFact = ff
	} else {
		input.VLog(" - LC - Fun Fact -")
		highConfidence = false
	}
	/*
		data.WheatScore = mergeScores(wheatScores)
		data.WormScore = mergeScores(wormScores)
		data.BerryScore = mergeScores(berryScores)
		data.MouseScore = mergeScores(mouseScores)
		data.FishScore = mergeScores(fishScores)
		data.NectarScore = mergeScores(nectarScores)

		data.CavityScore = mergeScores(cavityScores)
		data.CupScore = mergeScores(cupScores)
		data.GroundScore = mergeScores(groundScores)
		data.PlatformScore = mergeScores(platformScores)

		data.ForestScore = mergeScores(forestScores)
		data.GrassScore = mergeScores(grassScores)
		data.WaterScore = mergeScores(waterScores)

		data.PredatorScore = mergeScores(predatorScores)
		data.FlockingScore = mergeScores(flockingScores)
	*/
	return
}
