package parser

import (
	"fmt"
	"strconv"
	"strings"
)

func processRelayLine(line string, fileType string) (*RelayTime, error) {
	switch fileType {
	case FILETYPE_TYPE2:
		return processRelayLineType2(line)
	default:
		return processRelayLineType1(line)
	}
}
func processRelayLineType2(line string) (*RelayTime, error) {
	var err error
	relayTime := &RelayTime{
		Swimmers: []*RelaySwimmer{},
	}
	// line: 1 SwimTeam A SWT 2:05.49 1:26.68
	// line: 6 SwimTeam A SWT 1:03.12 1:01.48 SWT
	index1 := strings.Index(line, " ")
	if index1 == -1 {
		return relayTime, fmt.Errorf("place not found")
	}
	relayTime.Place = line[0:index1]
	line = line[index1+1:]
	// line: SwimTeam A SWT 2:05.49 1:26.68
	index2 := relayLetterIndex(line)
	if index2 == -1 {
		return relayTime, fmt.Errorf("couldn't find relay letter indicator")
	}
	relayTime.RelayEntry = line[index2 : index2+1]
	relayTime.TeamName = line[0 : index2-1]
	line = line[index2+2:]
	// line: SWT 2:05.49 1:26.68
	index3 := strings.Index(line, " ")
	if index3 == -1 {
		return relayTime, fmt.Errorf("couldn't find relay team (short) indicator")
	}
	relayTime.TeamNameShort = line[0:index3]
	line = line[index3+1:]
	// line: 2:05.49 1:26.68
	index4 := strings.Index(line, " ")
	if index4 == -1 {
		return relayTime, fmt.Errorf("couldn't find relay seed time indicator")
	}
	relayTime.SeedTime = line[0:index4]
	line = line[index4+1:]
	// line: 1:26.68
	index5 := strings.Index(line, " ")
	if index5 == -1 {
		relayTime.Time = line
		line = ""
	} else {
		relayTime.Time = line[0:index5]
		line = line[index5+1:]
	}

	// line: SWT
	// line: 9 INV
	if line != "" {
		index7 := strings.Index(line, " ")
		if index7 == -1 { // line: SWT    or   line: 9
			if !isNumeric(line) {
				relayTime.Achievements = line
			} else {
				relayTime.Points = line
			}
			line = ""
		} else {
			if !isNumeric(line[0:index7]) { // line: INV Somethingelse
				relayTime.Achievements = line[0:index7]
				line = line[index7+1:]
			} else { // line: 9 INV
				relayTime.Points = line[0:index7]
				line = line[index7+1:]
			}
		}
	}

	err = checkResidual(relayTime.SeedTime + " " + relayTime.Time + " " + line)
	if err != nil {
		return relayTime, fmt.Errorf("residual information found: '%s'", err)
	}
	return relayTime, nil
}
func processRelayLineType1(line string) (*RelayTime, error) {
	var err error
	relayTime := &RelayTime{
		Swimmers: []*RelaySwimmer{},
	}
	// line: 1 Nitro Swimming-ST     A 9:02.07 8:43.46 TAGS 40
	index1 := strings.Index(line, " ")
	if index1 == -1 {
		return relayTime, fmt.Errorf("place not found")
	}
	relayTime.Place = line[0:index1]
	line = line[index1+1:]
	// line: Nitro Swimming-ST     A 9:02.07 8:43.46 TAGS 40
	index2 := strings.Index(line, "     ")
	index2Offset := len("     ")
	if index2 == -1 {
		index2 = strings.Index(line, "    ")
		if index2 == -1 {
			return relayTime, fmt.Errorf("error while parsing team name")
		}
		index2Offset = 2
	}
	split1 := strings.Split(line[0:index2], "-")
	if len(split1) == 2 {
		relayTime.TeamName = split1[0]
		relayTime.TeamLSC = split1[1]
	} else {
		relayTime.TeamName = line[0:index2]
	}
	line = line[index2+index2Offset:]

	// line: A 9:02.07 8:43.46 TAGS 40
	relayLetterIndex := strings.Index(line, " ")
	if relayLetterIndex == -1 {
		return relayTime, fmt.Errorf("relay letter index not found")
	}
	relayTime.RelayEntry = line[0:relayLetterIndex]
	line = line[relayLetterIndex+1:]

	indexRelayLetter := -1
	indexRelayLetter, relayTime.SeedTime, relayTime.Time, err = processTimes(line)
	if err != nil {
		fmt.Printf("Line: %s\n", line)
		return relayTime, fmt.Errorf("process time error: %s", err)
	}

	line = line[indexRelayLetter:]

	// Extract points
	// line: 9:02.07 8:43.46 TAGS 40
	if !strings.HasSuffix(line, "  ") {
		pointsIndex := strings.LastIndex(line, " ")
		relayTime.Points = line[pointsIndex+1:]
		line = line[0:pointsIndex]
	}

	// Extract qualifying standards
	// line: "9:02.07 8:43.46 TAGS"
	line = strings.TrimSpace(line) // remove unnecessary spacing
	qualifyingStandardsIndex := strings.LastIndex(line, " ")
	if !timesRegex.MatchString(line[qualifyingStandardsIndex+1:]) {
		if line[qualifyingStandardsIndex+1:] != "DQ" && line[qualifyingStandardsIndex+1:] != "NS" && line[qualifyingStandardsIndex+1:] != "DNF" && line[qualifyingStandardsIndex+1:] != "DFS" {
			relayTime.QualifyingStandards = line[qualifyingStandardsIndex+1:]
			line = line[0:qualifyingStandardsIndex]
		}
	}
	// line: 10:43.41 Y 9:29.11
	seedTagIndex := strings.Index(line, " ")
	if seedTagIndex == -1 {
		return relayTime, fmt.Errorf("seed tag index not found")
	}
	if strings.HasPrefix(line[seedTagIndex+1:], "Y ") || strings.HasPrefix(line[seedTagIndex+1:], "S ") || strings.HasPrefix(line[seedTagIndex+1:], "L ") {
		relayTime.SeedTimeTag = line[seedTagIndex+1 : seedTagIndex+2]
		line = line[0:seedTagIndex+1] + line[seedTagIndex+3:]
	}

	//fmt.Printf("line: %s\n", line)

	err = checkResidual(line)
	if err != nil {
		return relayTime, fmt.Errorf("residual information found: '%s'", err)
	}

	return relayTime, nil
}

