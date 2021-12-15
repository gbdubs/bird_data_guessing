package bird_data_guessing

import (
	"testing"

	"github.com/gbdubs/inference"
)

func TestWingspan_Wikipedia(t *testing.T) {
	wikipediaWingspan := func(n string) ([]*inference.Float64Range, error) {
		r, err := getWikipediaResponse(n)
		if err != nil {
			return []*inference.Float64Range{}, err
		}
		s := &r.propertySearchers().wingspan
		return s.Wingspan(), err
	}
	runFloat64RangesTests(t, wikipediaWingspan,
		float64RangesCase("Aythya valisineria", 79.0, 89.0),
		float64RangesCase("Cardinalis sinuatus", 21.082, 21.082),
		float64RangesCase("Cistothorus palustris", 14.986, 14.986),
		float64RangesCase("Passer domesticus", 19.0, 25),
		float64RangesCase("Pooecetes gramineus", 23.876, 23.876),
		float64RangesCase("Loxia leucoptera", 26, 29),
		float64RangesCase("Columba livia", 60.96, 71.12),
		float64RangesCase("Riparia riparia", 22.86, 33.02),
		float64RangesCase("Spinus pinus", 18, 22),
		float64RangesCase("Junco hyemalis", 18, 25),
		float64RangesCase("Accipiter striatus", 58, 68, 42, 58))
}

/*
func TestWingspan_ZZZWhatBird(t *testing.T) {
	whatBirdWingspan := func(n string) (*Property, error) {
		r, err := getWhatBirdResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().wingspan
		return s.Wingspan(), err
	}
	runIntTests(t, whatBirdWingspan,
		intCase("Black-backed Woodpecker", 43),
		intCase("Allen's Hummingbird", 12),
		intCase("Turkey vulture", 178),
		intCase("Bald eagle", 214),
		intCase("Striped Owl", 87),
	)
}
*/

/*
func TestWingspan_WhatBird(t *testing.T) {
	whatBirdWingspan := func(n string) (*Property, error) {
		r, err := getWhatBirdResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().wingspan
		return s.Wingspan(), err
	}
	runIntTests(t, whatBirdWingspan,
		intCase("Black-backed Woodpecker", 43),
		intCase("Allen's Hummingbird", 12),
		intCase("Turkey vulture", 178),
		intCase("Bald eagle", 214),
		intCase("Striped Owl", 87),
	)
}

func TestWingspan_Wikipedia(t *testing.T) {
	wikipediaWingspan := func(n string) (*Property, error) {
		r, err := getWikipediaResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().wingspan
		return s.Wingspan(), err
	}
	runIntTests(t, wikipediaWingspan,
		intCase("Aythya valisineria", 84),
		intCase("Cardinalis sinuatus", 21),
		intCase("Cistothorus palustris", 15),
		intCase("Passer domesticus", 22),
		intCase("Pooecetes gramineus", 24),
		intCase("Loxia leucoptera", 28),
		intCase("Columba livia", 66),
		intCase("Riparia riparia", 28),
		intCase("Spinus pinus", 20),
		intCase("Junco hyemalis", 22),
		intCase("Accipiter striatus", 50))
}
*/
// Audubon and AllAboutBirds don't provide reliable Wingspan information.
