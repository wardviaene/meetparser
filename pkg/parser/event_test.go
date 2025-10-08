package parser

import (
	"fmt"
	"testing"
)

func TestProcessEventType1(t *testing.T) {
	lines := []string{
		"Event 1  Girls 13-14 200 Yard IM",
		"Event 1  Girls 15 & Over 200 Yard IM",
		"Event 2  Boys 13-14 200 Yard IM",
		// no age group
		"Event 1  Girls 1500 LC Meter Freestyle",
		//
		"Event 7  Women 15 & Over 200 LC Meter IM",
		//
		"#1  Girls 11-12 100 Yard Fly",
		"#3  Girls 13&O 200 Yard Fly",
		"Event	1		Girls	10	&	Under	200	Yard	IM",
	}
	expected := []Event{
		{
			Round:    "1",
			Gender:   "girls",
			AgeGroup: "13-14",
			Distance: "200 Yard",
			Stroke:   "IM",
			Relay:    false,
		},
		{
			Round:    "1",
			Gender:   "girls",
			AgeGroup: "15 & over",
			Distance: "200 Yard",
			Stroke:   "IM",
			Relay:    false,
		},
		{
			Round:    "2",
			Gender:   "boys",
			AgeGroup: "13-14",
			Distance: "200 Yard",
			Stroke:   "IM",
			Relay:    false,
		},
		{
			Round:    "1",
			Gender:   "girls",
			AgeGroup: "",
			Distance: "1500 LC Meter",
			Stroke:   "Freestyle",
			Relay:    false,
		},
		{ // "Event 7  Women 15 & Over 200 LC Meter IM",
			Round:    "7",
			Gender:   "women",
			AgeGroup: "15 & over",
			Distance: "200 LC Meter",
			Stroke:   "IM",
			Relay:    false,
		},
		{ // "#1  Girls 11-12 100 Yard Fly",
			Round:    "1",
			Gender:   "girls",
			AgeGroup: "11-12",
			Distance: "100 Yard",
			Stroke:   "Fly",
			Relay:    false,
		}, { // "#3  Girls 13&O 200 Yard Fly",
			Round:    "3",
			Gender:   "girls",
			AgeGroup: "13 & over",
			Distance: "200 Yard",
			Stroke:   "Fly",
			Relay:    false,
		},
		{ //Event	1		Girls	10	&	Under	200	Yard	IM
			Round:    "1",
			Gender:   "girls",
			AgeGroup: "10 & under",
			Distance: "200 Yard",
			Stroke:   "IM",
			Relay:    false,
		},
	}

	for k, line := range lines {
		parsed, err := processEventType1(line)
		if err != nil {
			t.Fatalf("error: %s. SwimmerTime: %+v", err, parsed)
		}
		fmt.Printf("parsed: %+v\n", parsed)
		if parsed.Round != expected[k].Round {
			t.Fatalf("Name: got: '%s', expected: '%s'", parsed.Round, expected[k].Round)
		}
		if parsed.AgeGroup != expected[k].AgeGroup {
			t.Fatalf("AgeGroup: got: '%s', expected: '%s'", parsed.AgeGroup, expected[k].AgeGroup)
		}
		if parsed.Distance != expected[k].Distance {
			t.Fatalf("Distance: got: '%s', expected: '%s'", parsed.Distance, expected[k].Distance)
		}
		if parsed.Gender != expected[k].Gender {
			t.Fatalf("Gender: got: '%s', expected: '%s'", parsed.Gender, expected[k].Gender)
		}
		if parsed.Stroke != expected[k].Stroke {
			t.Fatalf("Achievements: got: '%s', expected: '%s'", parsed.Stroke, expected[k].Stroke)
		}
		if parsed.Relay != expected[k].Relay {
			t.Fatalf("Relay: got: '%v', expected: '%v'", parsed.Relay, expected[k].Relay)
		}
	}

}

func TestProcessEventType2(t *testing.T) {
	lines := []string{
		"#1 Mixed 6 & Under 100yd Freestyle Relay",
	}
	expected := []Event{
		{
			Round:    "1",
			Gender:   "mixed",
			AgeGroup: "6 & under",
			Stroke:   "Freestyle",
			Distance: "100yd",
			Relay:    true,
		},
	}

	for k, line := range lines {
		parsed, err := processEventType2(line)
		if err != nil {
			t.Fatalf("error: %s. SwimmerTime: %+v", err, parsed)
		}
		fmt.Printf("parsed: %+v\n", parsed)
		if parsed.Round != expected[k].Round {
			t.Fatalf("Name: got: '%s', expected: '%s'", parsed.Round, expected[k].Round)
		}
		if parsed.AgeGroup != expected[k].AgeGroup {
			t.Fatalf("AgeGroup: got: '%s', expected: '%s'", parsed.AgeGroup, expected[k].AgeGroup)
		}
		if parsed.Distance != expected[k].Distance {
			t.Fatalf("Distance: got: '%s', expected: '%s'", parsed.Distance, expected[k].Distance)
		}
		if parsed.Gender != expected[k].Gender {
			t.Fatalf("Gender: got: '%s', expected: '%s'", parsed.Gender, expected[k].Gender)
		}
		if parsed.Stroke != expected[k].Stroke {
			t.Fatalf("Achievements: got: '%s', expected: '%s'", parsed.Stroke, expected[k].Stroke)
		}
		if parsed.Relay != expected[k].Relay {
			t.Fatalf("Relay: got: '%v', expected: '%v'", parsed.Relay, expected[k].Relay)
		}
	}

}

func TestParseEventAgeGroup(t *testing.T) {
	tests := []struct {
		input     string
		expected  string
		wantEmpty bool
	}{
		{"6 & Under 25 Free", "6 & under", false},
		{"9 & Under 50 Fly", "9 & under", false},
		{"12 & Over 100 IM", "12 & over", false},
		{"13 & O 200 Free", "13 & o", false},
		{"14&O 50 Back", "14&o", false},
		{"15-16 100 Free", "15-16", false},
		{"9-10 25 Fly", "9-10", false},
		{"Open 100 Free", "", true},
		{"Girls 11-12 100 Free", "", true}, // prefix not at start
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseEventAgeGroup(tt.input)
			if (got == "") != tt.wantEmpty {
				t.Fatalf("normalizeAgeGroup(%q) error = %v, wantErr %v", tt.input, err, tt.wantEmpty)
			}
			if got != tt.expected {
				t.Errorf("normalizeAgeGroup(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
