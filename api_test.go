package bird_data_guessing

import (
	"testing"
)

func SetGetText(t *testing.T) {
	e := "asdklfhqp9weuicna"

	wr := &wikipediaResponse{}
	wr.setText(e)
	a := wr.GetText()

	if a != e {
		t.Errorf("Expected %s, was %s", e, a)
	}
}

func TestCountMatches(t *testing.T) {
	testCases := []struct {
		caseName  string
		text      string
		searchFor []string
		expected  int
	}{
		{"Basic", "I have a feeling ooooohh yeah tonights gonna be a", []string{"night"}, 1},
		{"Multi-Word", "good night, that tonights gonna be a good night,", []string{"night", "good"}, 5},
		{"Case-insensitive", "TONIGHTS THE NIGHT. LETS TURN IT UP. ", []string{"night"}, 2},
	}
	for _, tc := range testCases {
		wr := wikipediaResponse{}
		wr.setText(tc.text)
		actual := wr.CountMatches(tc.searchFor)
		if tc.expected != actual.Strength {
			t.Errorf(`Case: %s Wanted "%v", was "%v"`, tc.caseName, tc.expected, actual.Strength)
		}
	}
}

func TestExtractMatch(t *testing.T) {
	testCases := []struct {
		caseName     string
		text         string
		regex        string
		captureGroup int
		expected     string
	}{
		{"Basic", "I got that boom boom pow", " (bo+m) ", 1, "boom"},
		{"Not Found", "Them chickens be jacking my style", "next shit now", 2, ""},
		{"Multiple Capture Groups", "they try to copy my swagger", "\\s+([^\\s]+)\\s+(\\w+)\\s", 2, "to"},
		{"Nested Capture Groups", "i'm on that next shit now", "(t(h(a(t))))", 3, "at"},
	}
	for _, tc := range testCases {
		wr := wikipediaResponse{}
		wr.setText(tc.text)
		actual := wr.ExtractMatch(tc.regex, tc.captureGroup)
		if tc.expected != actual.StringValue {
			t.Errorf(`Case: %s Wanted "%v", was "%v"`, tc.caseName, tc.expected, actual.StringValue)
		}
	}
}

func TestExtractAnyMatch(t *testing.T) {
	testCases := []struct {
		caseName string
		text     string
		patterns map[string]int
		expected string
	}{
		{"Basic", "Father father father help us", map[string]int{
			"mother":       1,
			"brother":      2,
			"sister":       0,
			"uncle":        1,
			"cousin":       0,
			"aunt":         1,
			"((father )+)": 1,
			"nibling":      0}, "Father father father "},
		{"Uses Associated Capture", "Need some guidance from above", map[string]int{
			"(n(e+)).*guidance (\\w+)": 3}, "from"}}
	for _, tc := range testCases {
		wr := wikipediaResponse{}
		wr.setText(tc.text)
		actual := wr.ExtractAnyMatch(tc.patterns)
		if tc.expected != actual.StringValue {
			t.Errorf(`Case: %s Wanted "%v", was "%v"`, tc.caseName, tc.expected, actual.StringValue)
		}
	}
}
