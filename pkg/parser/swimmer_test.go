package parser

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessLineType1(t *testing.T) {
	lines := []string{
		"1 Lastname, Firstname A  14 Mansfield Aquatic Club-NT 9:27.94 9:18.83 TAGS 20",
		"2 Lastname, Firstname  13 Metroplex Aquatics-NT 9:26.82 9:20.36 TAGS 17",
		"3 Lastname, Firstname S  14 Alamo Area Aquatic Association-ST 9:30.09 9:24.36 TAGS 16",
		"2 Lastname, Firstname J  10 Rockwall Aquatic Center of Exc-NT 2:45.55 2:40.93 TAGS 17",
		"3 Lastname, Firstname  10 Nitro Swimming-ST 2:48.71 2:44.99 TAGS 16",
		"1 Lastname, Firstname  14 Lynchburg YMCA 2:14.96 2:16.72 AG 9",
		"2 Lastname, Firstname  14 Lynchburg YMCA 2:20.31 2:22.04 7",
		"3 Lastname, Firstname  13 Lynchburg YMCA 2:26.59 2:24.11 6",
		"--- Lastname, Firstname   9 Lynchburg YMCA NT DQ ",
		"7 Lastname, Firstname   5 Lynchburg YMCA NT 33.49 2",
		"--- Lastname, Firstname   8 Lynchburg YMCA 30.88 DQ  ",
		"6 Lastname, Firstname O  14 Mansfield Aquatic Club-NT 10:43.41 Y 9:29.11 TAGS 13",
		"17 Lastname, Firstname F  13 The Woodlands Swim Team-GU 10:49.69 Y 9:43.85 TAGS  ",
		"1 Lastname, Firstname J  12 SJAC-MA 1:09.33",
		"--- Lastname, Firstname  14 Nation's Capital Swim Club DFS",
		"36 Lastname, Firstname  17 Nation's Capital Swim Club J22.27",
		"--- Lastname, Firstname  17 Occoquan Swimming DQ",
		"1 Lastname, Firstname  14 Nation's Capital Swim Club 21.27 21.26 # q",
	}
	expected := []SwimmerTime{
		{
			Name:                "Lastname, Firstname A",
			TeamName:            "Mansfield Aquatic Club",
			TeamLSC:             "NT",
			Time:                "9:18.83",
			SeedTime:            "9:27.94",
			Age:                 "14",
			Place:               "1",
			Points:              "20",
			QualifyingStandards: "TAGS",
		},
		{ // "2 Lastname, Firstname  13 Metroplex Aquatics-NT 9:26.82 9:20.36 TAGS 17",
			Name:                "Lastname, Firstname",
			TeamName:            "Metroplex Aquatics",
			TeamLSC:             "NT",
			Time:                "9:20.36",
			SeedTime:            "9:26.82",
			Age:                 "13",
			Place:               "2",
			Points:              "17",
			QualifyingStandards: "TAGS",
		},
		{ // "3 Lastname, Firstname S  14 Alamo Area Aquatic Association-ST 9:30.09 9:24.36 TAGS 16",
			Name:                "Lastname, Firstname S",
			TeamName:            "Alamo Area Aquatic Association",
			TeamLSC:             "ST",
			Time:                "9:24.36",
			SeedTime:            "9:30.09",
			Age:                 "14",
			Place:               "3",
			Points:              "16",
			QualifyingStandards: "TAGS",
		},
		{ // "2 Lastname, Firstname J  10 Rockwall Aquatic Center of Exc-NT 2:45.55 2:40.93 TAGS 17",
			Name:                "Lastname, Firstname J",
			TeamName:            "Rockwall Aquatic Center of Exc",
			TeamLSC:             "NT",
			Time:                "2:40.93",
			SeedTime:            "2:45.55",
			Place:               "2",
			Age:                 "10",
			Points:              "17",
			QualifyingStandards: "TAGS",
		},
		{
			Name:                "Lastname, Firstname",
			TeamName:            "Nitro Swimming",
			TeamLSC:             "ST",
			Time:                "2:44.99",
			SeedTime:            "2:48.71",
			Place:               "3",
			Age:                 "10",
			Points:              "16",
			QualifyingStandards: "TAGS",
		},
		{
			Name:                "Lastname, Firstname",
			TeamName:            "Lynchburg YMCA",
			Time:                "2:16.72",
			SeedTime:            "2:14.96",
			Age:                 "14",
			Place:               "1",
			Points:              "9",
			QualifyingStandards: "AG",
		},
		{
			Name:     "Lastname, Firstname",
			TeamName: "Lynchburg YMCA",
			Time:     "2:22.04",
			SeedTime: "2:20.31",
			Age:      "14",
			Place:    "2",
			Points:   "7",
		},
		{
			Name:     "Lastname, Firstname",
			TeamName: "Lynchburg YMCA",
			Time:     "2:24.11",
			SeedTime: "2:26.59",
			Age:      "13",
			Place:    "3",
			Points:   "6",
		},
		{
			Name:     "Lastname, Firstname",
			TeamName: "Lynchburg YMCA",
			Time:     "DQ",
			Place:    "---",
			SeedTime: "NT",
			Age:      "9",
		},
		{
			Name:     "Lastname, Firstname",
			TeamName: "Lynchburg YMCA",
			Time:     "33.49",
			Place:    "7",
			SeedTime: "NT",
			Age:      "5",
			Points:   "2",
		},
		{
			Name:     "Lastname, Firstname",
			TeamName: "Lynchburg YMCA",
			Time:     "DQ",
			Place:    "---",
			SeedTime: "30.88",
			Age:      "8",
		},
		{ // "6 Lastname, Firstname O  14 Mansfield Aquatic Club-NT 10:43.41 Y 9:29.11 TAGS 13",
			Name:                "Lastname, Firstname O",
			TeamName:            "Mansfield Aquatic Club",
			TeamLSC:             "NT",
			Time:                "9:29.11",
			Place:               "6",
			SeedTime:            "10:43.41",
			SeedTimeTag:         "Y",
			Points:              "13",
			QualifyingStandards: "TAGS",
			Age:                 "14",
		},
		{ // "17 Lastname, Firstname F  13 The Woodlands Swim Team-GU 10:49.69 Y 9:43.85 TAGS  "
			Name:                "Lastname, Firstname F",
			TeamName:            "The Woodlands Swim Team",
			TeamLSC:             "GU",
			Time:                "9:43.85",
			Place:               "17",
			SeedTime:            "10:49.69",
			SeedTimeTag:         "Y",
			Points:              "",
			QualifyingStandards: "TAGS",
			Age:                 "13",
		},
		{ // "1 Lastname, Firstname J  12 SJAC-MA 1:09.33",
			Name:     "Lastname, Firstname J",
			TeamName: "SJAC",
			TeamLSC:  "MA",
			Time:     "1:09.33",
			Place:    "1",
			Age:      "12",
		},
		{ // --- Lastname, Firstname 14 Nation's Capital Swim Club DFS
			Name:     "Lastname, Firstname",
			TeamName: "Nation's Capital Swim Club",
			TeamLSC:  "",
			Time:     "DFS",
			Place:    "---",
			Age:      "14",
		},
		{ // --- 36 Lastname, Firstname  17 Nation's Capital Swim Club J22.27
			Name:     "Lastname, Firstname",
			TeamName: "Nation's Capital Swim Club",
			TeamLSC:  "",
			Time:     "J22.27",
			Place:    "36",
			Age:      "17",
		},
		{ // "--- Lastname, Firstname  17 Occoquan Swimming DQ"
			Name:     "Lastname, Firstname",
			TeamName: "Occoquan Swimming",
			TeamLSC:  "",
			Time:     "DQ",
			Place:    "---",
			Age:      "17",
		},
		{ // "1 Lastname, Firstname  14 Nation's Capital Swim Club 21.27 21.26 # q",
			Name:      "Lastname, Firstname",
			TeamName:  "Nation's Capital Swim Club",
			TeamLSC:   "",
			SeedTime:  "21.27",
			Time:      "21.26",
			Place:     "1",
			Qualified: true,
			NewRecord: true,
			Age:       "14",
		},
	}

	for k, line := range lines {
		parsed, err := processLineType1(line)
		if err != nil {
			t.Fatalf("error: %s. SwimmerTime: %+v", err, parsed)
		}
		fmt.Printf("parsed: %+v\n", parsed)
		if diff := cmp.Diff(expected[k], *parsed); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestProcessLineType2(t *testing.T) {
	lines := []string{
		// SWT.txt (invitationals)
		"1 Lastname, Firstname 6 PFP 18.14 18.39",
		"2 Lastname, Firstname 6 AM 19.95 20.97",
		"3 Lastname, Firstname 6 BC25 22.47 22.13",
		"9 Lastnamé, Firstnamë 7 BC25 24.26 23.63",
		"11 Lastname, Firstname 12 SWT 34.57 34.08 SWT",
		"1 Lastname, Firstname 6 PFP 18.14 20.06 9 INV",
		"1 Lastname, Firstname 6 PFP 18.14 20.06 9",
		"1 Lastname, Firstname 6 PFP 18.14 20.06 INV",
	}
	expected := []SwimmerTime{
		{
			Name:     "Lastname, Firstname",
			TeamName: "PFP",
			SeedTime: "18.14",
			Time:     "18.39",
			Age:      "6",
			Place:    "1",
		},
		{
			Name:     "Lastname, Firstname",
			TeamName: "AM",
			SeedTime: "19.95",
			Time:     "20.97",
			Age:      "6",
			Place:    "2",
		},
		{
			Name:     "Lastname, Firstname",
			TeamName: "BC25",
			SeedTime: "22.47",
			Time:     "22.13",
			Age:      "6",
			Place:    "3",
		},
		{ // 9 Lastnamé, Firstnamë 7 BC25 24.26 23.63"
			Name:     "Lastnamé, Firstnamë",
			TeamName: "BC25",
			SeedTime: "24.26",
			Time:     "23.63",
			Age:      "7",
			Place:    "9",
		},
		{ // "11 Lastname, Firstname 12 SWT 34.57 34.08 SWT",
			Name:         "Lastname, Firstname",
			TeamName:     "SWT",
			SeedTime:     "34.57",
			Time:         "34.08",
			Age:          "12",
			Place:        "11",
			Achievements: "SWT",
		},
		{ // 1 Lastname, Firstname 6 PFP 18.14 20.06 9 INV
			Name:         "Lastname, Firstname",
			TeamName:     "PFP",
			SeedTime:     "18.14",
			Time:         "20.06",
			Age:          "6",
			Place:        "1",
			Points:       "9",
			Achievements: "INV",
		},
		{ // 1 Lastname, Firstname 6 PFP 18.14 20.06 9
			Name:         "Lastname, Firstname",
			TeamName:     "PFP",
			SeedTime:     "18.14",
			Time:         "20.06",
			Age:          "6",
			Place:        "1",
			Points:       "9",
			Achievements: "",
		},
		{ // 1 Lastname, Firstname 6 PFP 18.14 20.06 INV
			Name:         "Lastname, Firstname",
			TeamName:     "PFP",
			SeedTime:     "18.14",
			Time:         "20.06",
			Age:          "6",
			Place:        "1",
			Points:       "",
			Achievements: "INV",
		},
	}

	for k, line := range lines {
		parsed, err := processLineType2(line)
		if err != nil {
			t.Fatalf("error: %s. SwimmerTime: %+v", err, parsed)
		}
		fmt.Printf("parsed: %+v\n", parsed)
		if parsed.TeamName != expected[k].TeamName {
			t.Fatalf("Teamname: got: '%s', expected: '%s'", parsed.TeamName, expected[k].TeamName)
		}
		if parsed.Time != expected[k].Time {
			t.Fatalf("Time: got: '%s', expected: '%s'", parsed.Time, expected[k].Time)
		}
		if parsed.SeedTime != expected[k].SeedTime {
			t.Fatalf("SeedTime: got: '%s', expected: '%s'", parsed.SeedTime, expected[k].SeedTime)
		}
		if parsed.SeedTimeTag != expected[k].SeedTimeTag {
			t.Fatalf("SeedTime Tag: got: '%s', expected: '%s'", parsed.SeedTimeTag, expected[k].SeedTimeTag)
		}
		if parsed.Name != expected[k].Name {
			t.Fatalf("Name: got: '%s', expected: '%s'", parsed.Name, expected[k].Name)
		}
		if parsed.Age != expected[k].Age {
			t.Fatalf("Age: got: '%s', expected: '%s'", parsed.Age, expected[k].Age)
		}
		if parsed.Place != expected[k].Place {
			t.Fatalf("Place: got: '%s', expected: '%s'", parsed.Place, expected[k].Place)
		}
		if parsed.Achievements != expected[k].Achievements {
			t.Fatalf("Achievements: got: '%s', expected: '%s'", parsed.Achievements, expected[k].Achievements)
		}
		if parsed.Points != expected[k].Points {
			t.Fatalf("Points: got: '%s', expected: '%s'", parsed.Points, expected[k].Points)
		}
	}

}

func TestStringAgeIndex(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"abc 3 x", 4},        // 1-digit
		{"abc 12 y", 4},       // 2-digit
		{"abc 123 z", 4},      // 3-digit
		{"1234 z", -1},        // no leading space
		{"a123 456 b", 5},     // second number
		{"no match here", -1}, // no match
		{"café 12 x", 6},      // Unicode character before number
		{" café 1 x", 7},      // leading space
		{"a é 7 b", 5},        // multi-byte character between spaces
		{"emoji 😊 99 x", 11},  // emoji before number
		{"a 13 12 c", 2},      // first valid number only
		{"a 1234 b", -1},      // 4-digit number, should fail
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := stringAgeIndex(tt.input)
			if got != tt.expected {
				t.Errorf("stringAgeIndex(%q) = %d; want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTimeRegex(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"12.34", true},
		{"1:12.34", true},
		{"12.34 1:12.34", true},
		{"12.34 abc", false},
		{"abc 12.34", false},
		{"12.34 1:12.34 3:45.67", true},
		{"", false}, // optional: empty string should not match
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := splitTimesRegex.MatchString(tc.input)
			if actual != tc.expected {
				t.Errorf("input %q: expected %v, got %v", tc.input, tc.expected, actual)
			}
		})
	}
}
