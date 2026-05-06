package fieldmatch_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldmatch"
)

func TestNoRulesMatchesAll(t *testing.T) {
	m, err := fieldmatch.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(`level=info msg="hello world"`) {
		t.Error("expected empty matcher to match every line")
	}
}

func TestBadPatternReturnsError(t *testing.T) {
	_, err := fieldmatch.New(map[string]string{"level": "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatchKVLine(t *testing.T) {
	m, err := fieldmatch.New(map[string]string{"level": "^(warn|error)$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(`level=error msg="disk full"`) {
		t.Error("expected match for level=error")
	}
	if m.Match(`level=info msg="all good"`) {
		t.Error("expected no match for level=info")
	}
}

func TestMatchJSONLine(t *testing.T) {
	m, err := fieldmatch.New(map[string]string{"service": "auth"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(`{"service":"auth","level":"info"}`) {
		t.Error("expected match for service=auth")
	}
	if m.Match(`{"service":"billing","level":"info"}`) {
		t.Error("expected no match for service=billing")
	}
}

func TestMissingFieldDoesNotMatch(t *testing.T) {
	m, err := fieldmatch.New(map[string]string{"user": "alice"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Match(`level=info msg="no user field here"`) {
		t.Error("expected no match when required field is absent")
	}
}

func TestAllRulesMustMatch(t *testing.T) {
	m, err := fieldmatch.New(map[string]string{
		"level":   "error",
		"service": "auth",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// only level matches
	if m.Match(`{"level":"error","service":"billing"}`) {
		t.Error("expected no match when only one rule is satisfied")
	}
	// both match
	if !m.Match(`{"level":"error","service":"auth"}`) {
		t.Error("expected match when all rules are satisfied")
	}
}

func TestMatchedFieldsReturnsSubset(t *testing.T) {
	m, err := fieldmatch.New(map[string]string{
		"level": "error",
		"user":  "bob",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.MatchedFields(`{"level":"error","service":"auth"}`)
	if _, ok := got["level"]; !ok {
		t.Error("expected level in matched fields")
	}
	if _, ok := got["user"]; ok {
		t.Error("expected user absent from matched fields")
	}
}
