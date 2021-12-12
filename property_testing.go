package bird_data_guessing

import "testing"

type stringTestCase struct {
	name     string
	expected string
}

type intTestCase struct {
	name     string
	expected int
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

type testStringBehavior func(englishOrLatinName string) (*Property, error)

type testIntBehavior func(englishOrLatinName string) (*Property, error)

func runIntTests(t *testing.T, b testIntBehavior, testCases ...intTestCase) {
	for _, c := range testCases {
		p, err := b(c.name)
		if err != nil {
			t.Errorf(`case '%s' yielded error: %+v`, c.name, err)
			continue
		}
		if c.expected != p.IntValue {
			t.Errorf(`case %s assertion error - expected '%d' - was '%d' - debug: %+v`, c.name, c.expected, p.IntValue, p)
		}
	}
}

func runStringTests(t *testing.T, b testStringBehavior, testCases ...stringTestCase) {
	for _, c := range testCases {
		p, err := b(c.name)
		if err != nil {
			t.Errorf(`case '%s' yielded error: %+v`, c.name, err)
			continue
		}
		if c.expected != p.StringValue {
			t.Errorf(`case %s assertion error - expected '%s' - was '%s' - debug: %+v`, c.name, c.expected, p.StringValue, p)
		}
	}
}
