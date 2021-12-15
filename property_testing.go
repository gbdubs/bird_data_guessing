package bird_data_guessing

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gbdubs/inference"
)

type stringTestCase struct {
	name     string
	expected string
}

type intTestCase struct {
	name     string
	expected int
}

type intRangeTestCase struct {
	name        string
	expectedMin int
	expectedMax int
}

type float64RangesTestCase struct {
	name     string
	expected [][2]float64
}

func stringCase(n string, e string) stringTestCase {
	return stringTestCase{
		name:     n,
		expected: e,
	}
}

func intCase(n string, e int) intTestCase {
	return intTestCase{
		name:     n,
		expected: e,
	}
}

func intRangeCase(n string, min int, max int) intRangeTestCase {
	return intRangeTestCase{
		name:        n,
		expectedMin: min,
		expectedMax: max,
	}
}

func float64RangesCase(n string, vs ...float64) float64RangesTestCase {
	e := make([][2]float64, len(vs)/2)
	for i := 0; i < len(vs)/2; i++ {
		e[i] = [2]float64{vs[2*i], vs[2*i+1]}
	}
	return float64RangesTestCase{
		name:     n,
		expected: e,
	}
}

type testIntRangeBehavior func(englishOrLatinName string) (*inference.IntRange, error)

type testStringBehavior func(englishOrLatinName string) (*inference.String, error)

func runStringTests(t *testing.T, b testStringBehavior, testCases ...stringTestCase) {
	for _, c := range testCases {
		p, err := b(c.name)
		if err != nil {
			t.Errorf(`case '%s' yielded error: %+v`, c.name, err)
			continue
		}
		if c.expected != p.Value {
			t.Errorf(`case %s assertion error - expected '%s' - was '%s' - debug: %+v`, c.name, c.expected, p.Value, p)
		}
	}
}

func runIntRangeTests(t *testing.T, b testIntRangeBehavior, testCases ...intRangeTestCase) {
	for _, c := range testCases {
		r, err := b(c.name)
		if err != nil {
			t.Errorf(`case '%s' yielded error: %+v`, c.name, err)
			continue
		}
		if c.expectedMin != r.Min || c.expectedMax != r.Max {
			t.Errorf(`case %s assertion error - expected [%d, %d] - was [%d, %d] - debug: %+v %s`, c.name, c.expectedMin, c.expectedMax, r.Min, r.Max, r, r.Source.Dump())
		}
	}
}

type float64Ranges [][2]float64

func (f float64Ranges) Len() int      { return len(f) }
func (f float64Ranges) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f float64Ranges) Less(i, j int) bool {
	if f[i][0] < f[j][0] {
		return true
	}
	if f[i][0] > f[j][0] {
		return false
	}
	if f[i][1] < f[j][1] {
		return true
	}
	if f[i][1] > f[j][1] {
		return false
	}
	return false
}

type testFloat64RangesBehavior func(englishOrLatinName string) ([]*inference.Float64Range, error)

func runFloat64RangesTests(t *testing.T, b testFloat64RangesBehavior, testCases ...float64RangesTestCase) {
	for _, c := range testCases {
		r, err := b(c.name)
		if err != nil {
			t.Errorf(`case '%s' yielded error: %+v`, c.name, err)
			continue
		}
		aSorted := make([][2]float64, len(r))
		sourceDump := ""
		for i, rage := range r {
			aSorted[i] = [2]float64{rage.Min, rage.Max}
			sourceDump += rage.Source.Dump()
		}
		eSorted := c.expected
		sort.Sort(float64Ranges(eSorted))
		sort.Sort(float64Ranges(aSorted))

		if !reflect.DeepEqual(eSorted, aSorted) {
			t.Errorf(`case %s assertion error - expected '%v' - was '%v' - debug %+v %s`, c.name, eSorted, aSorted, r, sourceDump)
		}
	}
}
