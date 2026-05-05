package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestEmptyFilterMatchesAll(t *testing.T) {
	f := filter.New(filter.Options{})
	lines := []string{
		`2024-01-01T00:00:00Z level=info msg="hello world"`,
		`2024-01-01T00:00:01Z level=error msg="boom" service=api`,
		`plain text line without structure`,
	}
	for _, line := range lines {
		if !f.Match(line) {
			t.Errorf("empty filter should match %q", line)
		}
	}
}

func TestSubstringFilter(t *testing.T) {
	f := filter.New(filter.Options{Substring: "error"})
	if !f.Match(`level=error msg="something failed"`) {
		t.Error("expected match for line containing 'error'")
	}
	if f.Match(`level=info msg="all good"`) {
		t.Error("expected no match for line without 'error'")
	}
}

func TestFieldFilter(t *testing.T) {
	f := filter.New(filter.Options{
		Fields: map[string]string{"service": "api", "level": "error"},
	})
	if !f.Match(`ts=2024-01-01 level=error service=api msg="bad request"`) {
		t.Error("expected match when all fields present")
	}
	if f.Match(`ts=2024-01-01 level=error service=worker msg="task done"`) {
		t.Error("expected no match when service differs")
	}
	if f.Match(`ts=2024-01-01 level=info service=api msg="request ok"`) {
		t.Error("expected no match when level differs")
	}
}

func TestCombinedFilter(t *testing.T) {
	f := filter.New(filter.Options{
		Fields:    map[string]string{"level": "warn"},
		Substring: "timeout",
	})
	if !f.Match(`level=warn msg="connection timeout" host=db`) {
		t.Error("expected match for warn + timeout")
	}
	if f.Match(`level=warn msg="disk full" host=db`) {
		t.Error("expected no match: missing substring")
	}
	if f.Match(`level=error msg="connection timeout" host=db`) {
		t.Error("expected no match: wrong level field")
	}
}

func TestEmptyReportsCorrectly(t *testing.T) {
	empty := filter.New(filter.Options{})
	if !empty.Empty() {
		t.Error("expected Empty() == true for zero-value filter")
	}
	nonEmpty := filter.New(filter.Options{Substring: "x"})
	if nonEmpty.Empty() {
		t.Error("expected Empty() == false when substring set")
	}
}
