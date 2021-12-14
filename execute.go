package bird_data_guessing

import (
	"log"
	"math"
	"strings"

	"github.com/gbdubs/attributions"
)

func (i *Input) Execute() (*Output, error) {
	oo := Output{}
	o := &oo
	dd := DebugDatas{}
	var datas []Data
	var attribs []attributions.Attribution
	englishName := i.EnglishName

	wResp, wErr := getWikipediaResponse(i.LatinName)
	if wErr == nil {
		englishName = wResp.englishName()
		wData, wDebug := wResp.propertySearchers().getData(englishName)
		dd.Wikipedia = *wDebug
		datas = append(datas, *wData)
		attribs = append(attribs, *wResp.attribution())
	} else if strings.Contains(wErr.Error(), "404") {
		// Nothing to do.
	} else {
		log.Printf("Wikipedia Error %+v", wErr)
	}

	aabResp, aabErr := getAllAboutBirdsResponse(englishName)
	if aabErr == nil {
		aabData, aabDebug := aabResp.propertySearchers().getData(englishName)
		dd.AllAboutBirds = *aabDebug
		datas = append(datas, *aabData)
		attribs = append(attribs, *aabResp.attribution())
	} else if strings.Contains(aabErr.Error(), "404") {
		// Nothing to do.
	} else {
		log.Printf("All About Birds Error %+v", aabErr)
	}

	auResp, auErr := getAudubonResponse(englishName)
	if auErr == nil {
		auData, auDebug := auResp.propertySearchers().getData(englishName)
		dd.Audubon = *auDebug
		datas = append(datas, *auData)
		attribs = append(attribs, *auResp.attribution())
	} else if strings.Contains(auErr.Error(), "404") {
		// Nothing to do.
	} else {
		log.Printf("Audubon Error %+v", auErr)
	}

	wbResp, wbErr := getWhatBirdResponse(englishName)
	if wbErr == nil {
		wbData, wbDebug := wbResp.propertySearchers().getData(englishName)
		dd.WhatBird = *wbDebug
		datas = append(datas, *wbData)
		attribs = append(attribs, *wbResp.attribution())
	} else if strings.Contains(wbErr.Error(), "404") {
		// Nothing to do.
	} else {
		log.Printf("WhatBird Error %+v", auErr)
	}

	if len(datas) == 0 {
		return o, wErr
	}

	o.Attributions = attribs
	o.Data = synthesizeDatas(i.LatinName, datas)
	if i.Debug {
		o.DebugDatas = dd
	}
	return o, nil
}

type dataStringField func(d Data) string
type dataIntField func(d Data) int

func synthesizeDatas(latinName string, ds []Data) Data {
	return Data{
		LatinName:     latinName,
		EnglishName:   first(ds, func(d Data) string { return d.EnglishName }),
		WheatScore:    mult(ds, func(d Data) int { return d.WheatScore }),
		WormScore:     mult(ds, func(d Data) int { return d.WormScore }),
		BerryScore:    mult(ds, func(d Data) int { return d.BerryScore }),
		RatScore:      mult(ds, func(d Data) int { return d.RatScore }),
		FishScore:     mult(ds, func(d Data) int { return d.FishScore }),
		NectarScore:   mult(ds, func(d Data) int { return d.NectarScore }),
		ForestScore:   mult(ds, func(d Data) int { return d.ForestScore }),
		GrassScore:    mult(ds, func(d Data) int { return d.GrassScore }),
		WaterScore:    mult(ds, func(d Data) int { return d.WaterScore }),
		CupScore:      mult(ds, func(d Data) int { return d.CupScore }),
		GroundScore:   mult(ds, func(d Data) int { return d.GroundScore }),
		PlatformScore: mult(ds, func(d Data) int { return d.PlatformScore }),
		SlotScore:     mult(ds, func(d Data) int { return d.SlotScore }),
		Wingspan:      geomAvg(ds, func(d Data) int { return d.Wingspan }),
		ClutchSize:    geomAvg(ds, func(d Data) int { return d.ClutchSize }),
		FlockingScore: mult(ds, func(d Data) int { return d.FlockingScore }),
		PredatorScore: mult(ds, func(d Data) int { return d.PredatorScore }),
		FunFact:       last(ds, func(d Data) string { return d.FunFact }),
		EggColor:      last(ds, func(d Data) string { return d.EggColor }),
	}
}

func last(ds []Data, f dataStringField) string {
	for i, _ := range ds {
		d := ds[len(ds)-1-i]
		s := f(d)
		if s != "" {
			return s
		}
	}
	return ""
}

func first(ds []Data, f dataStringField) string {
	for _, d := range ds {
		s := f(d)
		if s != "" {
			return s
		}
	}
	return ""
}

func geomAvg(ds []Data, f dataIntField) int {
	r := 1
	count := 0
	for _, d := range ds {
		v := f(d)
		if v == 0 {
			continue
		}
		r *= v + 1
		count++
	}
	if count < 2 {
		return r - 1
	}
	root := math.Pow(float64(r), 1.0/float64(count))
	result := int(math.Round(root) - 1)
	return result
}

func mult(ds []Data, f dataIntField) int {
	r := 1
	for _, d := range ds {
		r *= f(d) + 1
	}
	return r - 1
}
