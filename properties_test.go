package bird_data_guessing

import (
	"testing"
)

func TestFunFact(t *testing.T) {
	testCases := []struct {
		caseName string
		text     string
		expected string
	}{
		{"Basic", "an amazing animal. The guaddog can move quickly for cheese.", "The guaddog can move quickly for cheese."},
	}
	for _, tc := range testCases {
		wr := wikipediaResponse{}
		wr.setText(tc.text)
		p := wr.FunFact("guaddog")
		if tc.expected != p.StringValue {
			t.Errorf(`Case %s: Expected "%v", was "%v"`, tc.caseName, tc.expected, p.StringValue)
		}
	}
}

func TestWingspan(t *testing.T) {
	testCases := []struct {
		caseName string
		text     string
		expected int
	}{
		{"Basic", "average wingspan of 20cm ok", 20},
		{"With Decimal", "adults have a wingspan of 32.2cm", 32},
		{"With Inches", "female wingspan of 20 inches (30 cm)", 30},
		{"Centimetres", "adult wingspan of 20 inches (33 centimeters)", 33},
		{"Range", "females typically have a wingspan in the range of 22 - 32in (35 - 46cm)", 46},
		{"not found", "lol no egg wingspan here looser", 0},
	}
	for _, tc := range testCases {
		wr := wikipediaResponse{}
		wr.setText(tc.text)
		p := wr.Wingspan()
		if tc.expected != p.IntValue {
			t.Errorf(`Case %s: Expected "%v", was "%v"`, tc.caseName, tc.expected, p.IntValue)
		}
	}
}

func TestClutchSize(t *testing.T) {
	testCases := []struct {
		caseName string
		text     string
		expected int
	}{
		{"range english", "asdf lays one to five eggs asdf", 3},
		{"average number", "has an average clutch size of 7. Typically", 7},
		{"numerical range", "typically produces 12-14 eggs per year, in", 13},
		{"clutch number", "female goose may clutch up to six", 6},
		{"list format", "Average Laid Eggs: 12", 12},
		{"decimal format", "averages 8.02 eggs laid", 8},
		{"clutches of", "hens yield clutches of 4 without variation", 4},
		{"mixed terminology", "lay 3 to five eggs", 4},
		{"not found", "lol no egg count here looser", 0},
	}
	for _, tc := range testCases {
		wr := wikipediaResponse{}
		wr.setText(tc.text)
		p := wr.ClutchSize()
		if tc.expected != p.IntValue {
			t.Errorf(`Case %s: Expected "%v", was "%v"`, tc.caseName, tc.expected, p.IntValue)
		}
	}
}
