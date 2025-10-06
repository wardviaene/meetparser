package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func processLine(line string, fileType string) (*SwimmerTime, error) {
	switch fileType {
	case FILETYPE_TYPE2:
		return processLineType2(line)
	default:
		return processLineType1(line)
	}
}

func processLineType2(line string) (*SwimmerTime, error) {
	swimmer := &SwimmerTime{}
	// line: 1 Lastname, Firstname 6 PFP 18.14 18.39
	index1 := strings.Index(line, " ")
	if index1 == -1 {
		return swimmer, fmt.Errorf("couldn't determine place")
	}
	swimmer.Place = line[0:index1]
	line = line[index1+1:]
	// line: Lastname, Firstname 6 PFP 18.14 18.39
	index2 := stringAgeIndex(line)
	if index2 == -1 {
		return swimmer, fmt.Errorf("couldn't determine age/name position")
	}
	swimmer.Name = line[0 : index2-1]
	line = line[index2:]
	// line: 6 PFP 18.14 18.39
	index3 := strings.Index(line, " ")
	if index1 == -1 {
		return swimmer, fmt.Errorf("couldn't determine age")
	}
	swimmer.Age = line[0:index3]
	line = line[index3+1:]
	// line: PFP 18.14 18.39
	index4 := strings.Index(line, " ")
	if index1 == -1 {
		return swimmer, fmt.Errorf("couldn't determine team name")
	}
	swimmer.TeamName = line[0:index4]
	line = line[index4+1:]
	// line: 18.14 18.39
	index5 := strings.Index(line, " ")
	if index1 == -1 {
		return swimmer, fmt.Errorf("couldn't determine seed time")
	}
	swimmer.SeedTime = line[0:index5]
	line = line[index5+1:]
	// line: 18.39
	index6 := strings.Index(line, " ")
	if index6 == -1 {
		swimmer.Time = line
		line = ""
	} else {
		swimmer.Time = line[0:index6]
		line = line[index6+1:]
	}

	// line: SWT
	// line: 9 INV
	if line != "" {
		index7 := strings.Index(line, " ")
		if index7 == -1 { // line: SWT    or   line: 9
			if !isNumeric(line) {
				swimmer.Achievements = line
			} else {
				swimmer.Points = line
			}
			line = ""
		} else {
			if !isNumeric(line[0:index7]) { // line: INV Somethingelse
				swimmer.Achievements = line[0:index7]
				line = line[index7+1:]
			} else { // line: 9 INV
				swimmer.Points = line[0:index7]
				line = line[index7+1:]
			}
		}
	}

	if line != "" && swimmer.Achievements == "" {
		swimmer.Achievements = line
		line = ""
	}

	err := checkResidual(swimmer.SeedTime + " " + swimmer.Time + " " + line)
	if err != nil {
		return swimmer, fmt.Errorf("residual information found: '%s'", err)
	}

	return swimmer, nil
}

