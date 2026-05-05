package timeparse_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/timeparse"
)

type parseCase struct {
	input   string
	wantErr bool
	wantUTC string // expected time formatted as RFC3339, empty means skip check
}

var cases = []parseCase{
	{input: "2024-03-15T12:00:00Z", wantUTC: "2024-03-15T12:00:00Z"},
	{input: "2024-03-15T12:00:00.123456789Z", wantUTC: "2024-03-15T12:00:00Z"},
	{input: "2024-03-15 12:00:00", wantUTC: "2024-03-15T12:00:00Z"},
	{input: "2024-03-15T12:00:00", wantUTC: "2024-03-15T12:00:00Z"},
	{input: "15/Mar/2024:12:00:00 +0000", wantUTC: "2024-03-15T12:00:00Z"},
	{input: "not-a-timestamp", wantErr: true},
	{input: "", wantErr: true},
}

func TestParseAny(t *testing.T) {
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := timeparse.ParseAny(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if tc.wantUTC != "" {
				got = got.UTC().Truncate(time.Second)
				want, _ := time.Parse(time.RFC3339, tc.wantUTC)
				if !got.Equal(want) {
					t.Errorf("input %q: got %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

func TestParserCaching(t *testing.T) {
	p := timeparse.NewParser()
	inputs := []string{
		"2024-01-01T00:00:00Z",
		"2024-06-15T08:30:00Z",
		"2024-12-31T23:59:59Z",
	}
	for _, raw := range inputs {
		if _, err := p.Parse(raw); err != nil {
			t.Errorf("cached parser failed on %q: %v", raw, err)
		}
	}
}
