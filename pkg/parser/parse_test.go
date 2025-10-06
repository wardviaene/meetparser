package parser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestParserSummerSwimTeamResults(t *testing.T) {
	t.Skip() // local test
	files := []string{
		"../testdata/1.txt",
		"../testdata/2.txt",
		"../testdata/3.txt",
	}
	expectedEvents := []int{
		91,
		91,
		91,
	}
	expectedTimes := []int{
		1456,
		1076,
		697,
	}
	expectedRelays := []int{
		151,
		103,
		79,
	}
	for k, file := range files {
		out, err := ParsePDFText(file)
		if err != nil {
			t.Fatalf("error: %s\n", err)
		}
		if len(out.Events) != expectedEvents[k] {
			t.Fatalf("got %d events, expected %d", len(out.Events), expectedEvents[k])
		}
		if len(out.Times) != expectedTimes[k] {
			t.Fatalf("got %d swimmer times, expected %d", len(out.Times), expectedTimes[k])
		}
		if len(out.RelayTimes) != expectedRelays[k] {
			t.Fatalf("got %d relay times, expected %d", len(out.RelayTimes), expectedRelays[k])
		}
	}
}

func TestParserUSASwimTeamHeatSheets(t *testing.T) {
	t.Skip() // local test
	files := []string{
		"../testdata/3.txt",
		"../testdata/4.txt",
		"../testdata/5.txt",
	}
	expectedEvents := []int{
		144,
		183,
		12,
	}
	expectedTimes := []int{
		1519,
		3866,
		142,
	}
	expectedSplitTimes := []int{
		0,
		0,
		486,
	}
	expectedRelays := []int{
		58,
		495,
		0,
	}
	for k, file := range files {
		out, err := ParsePDFText(file)
		if err != nil {
			t.Fatalf("error: %s\n", err)
		}
		if len(out.ParseErrors) > 0 {
			fmt.Printf("Parse errors for %s:\n", files[k])
			for _, v := range out.ParseErrors {
				t.Fatalf("%+v\n", v)
			}

		}
		if len(out.Events) != expectedEvents[k] {
			for _, event := range out.Events {
				fmt.Printf("event: %+v\n", event)
			}
			t.Fatalf("got %d events, expected %d", len(out.Events), expectedEvents[k])
		}
		if len(out.Times) != expectedTimes[k] {
			t.Fatalf("got %d swimmer times, expected %d", len(out.Times), expectedTimes[k])
		}
		totalSplits := 0
		for _, times := range out.Times {
			totalSplits += len(times.SplitTimes)
		}
		if totalSplits != expectedSplitTimes[k] {
			t.Fatalf("got %d split times, expected %d", totalSplits, expectedSplitTimes[k])
		}
		if len(out.RelayTimes) != expectedRelays[k] {
			t.Fatalf("got %d relay times, expected %d", len(out.RelayTimes), expectedRelays[k])
		}
	}
}

func TestSplitLastXChars(t *testing.T) {
	out1, out2 := splitLastXChars("123", 2)
	if out1 != "1" {
		t.Fatalf("Expected 1, got: %s", out1)
	}
	if out2 != "23" {
		t.Fatalf("Expected 23, got: %s", out2)
	}
}

func TestIsEvent(t *testing.T) {
	tests := []struct {
		line     string
		fileType string
		expected bool
	}{
		// FILETYPE_TYPE2 tests
		{"", FILETYPE_TYPE2, false},
		{"#1 Swimming", FILETYPE_TYPE2, true},  // starts with # + numeric
		{"#A Swimming", FILETYPE_TYPE2, false}, // starts with # but not numeric
		{"Swimming #1", FILETYPE_TYPE2, false}, // does not start with #
		{"#9", FILETYPE_TYPE2, true},           // minimal numeric
		{"#", FILETYPE_TYPE2, false},           // too short, len=1

		// default file type tests
		{"Event 1", FILETYPE_TYPE1, true},     // starts with Event
		{"#1 Swimming", FILETYPE_TYPE1, true}, // default type, Event
		{"Other text", FILETYPE_TYPE1, false}, // no Event
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got := isEvent(tt.line, tt.fileType)
			if got != tt.expected {
				t.Errorf("isEvent(%q, %q) = %v; want %v", tt.line, tt.fileType, got, tt.expected)
			}
		})
	}
}

func TestParsePDFTextBadData(t *testing.T) {
	a := bytes.NewBufferString("faulty\ndata\nEvent 3\nsome data here\n")
	res, err := parsePDFText(a)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
	if len(res.ParseErrors) == 0 {
		t.Fatalf("expected parse error. Got 0")
	}
}

func TestEventAddQualifyingTimes(t *testing.T) {
	expectedTime := "28.51"
	event := &Event{
		QualifyingTimes: make(map[string]string),
	}
	line := "INV NWSC Invitational Meet Qualifying Times '25 (Girls 6&U) 28.51"
	err := eventAddQualifyingTimes(event, line)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if val, ok := event.QualifyingTimes["INV NWSC Invitational Meet"]; ok {
		if val != expectedTime {
			t.Fatalf("expected different time. Got %s, expected %s", val, expectedTime)
		}
	} else {
		t.Fatalf("error: qualifying time 'key' not found. Qualifying times: %+v", event.QualifyingTimes)
	}
}
