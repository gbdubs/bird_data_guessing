package bird_data_guessing

import "testing"

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
		intCase("Spinus pinus", 20),
		intCase("Junco hyemalis", 22))
}

// Audubon and AllAboutBirds don't provide reliable Wingspan information.
