package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var ageRegex = regexp.MustCompile(`(?i)^(?:\d{1,2}\s*&\s*(?:Under|Over|O)|\d{1,2}\s*-\s*\d{1,2})`)
var distanceRegex = regexp.MustCompile(`(?i)^(?:\d+\s*(?:LC|SC)?\s*Meter|\d+\s*(?:Yard|yd))\b`)

func processEvent(line string, fileType string) (*Event, error) {
	switch fileType {
	case FILETYPE_TYPE2:
		return processEventType2(line)
	default:
		return processEventType1(line)
	}
}
func processEventType2(line string) (*Event, error) {
	event := &Event{
		QualifyingTimes: make(map[string]string),
	}
	var err error
	// line: #1 Mixed 6 & Under 100yd Freestyle Relay
	index1 := strings.Index(line, " ")
	if index1 == -1 {
		return event, fmt.Errorf("event number not found")
	}
	event.Round = strings.TrimLeft(line[0:index1], "#")
	line = line[index1+1:]
	// line: Mixed 6 & Under 100yd Freestyle Relay
	index2 := strings.Index(line, " ")
	if index2 == -1 {
		return event, fmt.Errorf("unexpected end of string (expected gender)")
	}
	event.Gender, err = parseGender(line[0:index2], true /* exact */)
	if err != nil {
		return event, fmt.Errorf("couldn't parse gender: %s", err)
	}
	line = line[index2+1:]
	// line: 6 & Under 100yd Freestyle Relay
	event.AgeGroup, err = parseEventAgeGroup(line)
	if err != nil {
		return nil, fmt.Errorf("parse age group error: %s", err)
	}
	if event.AgeGroup != "" { // we might have an empty age group if the event doesn't have it
		line = line[len(event.AgeGroup)+1:]
	}
	event.AgeGroup = normalizeAge(event.AgeGroup)
	// line: 100yd Freestyle Relay
	event.Distance, err = parseEventDistance(line)
	if err != nil {
		return nil, fmt.Errorf("parse distance error: %s", err)
	}
	line = line[len(event.Distance)+1:]
	// line: 100yd Freestyle Relay
	event.Stroke, event.Relay, err = parseStroke(line)
	if err != nil {
		return nil, fmt.Errorf("parse stroke error: %s", err)
	}

	return event, nil

}
func processEventType1(line string) (*Event, error) {
	line = strings.ReplaceAll(line, "\t", " ")
	event := &Event{
		QualifyingTimes: make(map[string]string),
	}
	if strings.HasPrefix(line, "(Event") || strings.HasPrefix(line, "(event") {
		line = line[len("(Event")+1:]
	} else if strings.HasPrefix(line, "Event") || strings.HasPrefix(line, "event") {
		line = line[len("Event")+1:]
	} else if strings.HasPrefix(line, "#") {
		line = line[len("#"):]
	} else {
		return nil, fmt.Errorf("can't find event. Output: %s", line)
	}
	eventSplit := strings.Split(line, "  ")
	if len(eventSplit) < 2 {
		return nil, fmt.Errorf("event # not found")
	}
	// eventSplit[0]: 57
	event.Round = eventSplit[0]
	// eventSplit[1]: Girls 10 & Under 50 LC Meter Butterfly)
	index := 0
	var err error
	event.Gender, err = parseGender(eventSplit[1], false)
	if err != nil {
		return nil, fmt.Errorf("can't extract gender from: %s", line)
	}
	index += len(event.Gender) + 1
	// eventSplit[1]: 10 & Under 50 LC Meter Butterfly)
	event.AgeGroup, err = parseEventAgeGroup(eventSplit[1][index:])
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	if event.AgeGroup != "" { // we might have an empty age group if the event doesn't have it
		index += len(event.AgeGroup) + 1
	}
	event.AgeGroup = normalizeAge(event.AgeGroup)
	event.Distance, err = parseEventDistance(eventSplit[1][index:])
	if err != nil {
		return nil, fmt.Errorf("can't extract distance from from: %s", eventSplit[1][index:])
	}
	index += len(event.Distance) + 1
	// eventSplit[1]: Butterfly)
	event.Stroke, event.Relay, err = parseStroke(eventSplit[1][index:])
	if err != nil {
		return nil, fmt.Errorf("can't extract stroke from: %s", err)
	}
	return event, nil
}

