package bird_data_guessing

import (
	"testing"

	"github.com/gbdubs/inference"
)

type rangeIntTestCase struct {
	name        string
	expectedMin int
	expectedMax int
}

type testRangeIntBehavior func(englishOrLatinName string) *inference.IntRange

func rangeIntCase(n string, min int, max int) rangeIntTestCase {
	return rangeIntTestCase{
		name:        n,
		expectedMin: min,
		expectedMax: max,
	}
}

func testRangeIntCases(t *testing.T, b testRangeIntBehavior, testCases ...rangeIntTestCase) {
	for _, c := range testCases {
		r := b(c.name)
		if c.expectedMin != r.Min || c.expectedMax != r.Max {
			t.Errorf(`case %s assertion error - expected [%d, %d] - was [%d, %d] - debug: %+v %s`, c.name, c.expectedMin, c.expectedMax, r.Min, r.Max, r, r.Source.Dump())
		}
	}
}
