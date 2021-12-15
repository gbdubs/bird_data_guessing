package bird_data_guessing

import (
	"fmt"

	"github.com/gbdubs/inference"
)

func (r *searcher) EggColor() []*inference.String {
	colors := "(white|whitish|cream|olive|brown|blue|pink|pinkish|redish|purple|green|tan|light)"
	connectors := "(which are|are|were|was|which)"

	patterns := make(map[string]int)
	patterns[fmt.Sprintf(`eggs?,? %s (([^.]{0,30} )?%s[^.]*).`, connectors, colors)] = 2
	patterns[fmt.Sprintf("(%s).(colou?r(ed)?) egg", colors)] = 1
	// Uniuqe to Audubon
	patterns[fmt.Sprintf(`EggsAudubonEggs [^.]+\. ([^.]*%s[^.]*)\.`, colors)] = 1
	// Unique to WhatBird
	patterns["Egg Color: (.+) Number of Eggs:"] = 1
	// For All About Birds
	patterns[`Egg Description:([^.]+)\.`] = 1

	return r.ZZZExtractAllMatches(patterns)
}
