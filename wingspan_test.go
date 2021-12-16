package bird_data_guessing

import (
	"testing"

	"github.com/gbdubs/inference"
)

func TestWingspan_Wikipedia(t *testing.T) {
	wikipediaWingspan := func(latinName string) []*inference.Float64Range {
		return wikipediaRequestForTesting(latinName).propertySearchers().wingspan.Wingspan()
	}
	testSliceRangeFloat64Cases(t, wikipediaWingspan,
		unorderedSliceRangeFloat64Case("Aythya valisineria", 79.0, 89.0),
		unorderedSliceRangeFloat64Case("Cardinalis sinuatus", 21.082, 21.082),
		unorderedSliceRangeFloat64Case("Cistothorus palustris", 14.986, 14.986),
		unorderedSliceRangeFloat64Case("Passer domesticus", 19.0, 25),
		unorderedSliceRangeFloat64Case("Pooecetes gramineus", 23.876, 23.876),
		unorderedSliceRangeFloat64Case("Loxia leucoptera", 26, 29),
		unorderedSliceRangeFloat64Case("Columba livia", 60.96, 71.12),
		unorderedSliceRangeFloat64Case("Riparia riparia", 22.86, 33.02),
		unorderedSliceRangeFloat64Case("Spinus pinus", 18, 22),
		unorderedSliceRangeFloat64Case("Junco hyemalis", 18, 25),
		unorderedSliceRangeFloat64Case("Accipiter striatus", 58, 68, 42, 58))
}

func TestWingspan_WhatBird(t *testing.T) {
	whatBirdWingspan := func(englishName string) []*inference.Float64Range {
		return whatBirdRequestForTesting(englishName).propertySearchers().wingspan.Wingspan()
	}
	testSliceRangeFloat64Cases(t, whatBirdWingspan,
		unorderedSliceRangeFloat64Case("Black-backed Woodpecker", 43.0, 43.0),
		unorderedSliceRangeFloat64Case("Allen's Hummingbird", 12.0, 12.0),
		unorderedSliceRangeFloat64Case("Turkey vulture", 173.0, 183.0),
		unorderedSliceRangeFloat64Case("Bald eagle", 183.0, 244.0),
		unorderedSliceRangeFloat64Case("Striped Owl", 76.0, 97.0),
	)
}

func TestWingspan_AllAboutBirds(t *testing.T) {
	whatBirdWingspan := func(englishName string) []*inference.Float64Range {
		return allAboutBirdsRequestForTesting(englishName).propertySearchers().wingspan.Wingspan()
	}
	testSliceRangeFloat64Cases(t, whatBirdWingspan,
		unorderedSliceRangeFloat64Case("Sprague's Pipit"),
		unorderedSliceRangeFloat64Case("Bald eagle", 203.962, 203.962),
		unorderedSliceRangeFloat64Case("Brandt's cormorant", 107.95, 107.95),
	)
}

// Audubon does not have reliable wingspan information.
func TestWingspan_Audubon_NoContentInSearcher(t *testing.T) {
	searcherText := audubonRequestForTesting("Bald eagle").propertySearchers().wingspan.text
	if searcherText != "" {
		t.Errorf("Expected Audubon Wingspan to be empty, but was %s.", searcherText)
	}
}
