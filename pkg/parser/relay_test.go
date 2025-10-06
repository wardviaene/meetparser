package parser

import (
	"fmt"
	"testing"
)

func TestProcessRelayLine(t *testing.T) {
	lines := []string{
		// out5.txt
		"1 Swimteam1-ST     A 9:02.07 8:43.46 TAGS 40",
		"2 Swimteam2-GU     A 8:46.56 8:45.18 TAGS 34",
		"3 Swimteam3-NT     A 8:53.15 8:49.90 TAGS 32",
		"20 Swimteam4-GU     B 9:43.17 9:38.86  ",
	}
	expected := []RelayTime{
		{
			TeamName:            "Swimteam1",
			TeamLSC:             "ST",
			RelayEntry:          "A",
			Time:                "8:43.46",
			SeedTime:            "9:02.07",
			Place:               "1",
			QualifyingStandards: "TAGS",
			Points:              "40",
		},
		{
			TeamName:            "Swimteam2",
			TeamLSC:             "GU",
			RelayEntry:          "A",
			Time:                "8:45.18",
			SeedTime:            "8:46.56",
			Place:               "2",
			QualifyingStandards: "TAGS",
			Points:              "34",
		},
		{ // 3 Swimteam3-NT     A 8:53.15 8:49.90 TAGS 32
			TeamName:            "Swimteam3",
			TeamLSC:             "NT",
			RelayEntry:          "A",
			Time:                "8:49.90",
			SeedTime:            "8:53.15",
			Place:               "3",
			QualifyingStandards: "TAGS",
			Points:              "32",
		},
		{ // 20 Swim Streamline at Northampton-GU     B 9:43.17 9:38.86
			TeamName:            "Swimteam4",
			TeamLSC:             "GU",
			RelayEntry:          "B",
			Time:                "9:38.86",
			SeedTime:            "9:43.17",
			Place:               "20",
			QualifyingStandards: "",
			Points:              "",
		},
	}

	for k, line := range lines {
		relayTime, err := processRelayLineType1(line)
		if err != nil {
			t.Fatalf("error: %s. SwimmerTime: %+v", err, relayTime)
		}
		fmt.Printf("parsed: %+v\n", relayTime)
		if relayTime.TeamName != expected[k].TeamName {
			t.Fatalf("Teamname: got: '%s', expected: '%s'", relayTime.TeamName, expected[k].TeamName)
		}
		if relayTime.Time != expected[k].Time {
			t.Fatalf("Time: got: '%s', expected: '%s'", relayTime.Time, expected[k].Time)
		}
		if relayTime.SeedTime != expected[k].SeedTime {
			t.Fatalf("SeedTime: got: '%s', expected: '%s'", relayTime.SeedTime, expected[k].SeedTime)
		}
		if relayTime.TeamLSC != expected[k].TeamLSC {
			t.Fatalf("TeamLSC: got: '%s', expected: '%s'", relayTime.TeamLSC, expected[k].TeamLSC)
		}
		if relayTime.RelayEntry != expected[k].RelayEntry {
			t.Fatalf("RelayEntry: got: '%s', expected: '%s'", relayTime.RelayEntry, expected[k].RelayEntry)
		}
		if relayTime.QualifyingStandards != expected[k].QualifyingStandards {
			t.Fatalf("QualifyingStandards: got: '%s', expected: '%s'", relayTime.QualifyingStandards, expected[k].QualifyingStandards)
		}
		if relayTime.Points != expected[k].Points {
			t.Fatalf("Points: got: '%s', expected: '%s'", relayTime.Points, expected[k].Points)
		}
	}

}

func TestProcessRelayLineType2(t *testing.T) {
	lines := []string{
		// SWT.txt
		"1 SwimTeam A SWT 2:05.49 1:26.68",
		"6 SwimTeam A SWT 1:03.12 1:01.48 SWT",
	}
	expected := []RelayTime{
		{
			Place:         "1",
			TeamName:      "SwimTeam",
			TeamNameShort: "SWT",
			RelayEntry:    "A",
			SeedTime:      "2:05.49",
			Time:          "1:26.68",
		},
		{
			Place:         "6",
			TeamName:      "SwimTeam",
			TeamNameShort: "SWT",
			RelayEntry:    "A",
			SeedTime:      "1:03.12",
			Time:          "1:01.48",
			Achievements:  "SWT",
		},
	}

	for k, line := range lines {
		relayTime, err := processRelayLineType2(line)
		if err != nil {
			t.Fatalf("error: %s. SwimmerTime: %+v", err, relayTime)
		}
		fmt.Printf("parsed: %+v\n", relayTime)
		if relayTime.Place != expected[k].Place {
			t.Fatalf("Place: got: '%s', expected: '%s'", relayTime.Place, expected[k].Place)
		}
		if relayTime.TeamName != expected[k].TeamName {
			t.Fatalf("Teamname: got: '%s', expected: '%s'", relayTime.TeamName, expected[k].TeamName)
		}
		if relayTime.Time != expected[k].Time {
			t.Fatalf("Time: got: '%s', expected: '%s'", relayTime.Time, expected[k].Time)
		}
		if relayTime.SeedTime != expected[k].SeedTime {
			t.Fatalf("SeedTime: got: '%s', expected: '%s'", relayTime.SeedTime, expected[k].SeedTime)
		}
		if relayTime.TeamLSC != expected[k].TeamLSC {
			t.Fatalf("TeamLSC: got: '%s', expected: '%s'", relayTime.TeamLSC, expected[k].TeamLSC)
		}
		if relayTime.RelayEntry != expected[k].RelayEntry {
			t.Fatalf("RelayEntry: got: '%s', expected: '%s'", relayTime.RelayEntry, expected[k].RelayEntry)
		}
		if relayTime.QualifyingStandards != expected[k].QualifyingStandards {
			t.Fatalf("QualifyingStandards: got: '%s', expected: '%s'", relayTime.QualifyingStandards, expected[k].QualifyingStandards)
		}
		if relayTime.Points != expected[k].Points {
			t.Fatalf("Points: got: '%s', expected: '%s'", relayTime.Points, expected[k].Points)
		}
	}

}

