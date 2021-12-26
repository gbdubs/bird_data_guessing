package bird_data_guessing

import (
	"fmt"

	"github.com/gbdubs/amass"
)

func (input *Input) Execute() (*Output, error) {
	oo := &Output{}

	requests := make([]*amass.GetRequest, 0)
	for i, bird := range input.Names {
		// The compiler doesn't allow two variadic args of the same type
		isMemoized, birdData, err := readMemoized(bird)
		if err != nil {
			return oo, fmt.Errorf("While reading is memoized for %s: %v", bird.English, err)
		}
		if isMemoized {
			input.VLog("[%d/%d] Memoized Read %s\n", i, len(input.Names), bird.English)
			oo.BirdData = append(oo.BirdData, birdData)
			continue
		}
		input.VLog("[%d/%d] Created Requests for %s\n", i, len(input.Names), bird.English)
		requests = append(requests, createWikipediaRequests(bird)...)
		requests = append(requests, createAudubonRequests(bird)...)
		requests = append(requests, createAllAboutBirdsRequests(bird)...)
		requests = append(requests, createWhatBirdRequests(bird)...)
		requests = append(requests, createRSPBRequests(bird)...)
	}

	amasser := amass.Amasser{
		TotalMaxConcurrentRequests: 20,
		Verbose:                    input.VIndent(),
		AllowedErrorProportion:     0.10,
	}
	responses, err := amasser.GetAll(requests)
	if err != nil {
		return oo, fmt.Errorf("While amassing results: %v", err)
	}

	latinToWikipedia := reconstructWikipediaResponsesKeyedByLatinName(responses)
	latinToAllAboutBirds := reconstructAllAboutBirdsResponsesKeyedByLatinName(responses)
	latinToAudubon := reconstructAudubonResponsesKeyedByLatinName(responses)
	latinToWhatBird := reconstructWhatBirdsResponsesKeyedByLatinName(responses)
	latinToRspb := reconstructRSPBResponsesKeyedByLatinName(responses)

	for i, bird := range input.Names {
		input.VLog("[%d/%d] Collecting + Merging %s", i, len(input.Names), bird.English)
		latin := bird.Latin
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
		if w, ok := latinToRspb[latin]; ok {
			allSources = append(allSources, w.propertySearchers().getData(bird))
		}
		if len(allSources) == 0 {
			input.VLog(" - EMPTY. Continuing\n")
			continue
		}
		merged, highConfidence := mergeSources(allSources)
		merged.Name = bird
		if highConfidence {
			oo.BirdData = append(oo.BirdData, *merged)
			err := writeMemoized(*merged)
			if err != nil {
				return oo, fmt.Errorf("While merging bird %s: %v", bird.English, err)
			}
			input.VLog(" - merged + memoized.\n")
		} else {
			input.VLog(" - low confidence, not merged or memoized\n")
		}
	}
	return oo, nil
}
