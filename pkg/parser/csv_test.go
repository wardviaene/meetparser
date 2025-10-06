package parser

import "testing"

func TestMarshalCSV(t *testing.T) {
	result := &Result{
		Times: []*SwimmerTime{
			{
				Event: &Event{
					Round: "1",
				},
			},
		},
	}
	_, err := MarshalCSV(result.Times)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}
