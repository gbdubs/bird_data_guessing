package bird_data_guessing

import (
	"sort"
	"testing"

	"github.com/gbdubs/inference"
)

// string

type testStringCase struct {
	name     string
	expected string
}

type testStringBehavior func(englishOrLatinName string) *inference.String

func stringCase(englishOrLatinName string, expectedResult string) testStringCase {
	return testStringCase{
		name:     englishOrLatinName,
		expected: expectedResult,
	}
}

func testStringCases(t *testing.T, b testStringBehavior, cases ...testStringCase) {
	for _, c := range cases {
		actual := b(c.name)
		if c.expected != actual.Value {
			t.Errorf(
				`case %s assertion error: expected '%s', was '%s'. Debug %+v source: %s`,
				c.name, c.expected, actual.Value, actual, actual.Source.Dump())
		}
	}
}

// []string

type testSliceStringCase struct {
	name            string
	expected        []string
	expectedOrdered bool
}

type testSliceStringBehavior func(englishOrLatinName string) []*inference.String

func unorderedSliceStringCase(englishOrLatinName string, expectedResult ...string) testSliceStringCase {
	return testSliceStringCase{
		name:            englishOrLatinName,
		expected:        expectedResult,
		expectedOrdered: false,
	}
}

func testSliceStringCases(t *testing.T, b testSliceStringBehavior, cases ...testSliceStringCase) {
	for _, c := range cases {
		actual := b(c.name)
		if len(actual) != len(c.expected) {
			t.Errorf(`case %s assertion error - lengths differ - expected len '%v', was '%v'. Debug %+v`,
				c.name, len(c.expected), len(actual), actual)
		}
		as := make([]string, len(actual))
		es := make([]string, len(actual))
		for i, a := range actual {
			as[i] = a.Value
			es[i] = c.expected[i]
		}
		if !c.expectedOrdered {
			sort.Strings(as)
			sort.Strings(es)
		}
		errorIndexes := make([]int, 0)
		for i, _ := range as {
			if as[i] != es[i] {
				errorIndexes = append(errorIndexes, i)
			}
		}
		if len(errorIndexes) > 0 {
			t.Errorf(
				`case %s assertion error - errors at indexes %v - expected '%s' - was '%s' - debug: %+v`,
				c.name, errorIndexes, es, as, actual)
		}
	}
}