func processLineType1(line string) (*SwimmerTime, error) {
	swimmer := &SwimmerTime{}
	// line: 1 Lastname, Firstname  14 Lynchburg YMCA 2:14.96 2:16.72 AG 9
	index1 := strings.Index(line, " ")
	if index1 == -1 {
		return swimmer, fmt.Errorf("couldn't determine place")
	}
	swimmer.Place = line[0:index1]
	line = line[index1+1:]
	// line: Lastname, Firstname  14 Lynchburg YMCA 2:14.96 2:16.72 AG 9
	index2 := strings.Index(line, "   ")
	index2Offset := 3
	if index2 == -1 {
		index2 = strings.Index(line, "  ")
		if index2 == -1 {
			return swimmer, fmt.Errorf("error while parsing swimmer name")
		}
		index2Offset = 2
	}
	swimmer.Name = line[0:index2]
	line = line[index2+index2Offset:]
	// line: 14 Lynchburg YMCA 2:14.96 2:16.72 AG 9
	index3 := strings.Index(line, " ")
	if index1 == -1 {
		return swimmer, fmt.Errorf("couldn't determine age")
	}
	swimmer.Age = line[0:index3]
	line = line[index3+1:]
	// line: Lynchburg YMCA 2:14.96 2:16.72 AG 9
	indexAfterTeamName := 0
	var err error
	indexAfterTeamName, swimmer.SeedTime, swimmer.Time, err = processTimes(line)
	if err != nil {
		fmt.Printf("Line: %s\n", line)
		return swimmer, fmt.Errorf("process time error: %s", err)
	}
	teamName := strings.TrimSpace(line[0:indexAfterTeamName])
	if index := strings.Index(teamName, "-"); index != -1 {
		swimmer.TeamLSC = teamName[index+1:]
		swimmer.TeamName = teamName[0:index]
	} else {
		swimmer.TeamName = teamName
	}
	line = line[indexAfterTeamName:]
	// Extract points
	// line: 2:14.96 2:16.72 AG 9
	// line: 2:14.96 2:16.72
	// line:   2:16.72
	if !strings.HasSuffix(line, "  ") {
		pointsIndex := strings.LastIndex(line, " ")
		if isNumeric(line[pointsIndex+1:]) {
			swimmer.Points = line[pointsIndex+1:]
			line = line[0:pointsIndex]
		}
	}

	// line can be empty
	if line == "" {
		residual := swimmer.Time
		if swimmer.SeedTime != "" {
			residual = swimmer.SeedTime + " " + swimmer.Time
		}
		err = checkResidual(residual)
		if err != nil {
			return swimmer, fmt.Errorf("residual information found: '%s'", err)
		}
		err = validateSwimmer(swimmer)
		if err != nil {
			return swimmer, fmt.Errorf("invalid swimmer data: '%s'", err)
		}
		return swimmer, nil
	}

	// Extract qualifying standards
	// line: "10:49.69 Y 9:43.85 TAGS  "
	line = strings.TrimSpace(line) // remove unnecessary spacing
	qualifyingStandardsIndex := strings.LastIndex(line, " ")
	if !timesRegex.MatchString(line[qualifyingStandardsIndex+1:]) {
		if line[qualifyingStandardsIndex+1:] != "DQ" && line[qualifyingStandardsIndex+1:] != "NS" && line[qualifyingStandardsIndex+1:] != "DNF" && line[qualifyingStandardsIndex+1:] != "DFS" {
			swimmer.QualifyingStandards = line[qualifyingStandardsIndex+1:]
			line = line[0:qualifyingStandardsIndex]
		}
	}
	// line: 10:43.41 Y 9:29.11
	seedTagIndex := strings.Index(line, " ")
	if index1 == -1 {
		return swimmer, fmt.Errorf("couldn't determine seed time tag")
	}
	if strings.HasPrefix(line[seedTagIndex+1:], "Y ") || strings.HasPrefix(line[seedTagIndex+1:], "S ") || strings.HasPrefix(line[seedTagIndex+1:], "L ") {
		swimmer.SeedTimeTag = line[seedTagIndex+1 : seedTagIndex+2]
		line = line[0:seedTagIndex+1] + line[seedTagIndex+3:]
	}

	err = checkResidual(line)
	if err != nil {
		return swimmer, fmt.Errorf("residual information found: '%s'", err)
	}

	err = validateSwimmer(swimmer)
	if err != nil {
		return swimmer, fmt.Errorf("invalid swimmer data: '%s'", err)
	}

	return swimmer, nil
}

func validateSwimmer(swimmer *SwimmerTime) error {
	if swimmer.Points != "" && !isNumeric(swimmer.Points) {
		return fmt.Errorf("points is not numeric")
	}
	if swimmer.Time != "" && !timesRegex.MatchString(swimmer.Time) {
		if !isValidTimeCode(swimmer.Time) {
			return fmt.Errorf("final time is not a time string / valid code")
		}
	}
	if swimmer.SeedTime != "" && !timesRegex.MatchString(swimmer.SeedTime) {
		if !isValidSeedTimeCode(swimmer.SeedTime) {
			return fmt.Errorf("seed time is not a time string")
		}
	}
	if swimmer.Age != "" && !isNumeric(swimmer.Age) {
		return fmt.Errorf("age is not numeric")
	}
	return nil
}

func isValidTimeCode(s string) bool {
	if s == "DQ" || s == "NS" || s == "DNF" || s == "DFS" || s == "SCR" {
		return true
	}
	return false
}
func isValidSeedTimeCode(s string) bool {
	if s == "NT" {
		return true
	}
	return false
}

func checkResidual(residual string) error {
	split := strings.Split(strings.TrimSpace(residual), " ")
	if len(split) != 2 && len(split) != 1 {
		return fmt.Errorf("zero parts or more than 2 parts found: '%s'", residual)
	}
	if !timesRegex.MatchString(split[0]) && split[0] != "NT" {
		return fmt.Errorf("unknown value (first part): '%s'", split[0])
	}
	if len(split) >= 2 {
		if !timesRegex.MatchString(split[1]) && split[1] != "DQ" && split[1] != "NS" && split[1] != "DNF" && split[1] != "DFS" {
			return fmt.Errorf("unknown value (second part): '%s'", split[1])
		}
	}
	return nil
}

func stringAgeIndex(s string) int {
	for byteIdx := 0; byteIdx < len(s); {
		r, size := utf8.DecodeRuneInString(s[byteIdx:])

		// skip first character since it can't have a space before it
		if byteIdx > 0 {
			// check if preceding byte is a space (safe because space is 1 byte)
			if s[byteIdx-1] == ' ' && unicode.IsDigit(r) {
				// look ahead to count up to 3 digits
				digitBytes := size
				nextIdx := byteIdx + size
				for count := 1; count < 3 && nextIdx < len(s); count++ {
					nextR, nextSize := utf8.DecodeRuneInString(s[nextIdx:])
					if !unicode.IsDigit(nextR) {
						break
					}
					digitBytes += nextSize
					nextIdx += nextSize
				}
				// check for trailing space
				if nextIdx < len(s) {
					nextR, _ := utf8.DecodeRuneInString(s[nextIdx:])
					if nextR == ' ' {
						return byteIdx
					}
				}
			}
		}

		byteIdx += size
	}
	return -1
}

func isNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
