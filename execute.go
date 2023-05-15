package bird_data_guessing

import (
	"fmt"
	"sync"
	"time"

	"github.com/gbdubs/amass"
)

func (input *Input) Execute() (*Output, error) {
	input.VLog("Bird data guessing: starting.\n")
	oo := &Output{}

	wasMemoizedLatins := make(map[string]bool)
	requests := make([]*amass.GetRequest, 0)
	requestsDone := 0
	var requestsErr error = nil
	requestsLock := sync.RWMutex{}
	for index, b := range input.Names {
		i := index
		bird := b
		go func() {
			isMemoized, birdData, err := readMemoized(bird)
			if err != nil {
				requestsLock.Lock()
				requestsErr = fmt.Errorf("While reading is memoized for %s: %v", bird.English, err)
				requestsDone++
				requestsLock.Unlock()
				return
			}
			if isMemoized {
				requestsLock.Lock()
				wasMemoizedLatins[bird.Latin] = true
				requestsDone++
				input.VLog("[%d/%d] Memoized Read %s\n", requestsDone, len(input.Names), bird.English)
				oo.BirdData = append(oo.BirdData, birdData)
				requestsLock.Unlock()
			} else {
				wikipedia := createWikipediaRequests(bird)
				audubon := createAudubonRequests(bird)
				allAboutBirds := createAllAboutBirdsRequests(bird)
				rspb := createRSPBRequests(bird)
				// Whatbird is currently down
				// whatBird := createWhatBirdRequests(bird)

				requestsLock.Lock()
				input.VLog("[%d/%d] Created Requests for %s\n", i, len(input.Names), bird.English)
				// The compiler doesn't allow mulitple variadic args
				requests = append(requests, wikipedia...)
				requests = append(requests, audubon...)
				requests = append(requests, allAboutBirds...)
				requests = append(requests, rspb...)
				// requests = append(requests, whatBird...)
				requestsDone++
				requestsLock.Unlock()
			}
		}()
	}

	for requestsDone < len(input.Names) {
		if requestsErr != nil {
			return oo, requestsErr
		}
		input.VLog("[%d/%d] Waiting for all requests to be generated.\n", requestsDone, len(input.Names))
		time.Sleep(100 * time.Millisecond)
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
	englishToAllAboutBirds := reconstructAllAboutBirdsResponsesKeyedByEnglishName(responses)
	englishToAudubon := reconstructAudubonResponsesKeyedByEnglishName(responses)
	englishToWhatBird := reconstructWhatBirdsResponsesKeyedByEnglishName(responses)
	englishToRspb := reconstructRSPBResponsesKeyedByEnglishName(responses)

	for i, bird := range input.Names {
		input.VLog("[%d/%d] Collecting + merging ", i, len(input.Names))
		latin := bird.Latin
		english := bird.English
		if wasMemoizedLatins[latin] {
			input.VLog(" %s - came from memoized.\n", bird.English)
			continue
		}
		allSources := make([]*singleSourceData, 0)
		if w, ok := latinToWikipedia[latin]; ok {
			allSources = append(allSources, w.propertySearchers().getData(bird))
		}
		if a, ok := englishToAllAboutBirds[english]; ok {
			allSources = append(allSources, a.propertySearchers().getData(bird))
		}
		if a, ok := englishToAudubon[english]; ok {
			allSources = append(allSources, a.propertySearchers().getData(bird))
		}
		if w, ok := englishToWhatBird[english]; ok {
			allSources = append(allSources, w.propertySearchers().getData(bird))
		}
		if w, ok := englishToRspb[english]; ok {
			allSources = append(allSources, w.propertySearchers().getData(bird))
		}
		if len(allSources) == 0 {
			input.VLog(" - EMPTY. Continuing\n")
			continue
		}
		merged, highConfidence := input.mergeSources(allSources)
		merged.Name = bird
		if highConfidence {
			oo.BirdData = append(oo.BirdData, *merged)
			err := writeMemoized(*merged)
			if err != nil {
				return oo, fmt.Errorf("While merging bird %s: %v", bird.English, err)
			}
			input.VLog(" = merged + memoized %s.\n", bird.English)
		} else {
			input.VLog(" = LC!, not merged or memoized: %s.\n", bird.English)
		}
	}
	return oo, nil
}
