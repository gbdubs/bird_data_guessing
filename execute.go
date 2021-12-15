package bird_data_guessing

import "github.com/gbdubs/amass"

func (i *Input) Execute() (*Output, error) {
	oo := &Output{}

	requests := make([]*amass.GetRequest, 0)
	for _, bird := range i.Birds {
		// The compiler doesn't allow two variadic args of the same type
		requests = append(
			requests,
			createWikipediaRequest(bird),
			createAudubonRequest(bird))
		requests = append(requests,
			createAllAboutBirdsRequests(bird)...)
		requests = append(requests, createWhatBirdRequests(bird)...)
	}

	amasser := amass.Amasser{
		TotalMaxConcurrentRequests: 10,
		Verbose:                    true,
		AllowedErrorProportion:     0.0,
	}
	responses, err := amasser.GetAll(requests)
	if err != nil {
		return oo, err
	}

	latinToWikipedia := reconstructWikipediaResponsesKeyedByLatinName(responses)
	latinToAllAboutBirds := reconstructAllAboutBirdsResponsesKeyedByLatinName(responses)
	latinToAudubon := reconstructAudubonResponsesKeyedByLatinName(responses)
	latinToWhatBird := reconstructWhatBirdsResponsesKeyedByLatinName(responses)

	for _, bird := range i.Birds {
		latin := bird.LatinName
		allSources := make([]*singleSourceData, 0)
		if w, ok := latinToWikipedia[latin]; ok {
			allSources = append(allSources, w.propertySearchers().getData(bird))
		}
		if a, ok := latinToAllAboutBirds[latin]; ok {
			allSources = append(allSources, a.propertySearchers().getData(bird))
		}
		if a, ok := latinToAudubon[latin]; ok {
			allSources = append(allSources, a.propertySearchers().getData(bird))
		}
		if w, ok := latinToWhatBird[latin]; ok {
			allSources = append(allSources, w.propertySearchers().getData(bird))
		}
		if len(allSources) == 0 {
			continue
		}
		merged, highConfidence := mergeSources(allSources)
		if highConfidence {
			oo.BirdData = append(oo.BirdData, merged)
		}
	}

	return oo, nil
}