func TestProcessSwimmersLineType1(t *testing.T) {
	lines := []string{
		"1) Name, Firstname1 14 2) Name, Firstname2 13 3) Name, Firstname3 14 4) Name, Firstname4 14",
	}
	expected := [][]RelaySwimmer{
		{
			{
				Place: "1",
				Name:  "Name, Firstname1",
				Age:   "14",
			},
			{
				Place: "2",
				Name:  "Name, Firstname2",
				Age:   "13",
			},
			{
				Place: "3",
				Name:  "Name, Firstname3",
				Age:   "14",
			},
			{
				Place: "4",
				Name:  "Name, Firstname4",
				Age:   "14",
			},
		},
	}

	for k, line := range lines {
		relaySwimmers, err := processRelaySwimmersLineType1(line)
		if err != nil {
			t.Fatalf("error: %s. SwimmerTime: %+v", err, relaySwimmers)
		}
		if len(expected[k]) != len(relaySwimmers) {
			for _, swimmer := range relaySwimmers {
				fmt.Printf("parsed: %+v\n", swimmer)
			}
			t.Fatalf("expected array is not the length same as result")
		}
		for kk, swimmer := range relaySwimmers {
			fmt.Printf("Swimmer: %+v\n", swimmer)
			if swimmer.Place != expected[k][kk].Place {
				t.Fatalf("Place: got: '%s', expected: '%s'", swimmer.Place, expected[k][kk].Place)
			}
			if swimmer.Name != expected[k][kk].Name {
				t.Fatalf("Name: got: '%s', expected: '%s'", swimmer.Name, expected[k][kk].Name)
			}
			if swimmer.Age != expected[k][kk].Age {
				t.Fatalf("Name: got: '%s', expected: '%s'", swimmer.Age, expected[k][kk].Age)
			}
		}

	}

}

func TestProcessSwimmersLineType2(t *testing.T) {
	lines := [][]string{
		{
			"1) Name, Firstname1 (6) 2) Name, Firstname2 (6)",
			"3) Name, Firstname3 (6) 4) Name, Firstname4 (6)",
		},
		{
			"1) - 2) Name, Firstname1 (17)",
			"3) Name, Firstname2 (17) 4) —",
		},
		{
			"1) Name, Firstname1 (17) 2) -",
			"3) — 4) —",
		},
	}
	expected := [][]RelaySwimmer{
		{
			{
				Place: "1",
				Name:  "Name, Firstname1",
				Age:   "6",
			},
			{
				Place: "2",
				Name:  "Name, Firstname2",
				Age:   "6",
			},
			{
				Place: "3",
				Name:  "Name, Firstname3",
				Age:   "6",
			},
			{
				Place: "4",
				Name:  "Name, Firstname4",
				Age:   "6",
			},
		},
		{
			{
				Place: "2",
				Name:  "Name, Firstname1",
				Age:   "17",
			},
			{
				Place: "3",
				Name:  "Name, Firstname2",
				Age:   "17",
			},
		},
		{
			{
				Place: "1",
				Name:  "Name, Firstname1",
				Age:   "17",
			},
		},
	}

	for k, line := range lines {
		relaySwimmers := []*RelaySwimmer{}
		for _, line2 := range line {
			relaySwimmersToAdd, err := processRelaySwimmersLineType2(line2)
			if err != nil {
				t.Fatalf("error: %s. SwimmerTime: %+v", err, relaySwimmers)
			}
			relaySwimmers = append(relaySwimmers, relaySwimmersToAdd...)
		}
		if len(expected[k]) != len(relaySwimmers) {
			for _, swimmer := range relaySwimmers {
				fmt.Printf("parsed: %+v\n", swimmer)
			}
			t.Fatalf("expected array is not the length same as result")
		}
		for kk, swimmer := range relaySwimmers {
			fmt.Printf("Swimmer: %+v\n", swimmer)
			if swimmer.Place != expected[k][kk].Place {
				t.Fatalf("Place: got: '%s', expected: '%s'", swimmer.Place, expected[k][kk].Place)
			}
			if swimmer.Name != expected[k][kk].Name {
				t.Fatalf("Name: got: '%s', expected: '%s'", swimmer.Name, expected[k][kk].Name)
			}
			if swimmer.Age != expected[k][kk].Age {
				t.Fatalf("Name: got: '%s', expected: '%s'", swimmer.Age, expected[k][kk].Age)
			}
		}

	}

}
