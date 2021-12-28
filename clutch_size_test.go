package bird_data_guessing

import (
	"testing"

	"github.com/gbdubs/inference"
)

func TestClutchSize_Wikipedia(t *testing.T) {
	wikipediaClutchSize := func(latinName string) *inference.IntRange {
		return wikipediaRequestForTesting(latinName).propertySearchers().clutchSize.ClutchSize()
	}
	testRangeIntCases(t, wikipediaClutchSize,
		rangeIntCase("Aythya valisineria", 5, 11),
		rangeIntCase("Cardinalis sinuatus", 2, 4),
		rangeIntCase("Cistothorus palustris", 4, 6),
		rangeIntCase("Passer domesticus", 4, 5),
		rangeIntCase("Pooecetes gramineus", 3, 5),
		rangeIntCase("Loxia leucoptera", 3, 5),
		rangeIntCase("Buteo platypterus", 1, 4),
		rangeIntCase("Troglodytes aedon", 2, 8),
		rangeIntCase("Egretta thula", -1, 6),
		rangeIntCase("Corthylio calendula", -1, 12),
		rangeIntCase("Patagioenas fasciata", 1, 2),
		rangeIntCase("Ardea alba", -1, 6),
		rangeIntCase("Colinus virginianus", 8, 25),
		rangeIntCase("Sialia mexicana", 4, 6),
		rangeIntCase("Spizella pallida", 3, 5),
		rangeIntCase("Charadrius melodus", 2, 3),
		rangeIntCase("Spinus pinus", 0, 0))
}

func TestClutchSize_WhatBird(t *testing.T) {
	whatBirdClutchSize := func(englishName string) *inference.IntRange {
		return whatBirdRequestForTesting(englishName).propertySearchers().clutchSize.ClutchSize()
	}
	testRangeIntCases(t, whatBirdClutchSize,
		rangeIntCase("Striped Owl", 2, 4),
		rangeIntCase("Reed Bunting", 4, 6),
		rangeIntCase("Red-breasted Sapsucker", 4, 5),
		rangeIntCase("Thick-billed Murre", 1, 1),
		rangeIntCase("Thick-billed Parrot", 1, 4),
		rangeIntCase("Zenaida Dove", 2, 2))
}

func TestClutchSize_Audubon(t *testing.T) {
	audubonClutchSize := func(englishName string) *inference.IntRange {
		return audubonRequestForTesting(englishName).propertySearchers().clutchSize.ClutchSize()
	}
	testRangeIntCases(t, audubonClutchSize,
		rangeIntCase("Blue-gray gnatcatcher", 4, 5),
		rangeIntCase("Bushtit", 5, 7),
		rangeIntCase("Lucy's warbler", 4, 5),
		rangeIntCase("Red-faced warbler", 3, 4))
}

func TestClutchSize_AllAboutBirds(t *testing.T) {
	allAboutBirdsClutchSize := func(englishName string) *inference.IntRange {
		return allAboutBirdsRequestForTesting(englishName).propertySearchers().clutchSize.ClutchSize()
	}
	testRangeIntCases(t, allAboutBirdsClutchSize,
		rangeIntCase("Black-throated sparrow", 2, 5),
		rangeIntCase("Black-and-white warbler", 4, 6),
		rangeIntCase("Clay-colored sparrow", 0, 0),
		rangeIntCase("Downy woodpecker", 3, 8),
		rangeIntCase("Hermit thrush", 3, 6))
}
