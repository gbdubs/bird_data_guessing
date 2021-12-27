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
		unorderedSliceRangeFloat64Case("Cistothorus palustris", 15, 15),
		unorderedSliceRangeFloat64Case("Passer domesticus", 19.0, 25),
		unorderedSliceRangeFloat64Case("Pooecetes gramineus", 24, 24),
		unorderedSliceRangeFloat64Case("Loxia leucoptera", 26, 29),
		unorderedSliceRangeFloat64Case("Columba livia", 62, 72),
		unorderedSliceRangeFloat64Case("Riparia riparia", 25, 33),
		unorderedSliceRangeFloat64Case("Spinus pinus", 18, 22),
		unorderedSliceRangeFloat64Case("Junco hyemalis", 18, 25),
		unorderedSliceRangeFloat64Case("Ictinia mississippiensis", 91, 91, 91.44, 91.44),
		unorderedSliceRangeFloat64Case("Bubulcus ibis", 88, 96),
		unorderedSliceRangeFloat64Case("Spizella pallida", 19, 19),
		unorderedSliceRangeFloat64Case("Euphagus cyanocephalus", 39, 39),
		unorderedSliceRangeFloat64Case("Melanitta americana", 71, 71),
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
		unorderedSliceRangeFloat64Case("Bald eagle", 204, 204),
		unorderedSliceRangeFloat64Case("Brandt's cormorant", 108, 108),
	)
}

func TestWingspan_RSPB(t *testing.T) {
	rspbWingspan := func(englishName string) []*inference.Float64Range {
		return rspbRequestForTesting(englishName).propertySearchers().wingspan.Wingspan()
	}
	testSliceRangeFloat64Cases(t, rspbWingspan,
		unorderedSliceRangeFloat64Case("Stonechat", 18, 21),
		unorderedSliceRangeFloat64Case("Temminck's Stint", 34, 37),
		unorderedSliceRangeFloat64Case("Stone curlew", 77, 85),
	)
}

// Audubon does not have reliable wingspan information.
func TestWingspan_Audubon_NoContentInSearcher(t *testing.T) {
	searcherText := audubonRequestForTesting("Bald eagle").propertySearchers().wingspan.text
	if searcherText != "" {
		t.Errorf("Expected Audubon Wingspan to be empty, but was %s.", searcherText)
	}
}