func processRelaySwimmersLine(line string, fileType string) ([]*RelaySwimmer, error) {
	switch fileType {
	case FILETYPE_TYPE2:
		return processRelaySwimmersLineType2(line)
	default:
		return processRelaySwimmersLineType1(line)
	}
}

func processRelaySwimmersLineType2(line string) ([]*RelaySwimmer, error) {
	swimmers := []*RelaySwimmer{}
	end := false
	// line: 1) Lastname, Firstname (6) 2) Raley, Colton (6)
	for i := 1; !end; i++ {
		relaySwimmer := &RelaySwimmer{}
		index1 := strings.Index(line, ")")
		if index1 == -1 {
			return swimmers, fmt.Errorf("couldn't determine place")
		}
		relaySwimmer.Place = line[0:index1]
		line = line[index1+2:]

		// line: Lastname, Firstname (6) 2) Raley, Colton (6)
		// line: —
		// line: - 2) Raley, Colton (6)
		if strings.TrimSpace(line) == "-" || strings.TrimSpace(line) == "—" {
			break // no swimmer and we're at the end
		}
		if (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "—")) && strings.Contains(line, ")") {
			line = line[2:]
			continue // invalid swimmer, but another one afterwards
		}
		index2 := strings.Index(line, "(")
		if index2 == -1 {
			return swimmers, fmt.Errorf("couldn't find name")
		}
		relaySwimmer.Name = line[0 : index2-1]
		line = line[index2+1:]

		// line: 6) 2) Raley, Colton (6)
		index3 := strings.Index(line, ")")
		if index3 == -1 {
			return swimmers, fmt.Errorf("couldn't determine age")
		}
		relaySwimmer.Age = line[0:index3]
		line = line[index3+1:]

		if line == "" {
			end = true
		} else {
			line = line[1:] // move one more for the last space
		}

		swimmers = append(swimmers, relaySwimmer)
	}

	return swimmers, nil
}

func processRelaySwimmersLineType1(line string) ([]*RelaySwimmer, error) {
	swimmers := []*RelaySwimmer{}
	end := false
	// line: 1) Lastname, Firstname 14 2) Gunn, Pepper 13 3) Peeters, Hanne 14 4) Sitter, Gianna 14
	for i := 1; !end; i++ {
		relaySwimmer := &RelaySwimmer{}
		index1 := strings.Index(line, ")")
		if index1 == -1 {
			return swimmers, fmt.Errorf("couldn't determine place")
		}
		relaySwimmer.Place = line[0:index1]
		line = line[index1+2:]

		// line: Lastname, Firstname 14 2) Gunn, Pepper 13 3) Peeters, Hanne 14 4) Sitter, Gianna 14
		rightIndex := strings.Index(line, ")")
		offset2 := 0
		if rightIndex == -1 {
			end = true
			rightIndex = len(line) + 1
			offset2 = rightIndex - 1
		} else {
			offset2 = rightIndex - 2
		}
		// line: Lastname, Firstname 14

		findSpace := strings.LastIndex(line[:offset2], " ")
		if findSpace == -1 {
			return swimmers, fmt.Errorf("expected spacing before 'number)'")
		}
		relaySwimmer.Age = line[findSpace+1 : offset2]
		relaySwimmer.Name = line[:findSpace]
		// next part of the line
		if !end {
			line = line[rightIndex-len(strconv.Itoa(i+1)):]
		}
		swimmers = append(swimmers, relaySwimmer)
	}
	return swimmers, nil
}

func relayLetterIndex(s string) int {
	for i := 1; i < len(s)-1; i++ {
		if s[i-1] == ' ' && s[i+1] == ' ' {
			if s[i] >= 'A' && s[i] <= 'Z' {
				return i
			}
		}
	}
	return -1
}
