package parser

import "fmt"

const FILETYPE_TYPE1 = ""
const FILETYPE_TYPE2 = "SwimTopia Meet Maestro"

type Result struct {
	Events      []*Event       `json:"events"`
	Times       []*SwimmerTime `json:"times"`
	RelayTimes  []*RelayTime   `json:"relayTimes"`
	ParseErrors []*ParseError  `json:"parseErrors"`
}

type Event struct {
	Round           string            `json:"round"`
	Type            string            `json:"type"`
	Gender          string            `json:"gender"`
	AgeGroup        string            `json:"ageGroup"`
	Distance        string            `json:"distance"`
	Stroke          string            `json:"stroke"`
	Relay           bool              `json:"relay"`
	QualifyingTimes map[string]string `json:"qualifyingTimes"`
}

type RelayTime struct {
	Event               *Event          `json:"event"`
	Place               string          `json:"place"`
	TeamName            string          `json:"teamName"`
	TeamNameShort       string          `json:"teamNameShort,omitempty"`
	TeamLSC             string          `json:"teamLSC"`
	RelayEntry          string          `json:"relay"`
	Time                string          `json:"time"`
	SeedTime            string          `json:"seedTime"`
	SeedTimeTag         string          `json:"seedTimeTag"`
	QualifyingStandards string          `json:"qualifyingStandards"`
	Points              string          `json:"points"`
	Achievements        string          `json:"achievements,omitempty"`
	Swimmers            []*RelaySwimmer `json:"swimmers"`
}
type RelaySwimmer struct {
	Place string `json:"place"`
	Name  string `json:"name"`
	Age   string `json:"age"`
}
type SwimmerTime struct {
	Event               *Event   `json:"event"`
	Place               string   `json:"place"`
	Age                 string   `json:"age"`
	Name                string   `json:"name"`
	TeamName            string   `json:"teamName"`
	TeamLSC             string   `json:"teamLSC"`
	Finals              string   `json:"finals"`
	Time                string   `json:"time"`
	SeedTime            string   `json:"seedTime"`
	SeedTimeTag         string   `json:"seedTimeTag"`
	Points              string   `json:"points"`
	QualifyingStandards string   `json:"qualifyingStandards"`
	Achievements        string   `json:"achievements,omitempty"`
	SplitTimes          []string `json:"splitTimes,omitempty"`
}

func (e *Event) String() string {
	return fmt.Sprintf("Event: Round: '%s', Gender: '%s', AgeGroup: '%s', Distance: '%s', Stroke: '%s'",
		e.Round,
		e.Gender,
		e.AgeGroup,
		e.Distance,
		e.Stroke,
	)
}

func (s *SwimmerTime) String() string {
	return fmt.Sprintf("Team: '%s', LSC: '%s', Age: '%s', Name: '%s', Time: '%s', Seed Time: '%s', Seed Time Tag: '%s', Points: '%s', Place: '%s', Qualifying Standards: '%s'",
		s.TeamName,
		s.TeamLSC,
		s.Age,
		s.Name,
		s.Time,
		s.SeedTime,
		s.SeedTimeTag,
		s.Points,
		s.Place,
		s.QualifyingStandards,
	)
}

type ParseError struct {
	Type               string       `json:"type"`
	LineNumber         int          `json:"lineNumber"`
	Line               string       `json:"line"`
	ErrorMessage       string       `json:"errorMessage"`
	PartialSwimmerTime *SwimmerTime `json:"partialSwimmerTime"`
}
