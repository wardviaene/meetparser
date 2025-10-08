package parser

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var isvalidTime = regexp.MustCompile(`^(?:\d+|---)\s+(.+?),\s+(.+)`)

func ParsePDFText(filePath string) (Result, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Result{}, err
	}
	defer file.Close()
	return parsePDFText(file)
}
func parsePDFText(reader io.Reader) (Result, error) {
	var err error
	result := Result{
		Times:       []*SwimmerTime{},
		RelayTimes:  []*RelayTime{},
		Events:      []*Event{},
		ParseErrors: []*ParseError{},
	}
	fileType := ""
	pageRegex := regexp.MustCompile(`Page \d+$`)

	scanner := bufio.NewScanner(reader)
	processIndividual := false
	processRelay := false
	var event *Event

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if i == 0 {
			if strings.TrimSpace(line) == "FileType: SwimTopia Meet Maestro" {
				fileType = "SwimTopia Meet Maestro"
			}
		}
		if (processIndividual || processRelay) && (line == " " || line == "" || pageRegex.MatchString(line)) {
			processIndividual = false
			processRelay = false
		}
		if processIndividual {
			if line == "Preliminaries" {
				event.Type = "Preliminaries"
			} else if strings.HasSuffix(line, "Final") {
				event.Type = line
			} else if strings.HasSuffix(line, "Swim-off") {
				event.Type = "swim-off"
			} else if strings.HasSuffix(line, "Swim-Off Required") {
				event.Type = "swim-Off required"
			} else {
				if isvalidTime.MatchString(line) {
					swimmerTime, err := processLine(line, fileType)
					if err != nil {
						parseError := ParseError{
							Type:               "IndividualTime",
							PartialSwimmerTime: swimmerTime,
							LineNumber:         i,
							Line:               line,
							ErrorMessage:       err.Error(),
						}
						result.ParseErrors = append(result.ParseErrors, &parseError)
					} else {
						if event == nil || event.Round == "" {
							parseError := ParseError{
								Type:         "IndividualTime",
								LineNumber:   i,
								Line:         line,
								ErrorMessage: "event number is empty",
							}
							result.ParseErrors = append(result.ParseErrors, &parseError)
						}
						swimmerTime.Event = event
						result.Times = append(result.Times, swimmerTime)
					}
				} else if splitTimesRegex.MatchString(line) && len(result.Times) > 0 {
					splitTimes := getSplitTimes(line)
					result.Times[len(result.Times)-1].SplitTimes = splitTimes
				}
			}
		} else if processRelay {
			if isRelaySwimmerLine(line) {
				relaySwimmers, err := processRelaySwimmersLine(line, fileType)
				if err != nil {
					parseError := ParseError{
						Type:         "RelaySwimmer",
						LineNumber:   i,
						Line:         line,
						ErrorMessage: err.Error(),
					}
					result.ParseErrors = append(result.ParseErrors, &parseError)
				} else {
					result.RelayTimes[len(result.RelayTimes)-1].Swimmers = append(result.RelayTimes[len(result.RelayTimes)-1].Swimmers, relaySwimmers...)
				}
			} else {
				if startsWithNumber(line) {
					relayTime, err := processRelayLine(line, fileType)
					if err != nil {
						parseError := ParseError{
							Type:         "RelayTime",
							LineNumber:   i,
							Line:         line,
							ErrorMessage: err.Error(),
						}
						result.ParseErrors = append(result.ParseErrors, &parseError)
					}
					relayTime.Event = event
					result.RelayTimes = append(result.RelayTimes, relayTime)
				}
			}
		}

		if isEvent(line, fileType) {
			event, err = processEvent(line, fileType)
			if err != nil {
				parseError := ParseError{
					Type:         "Event",
					LineNumber:   i,
					Line:         line,
					ErrorMessage: err.Error(),
				}
				result.ParseErrors = append(result.ParseErrors, &parseError)
			} else {
				result.Events = append(result.Events, event)
			}
		} else if strings.Contains(line, "Name Age") || strings.Contains(line, "Name Ag  e") || strings.Contains(line, "Name Ag\te") {
			processIndividual = true
		} else if strings.Contains(line, "Team  Relay") || (fileType == FILETYPE_TYPE2 && strings.Contains(line, "Pl Team Relay")) {
			processIndividual = false
			processRelay = true
		} else if strings.Contains(line, "Qualifying Times") {
			if event != nil {
				err = eventAddQualifyingTimes(result.Events[len(result.Events)-1], line)
				if err != nil {
					parseError := ParseError{
						Type:         "Event",
						LineNumber:   i,
						Line:         line,
						ErrorMessage: "Qualifying Times: " + err.Error(),
					}
					result.ParseErrors = append(result.ParseErrors, &parseError)
				}
			}
		}

	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func isRelaySwimmerLine(line string) bool {
	return strings.HasPrefix(line, "1)") || strings.HasPrefix(line, "2)") || strings.HasPrefix(line, "3)") || strings.HasPrefix(line, "4)")
}

func isEvent(line string, fileType string) bool {
	switch fileType {
	case FILETYPE_TYPE2:
		return len(line) > 1 && strings.HasPrefix(line, "#") && isNumeric(line[1:2])
	default:
		return strings.HasPrefix(line, "Event") || strings.HasPrefix(line, "event") || len(line) > 1 && strings.HasPrefix(line, "#") && isNumeric(line[1:2])
	}
}

func splitLastXChars(myvar string, length int) (var1, var2 string) {
	if len(myvar) == 0 {
		return "", ""
	}
	runes := []rune(myvar) // handle Unicode safely
	var1 = string(runes[:len(runes)-length])

	last := runes[len(runes)-length:]
	var2 = string(last)

	return
}

func startsWithNumber(line string) bool {
	if len(line) == 0 {
		return false
	}
	r := []rune(line)[0]
	return unicode.IsDigit(r)
}
