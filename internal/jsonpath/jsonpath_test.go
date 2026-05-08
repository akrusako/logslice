package jsonpath

import (
	"testing"
)

func TestNewEmptyPathReturnsError(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewInvalidSegmentReturnsError(t *testing.T) {
	_, err := New("a..b")
	if err == nil {
		t.Fatal("expected error for path with empty segment")
	}
}

func TestGetTopLevelString(t *testing.T) {
	e, _ := New("level")
	v, ok := e.Get(`{"level":"info","msg":"hello"}`)
	if !ok || v != "info" {
		t.Fatalf("expected info, got %q ok=%v", v, ok)
	}
}

func TestGetNestedField(t *testing.T) {
	e, _ := New("request.method")
	v, ok := e.Get(`{"request":{"method":"GET","path":"/api"}}`)
	if !ok || v != "GET" {
		t.Fatalf("expected GET, got %q ok=%v", v, ok)
	}
}

func TestGetNumericField(t *testing.T) {
	e, _ := New("status")
	v, ok := e.Get(`{"status":200}`)
	if !ok || v != "200" {
		t.Fatalf("expected 200, got %q ok=%v", v, ok)
	}
}

func TestGetBoolField(t *testing.T) {
	e, _ := New("ok")
	v, ok := e.Get(`{"ok":true}`)
	if !ok || v != "true" {
		t.Fatalf("expected true, got %q ok=%v", v, ok)
	}
}

func TestGetMissingField(t *testing.T) {
	e, _ := New("missing")
	_, ok := e.Get(`{"level":"info"}`)
	if ok {
		t.Fatal("expected not found for missing field")
	}
}

func TestGetNonJSONLine(t *testing.T) {
	e, _ := New("level")
	_, ok := e.Get("plain text log line")
	if ok {
		t.Fatal("expected not found for non-JSON line")
	}
}

func TestGetEmptyLine(t *testing.T) {
	e, _ := New("level")
	_, ok := e.Get("")
	if ok {
		t.Fatal("expected not found for empty line")
	}
}

func TestPathReturnsOriginal(t *testing.T) {
	e, _ := New("a.b.c")
	if e.Path() != "a.b.c" {
		t.Fatalf("expected a.b.c, got %q", e.Path())
	}
}

func TestGetDeeplyNested(t *testing.T) {
	e, _ := New("a.b.c")
	v, ok := e.Get(`{"a":{"b":{"c":"deep"}}}`)
	if !ok || v != "deep" {
		t.Fatalf("expected deep, got %q ok=%v", v, ok)
	}
}
