package bird_data_guessing

func (i *Input) Execute() (*Output, error) {
	oo := Output{}

	// TODO! :)

	return &oo, nil
	/*
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
	*/
}
