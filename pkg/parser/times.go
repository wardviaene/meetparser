package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var timesRegex = regexp.MustCompile(`(?:\d{1,2}:)?\d{2}\.\d{2}`)
var timesNoSeed = regexp.MustCompile(` (?:\d{1,2}:)?\d{2}\.\d{2}`)
var judgedTime = regexp.MustCompile(` ?J?(?:\d{1,2}:)?\d{2}\.\d{2}`)
var timesNTRegex = regexp.MustCompile(`NT (?:\d{1,2}:)?\d{2}\.\d{2}`)
var timesDQRegex = regexp.MustCompile(`(?:\d{1,2}:)?\d{2}\.\d{2}(?: [YLS])? DQ`)
var timesNSRegex = regexp.MustCompile(`(?:\d{1,2}:)?\d{2}\.\d{2}(?: [YLS])? NS`)
var timesDFSRegex = regexp.MustCompile(`(?:\d{1,2}:)?\d{2}\.\d{2}(?: [YLS])? DFS`)
var splitTimesRegex = regexp.MustCompile(`^(?:\d{1,2}:)?\d{2}\.\d{2}(?:\s+(?:\d{1,2}:)?\d{2}\.\d{2})*$`)

func processTimes(line string) (int, string, string, error) {
	matchedTimes := timesRegex.FindAllStringIndex(line, -1)
	if len(matchedTimes) == 2 { // found both seed and actual
		return matchedTimes[0][0], line[matchedTimes[0][0]:matchedTimes[0][1]], line[matchedTimes[1][0]:matchedTimes[1][1]], nil
	}

	if len(matchedTimes) == 0 || len(matchedTimes) == 1 {
		return getNoTimeCodes(line, len(matchedTimes))
	}
	fmt.Printf("more than 2 times? %v\n", matchedTimes)
	return -1, "", "", fmt.Errorf("not supported (found more than 2 times)")

}
func getNoTimeCodes(line string, timesFound int) (int, string, string, error) {
	if timesFound == 0 {
		// line: Virginia Gators NT SCR
		// line: Nation's Capital Swim Club DFS
		indexAfterTeamName := strings.Index(line, "NT DQ")
		if indexAfterTeamName != -1 {
			return indexAfterTeamName, "NT", "DQ", nil
		}
		indexAfterTeamName = strings.Index(line, "NT NS")
		if indexAfterTeamName != -1 {
			return indexAfterTeamName, "NT", "NS", nil
		}
		indexAfterTeamName = strings.Index(line, "NT DNF")
		if indexAfterTeamName != -1 {
			return indexAfterTeamName, "NT", "DNF", nil
		}
		indexAfterTeamName = strings.Index(line, "NT SCR")
		if indexAfterTeamName != -1 {
			return indexAfterTeamName, "NT", "SCR", nil
		}
		indexAfterTeamName = strings.Index(line, " DFS")
		if indexAfterTeamName != -1 {
			return indexAfterTeamName, "", "DFS", nil
		}
		indexAfterTeamName = strings.Index(line, " DQ")
		if indexAfterTeamName != -1 {
			return indexAfterTeamName, "", "DQ", nil
		}
		indexAfterTeamName = strings.Index(line, " DNF")
		if indexAfterTeamName != -1 {
			return indexAfterTeamName, "", "DNF", nil
		}
		return -1, "", "", fmt.Errorf("no codes recognized instead of times")
	}
	if timesFound == 1 {
		// potentially one time matches
		// line: Lynchburg YMCA NT 33.49 2
		// line: Lynchburg YMCA 33.49 DQ 2
		// line: Lynchburg YMCA 2:54.90 NS
		// line: Dads Club Swim Team-GU 2:13.68 Y DQ
		// line: Lakeside Aquatic Club-NT 4:55.32 DFS
		if matchedDQTimes := timesDQRegex.FindStringIndex(line); matchedDQTimes != nil {
			return matchedDQTimes[0], line[matchedDQTimes[0] : matchedDQTimes[1]-3], "DQ", nil
		}
		if matchedNSTimes := timesNSRegex.FindStringIndex(line); matchedNSTimes != nil {
			return matchedNSTimes[0], line[matchedNSTimes[0] : matchedNSTimes[1]-3], "DQ", nil
		}
		if matchedDFSTimes := timesDFSRegex.FindStringIndex(line); matchedDFSTimes != nil {
			return matchedDFSTimes[0], line[matchedDFSTimes[0] : matchedDFSTimes[1]-3], "DFS", nil
		}
		if matchedNTTimes := timesNTRegex.FindStringIndex(line); matchedNTTimes != nil {
			if matchedNTTimes[0] == 0 || line[matchedNTTimes[0]-1:matchedNTTimes[0]] != "-" {
				return matchedNTTimes[0], "NT", line[matchedNTTimes[0]+3 : matchedNTTimes[1]], nil
			}
		}
		// or no seed time:
		// SJAC-MA 1:09.33

		if matchedNoSeedTimes := timesNoSeed.FindStringIndex(line); matchedNoSeedTimes != nil {
			return matchedNoSeedTimes[0], "", line[matchedNoSeedTimes[0]+1 : matchedNoSeedTimes[1]], nil
		}
		// or judged:
		// J1:09.33
		if matchedJudgedTimes := judgedTime.FindStringIndex(line); matchedJudgedTimes != nil {
			return matchedJudgedTimes[0], "", line[matchedJudgedTimes[0]+1 : matchedJudgedTimes[1]], nil
		}

		return -1, "", "", fmt.Errorf("code not supported (found one time, but no DQ/NT/NS)")
	}
	return -1, "", "", fmt.Errorf("no codes recognized instead of times (too many elements supplied)")
}

func getSplitTimes(line string) []string {
	return timesRegex.FindAllString(line, -1)
}
