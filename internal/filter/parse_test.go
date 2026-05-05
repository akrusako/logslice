package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestParseFieldsValid(t *testing.T) {
	pairs := []string{"level=error", "service=api", "region=us-east-1"}
	fields, err := filter.ParseFields(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields["level"] != "error" {
		t.Errorf("level: got %q, want %q", fields["level"], "error")
	}
	if fields["service"] != "api" {
		t.Errorf("service: got %q, want %q", fields["service"], "api")
	}
	if fields["region"] != "us-east-1" {
		t.Errorf("region: got %q, want %q", fields["region"], "us-east-1")
	}
}

func TestParseFieldsValueWithEquals(t *testing.T) {
	// Values that contain '=' should be preserved intact.
	pairs := []string{"expr=a=b"}
	fields, err := filter.ParseFields(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields["expr"] != "a=b" {
		t.Errorf("expr: got %q, want %q", fields["expr"], "a=b")
	}
}

func TestParseFieldsInvalid(t *testing.T) {
	invalid := []string{"noequalssign", "=nokey", ""}
	for _, bad := range invalid {
		_, err := filter.ParseFields([]string{bad})
		if err == nil {
			t.Errorf("expected error for invalid pair %q", bad)
		}
	}
}

func TestParseFieldsNil(t *testing.T) {
	fields, err := filter.ParseFields(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields != nil {
		t.Error("expected nil map for nil input")
	}
}

func TestNewFromArgs(t *testing.T) {
	f, err := filter.NewFromArgs([]string{"level=info"}, "started")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match(`level=info msg="server started" port=8080`) {
		t.Error("expected match")
	}
	if f.Match(`level=info msg="server stopped" port=8080`) {
		t.Error("expected no match: substring missing")
	}
}