func parseStroke(stroke string) (string, bool, error) {
	if strings.HasPrefix(stroke, "Medley Relay") {
		return "Medley", true, nil
	}
	if strings.HasPrefix(stroke, "Butterfly Relay") {
		return "Butterfly", true, nil
	}
	if strings.HasPrefix(stroke, "Breaststroke Relay") {
		return "Breaststroke", true, nil
	}
	if strings.HasPrefix(stroke, "Freestyle Relay") {
		return "Freestyle", true, nil
	}
	if strings.HasPrefix(stroke, "IM Relay") {
		return "IM", true, nil
	}
	if strings.HasPrefix(stroke, "Backstroke Relay") {
		return "Backstroke", true, nil
	}
	if strings.HasPrefix(stroke, "Butterfly") {
		return "Butterfly", false, nil
	}
	if strings.HasPrefix(stroke, "Butter ly") {
		return "Butterfly", false, nil
	}
	if strings.HasPrefix(stroke, "Fly") {
		return "Fly", false, nil
	}
	if strings.HasPrefix(stroke, "Back") {
		return "Back", false, nil
	}
	if strings.HasPrefix(stroke, "Breaststroke") {
		return "Breaststroke", false, nil
	}
	if strings.HasPrefix(stroke, "Freestyle") {
		return "Freestyle", false, nil
	}
	if strings.HasPrefix(stroke, "IM") {
		return "IM", false, nil
	}
	if strings.HasPrefix(stroke, "Backstroke") {
		return "Backstroke", false, nil
	}

	return "", false, fmt.Errorf("stroke not found: %s", stroke)

}

func parseGender(gender string, exact bool) (string, error) {
	if !exact {
		if strings.HasPrefix(gender, "Girls ") {
			return "girls", nil
		}
		if strings.HasPrefix(gender, "Boys ") {
			return "boys", nil
		}
		if strings.HasPrefix(gender, "Women ") {
			return "women", nil
		}
		if strings.HasPrefix(gender, "Men ") {
			return "men", nil
		}
		if strings.HasPrefix(gender, "Mixed ") {
			return "mixed", nil
		}
	} else {
		switch strings.ToLower(gender) {
		case "boys":
			return "boys", nil
		case "girls":
			return "girls", nil
		case "women":
			return "women", nil
		case "men":
			return "men", nil
		case "mixed":
			return "mixed", nil
		default:
			return "", fmt.Errorf("gender not found (no match): %s", gender)
		}
	}

	return "", fmt.Errorf("gender not found: %s", gender)
}
func parseEventAgeGroup(data string) (string, error) {
	match := ageRegex.FindString(data)
	if match != "" {
		return strings.ToLower(match), nil
	}

	return "", nil // some events don't have an age group
	//return "", fmt.Errorf("can't extract age group from: %s", data)
}

func normalizeAge(s string) string {
	out := s
	if strings.HasSuffix(out, "&o") {
		out = strings.ReplaceAll(out, "&o", "")
		return out + " & over"
	}
	return out
}

func parseEventDistance(data string) (string, error) {
	// eventSplit[1]: 50 LC Meter Butterfly)
	match := distanceRegex.FindString(data)
	if match != "" {
		return match, nil
	}
	return "", fmt.Errorf("can't extract distance from: %s", data)
}

func eventAddQualifyingTimes(event *Event, line string) error {
	// INV NWSC Invitational Meet Qualifying Times '25 (Girls 6&U) 28.51
	index := strings.Index(line, "Qualifying Times")
	if index == -1 || index == 0 {
		return nil // couldn't find any qualifying times
	}
	name := line[0 : index-1]
	matchedTimes := timesRegex.FindAllStringIndex(line, -1)
	for _, matchedTime := range matchedTimes {
		event.QualifyingTimes[name] = line[matchedTime[0]:matchedTime[1]]
	}

	return nil
}
