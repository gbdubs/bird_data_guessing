package bird_data_guessing

import (
	"sort"
	"testing"

	"github.com/gbdubs/inference"
)

// float64

type testFloat64Case struct {
	name     string
	expected float64
}

type testFloat64Behavior func(englishOrLatinName string) *inference.Float64

func float64Case(englishOrLatinName string, expectedResult float64) testFloat64Case {
	return testFloat64Case{
		name:     englishOrLatinName,
		expected: expectedResult,
	}
}

func testFloat64Cases(t *testing.T, b testFloat64Behavior, cases ...testFloat64Case) {
	for _, c := range cases {
		actual := b(c.name)
		if c.expected != actual.Value {
			t.Errorf(
				`case %s assertion error: expected '%v', was '%v'. Debug %+v source: %s`,
				c.name, c.expected, actual.Value, actual, actual.Source.Dump())
		}
	}
}

// []Range<Float64>

type testSliceRangeFloat64Case struct {
	name            string
	expected        [][2]float64
	expectedOrdered bool
}

type testSliceRangeFloat64Behavior func(englishOrLatinName string) []*inference.Float64Range

func unorderedSliceRangeFloat64Case(englishOrLatinName string, vs ...float64) testSliceRangeFloat64Case {
	e := make([][2]float64, len(vs)/2)
	for i := 0; i < len(vs)/2; i++ {
		e[i] = [2]float64{vs[2*i], vs[2*i+1]}
	}
	return testSliceRangeFloat64Case{
		name:            englishOrLatinName,
		expected:        e,
		expectedOrdered: false,
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

func testSliceRangeFloat64Cases(t *testing.T, b testSliceRangeFloat64Behavior, cases ...testSliceRangeFloat64Case) {
	for _, c := range cases {
		actual := b(c.name)
		if len(actual) != len(c.expected) {
			t.Errorf(`case %s assertion error - lengths differ - expected len '%v', was '%v'. Debug %+v`,
				c.name, len(c.expected), len(actual), actual)
			continue
		}
		as := make([][2]float64, len(actual))
		es := make([][2]float64, len(actual))
		for i, a := range actual {
			as[i] = [2]float64{a.Min, a.Max}
			es[i] = c.expected[i]
		}
		if !c.expectedOrdered {
			sort.Sort(float64Ranges(as))
			sort.Sort(float64Ranges(es))
		}
		errorIndexes := make([]int, 0)
		for i, _ := range as {
			if as[i][0] != es[i][0] || as[i][1] != es[i][1] {
				errorIndexes = append(errorIndexes, i)
			}
		}
		if len(errorIndexes) > 0 {
			t.Errorf(
				`case %s assertion error - errors at indexes %v - expected '%+v' - was '%+v' - debug: %+v`,
				c.name, errorIndexes, es, as, actual)
		}
	}
}
