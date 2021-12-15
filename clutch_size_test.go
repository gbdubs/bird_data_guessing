package bird_data_guessing

import (
	"testing"

	"github.com/gbdubs/inference"
)

func TestClutchSize_Wikipedia(t *testing.T) {
	wikipediaClutchSize := func(n string) (*inference.IntRange, error) {
		r, err := getWikipediaResponse(n)
		if err != nil {
			return &inference.IntRange{}, err
		}
		s := &r.propertySearchers().clutchSize
		return s.ClutchSize(), err
	}
	runIntRangeTests(t, wikipediaClutchSize,
		intRangeCase("Aythya valisineria", 5, 11),
		intRangeCase("Cardinalis sinuatus", 2, 4),
		intRangeCase("Cistothorus palustris", 4, 6),
		intRangeCase("Passer domesticus", 6, 6),
		intRangeCase("Pooecetes gramineus", 3, 5),
		intRangeCase("Loxia leucoptera", 3, 5),
		intRangeCase("Spinus pinus", 0, 0))
}

/*
func TestClutchSize_WhatBird(t *testing.T) {
	whatBirdClutchSize := func(n string) (*Property, error) {
		r, err := getWhatBirdResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().clutchSize
		return s.ClutchSize(), err
	}
	runIntTests(t, whatBirdClutchSize,
		intCase("Striped Owl", 3),
		intCase("Reed Bunting", 5),
		intCase("Red-breasted Sapsucker", 4),
		intCase("Thick-billed Murre", 1),
		intCase("Thick-billed Parrot", 2),
		intCase("Zenaida Dove", 2))
}

func TestClutchSize_Wikipedia(t *testing.T) {
	wikipediaClutchSize := func(n string) (*Property, error) {
		r, err := getWikipediaResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().clutchSize
		return s.ClutchSize(), err
	}
	runIntTests(t, wikipediaClutchSize,
		intCase("Aythya valisineria", 8),
		intCase("Cardinalis sinuatus", 3),
		intCase("Cistothorus palustris", 5),
		intCase("Passer domesticus", 4),
		intCase("Pooecetes gramineus", 4),
		intCase("Loxia leucoptera", 4),
		intCase("Spinus pinus", 0))
}

func TestClutchSize_Audubon(t *testing.T) {
	audubonClutchSize := func(n string) (*Property, error) {
		r, err := getAudubonResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().clutchSize
		return s.ClutchSize(), err
	}
	runIntTests(t, audubonClutchSize,
		intCase("Blue-gray gnatcatcher", 4),
		intCase("Bushtit", 6),
		intCase("Lucy's warbler", 4),
		intCase("Red-faced warbler", 3))
}

func TestClutchSize_AllAboutBirds(t *testing.T) {
	allAboutBirdsClutchSize := func(n string) (*Property, error) {
		r, err := getAllAboutBirdsResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().clutchSize
		return s.ClutchSize(), err
	}
	runIntTests(t, allAboutBirdsClutchSize,
		intCase("Black-throated sparrow", 3),
		intCase("Black-and-white warbler", 5),
		intCase("Clay-colored sparrow", 0),
		intCase("Downy woodpecker", 5),
		intCase("Hermit thrush", 4))
}
*/
