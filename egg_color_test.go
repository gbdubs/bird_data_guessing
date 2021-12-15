package bird_data_guessing

import (
	"testing"

	"github.com/gbdubs/inference"
)

func TestEggColor_ZZZWikipedia(t *testing.T) {
	wikipediaEggColor := func(n string) (*inference.String, error) {
		r, err := getWikipediaResponse(n)
		if err != nil {
			return &inference.String{}, err
		}
		s := &r.propertySearchers().eggColor
		return s.ZZZEggColor(), err
	}
	runStringTestsZZZ(t, wikipediaEggColor,
		stringCaseZZZ("Cardinalis sinuatus", "whitish with specks of green or gray"),
		stringCaseZZZ("Gracula religiosa", ""),
		stringCaseZZZ("Platalea ajaja", "whitish with brown markings"),
		stringCaseZZZ("Dryobates scalaris", "plain white"),
		stringCaseZZZ("Grus americana", "olive"),
		stringCaseZZZ("Ictinia mississippiensis", "white to pale-bluish in color, and are usually about 1"),
		stringCaseZZZ("Falco mexicanus", "subelliptical and pinkish with brown, reddish-brown, and purplish dots"),
		stringCaseZZZ("Ectopistes migratorius", "white and oval shaped and averaged 40 by 34 mm (1"),
		stringCaseZZZ("Leiothlypis virginiae", "white in color and dotted with fine brown speckles"),
		stringCaseZZZ("Motacilla alba", "cream-coloured, often with a faint bluish-green or turquoise tint, and heavily spotted with reddish brown; they measure, on average, 21 mm × 15 mm (0"),
		stringCaseZZZ("Numenius americanus", "vary in hue from white to olive"),
	)
}

/*
func TestEggColor_WhatBird(t *testing.T) {
	whatBirdEggColor := func(n string) (*Property, error) {
		r, err := getWhatBirdResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().eggColor
		return s.EggColor(), nil
	}
	runStringTests(t, whatBirdEggColor,
		stringCase("Ancient murrelet", "Light to olive brown, sometimes speckled."),
		stringCase("Anna's Hummingbird", "White"),
		stringCase("Audubon's Oriole", "Pale blue or gray with brown or purple marks"),
		stringCase("Black Rail", "Pale pink to white with brown spots"),
		stringCase("Black-crested titmouse", "White with brown spots"),
		stringCase("Brown booby", "White to pale blue green, nest stained."))
}

func TestEggColor_Wikipedia(t *testing.T) {
	wikipediaEggColor := func(n string) (*Property, error) {
		r, err := getWikipediaResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().eggColor
		return s.EggColor(), err
	}
	runStringTests(t, wikipediaEggColor,
		stringCase("Cardinalis sinuatus", "whitish with specks of green or gray"),
		stringCase("Gracula religiosa", ""),
		stringCase("Platalea ajaja", "whitish with brown markings"),
		stringCase("Dryobates scalaris", "plain white"),
		stringCase("Grus americana", "olive"),
		stringCase("Ictinia mississippiensis", "white to pale-bluish in color, and are usually about 1"),
		stringCase("Falco mexicanus", "subelliptical and pinkish with brown, reddish-brown, and purplish dots"),
		stringCase("Ectopistes migratorius", "white and oval shaped and averaged 40 by 34 mm (1"),
		stringCase("Leiothlypis virginiae", "white in color and dotted with fine brown speckles"),
		stringCase("Motacilla alba", "cream-coloured, often with a faint bluish-green or turquoise tint, and heavily spotted with reddish brown; they measure, on average, 21 mm × 15 mm (0"),
		stringCase("Numenius americanus", "vary in hue from white to olive"),
	)
}

func TestEggColor_Audubon(t *testing.T) {
	audubonEggColor := func(n string) (*Property, error) {
		r, err := getAudubonResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().eggColor
		return s.EggColor(), err
	}
	runStringTests(t, audubonEggColor,
		stringCase("Great blue heron", "Pale blue"),
		stringCase("Bald eagle", "White"),
		stringCase("Mountain chickadee", "White, dotted with reddish brown, sometimes unmarked"),
		stringCase("Bridled Titmouse", "Unmarked white"),
		stringCase("American Robin", `Pale blue or "robin's-egg blue`),
		stringCase("Boreal chickadee", "White, with fine reddish brown dots often concentrated at larger end"),
		stringCase("Black and white warbler", "Creamy white, flecked with brown at large end"),
	)
}

func TestEggColor_AllAboutBirds(t *testing.T) {
	allAboutBirdsEggColor := func(n string) (*Property, error) {
		r, err := getAllAboutBirdsResponse(n)
		if err != nil {
			return &Property{}, err
		}
		s := &r.propertySearchers().eggColor
		return s.EggColor(), err
	}
	runStringTests(t, allAboutBirdsEggColor,
		stringCase("Black-and-white warbler", "Creamy white, pale bluish- or greenish-white, with speckles of brown or lavender"),
		stringCase("Prothonotary warbler", "White spotted with rust-brown to lavender"),
		stringCase("Anhinga", "Conspicuously pointed at one end, pale bluish green, and overlaid with a chalky coating"),
	)
}*/
