package fieldextract

import (
	"testing"
)

func TestExtractJSON(t *testing.T) {
	line := `{"level":"info","msg":"started","port":8080}`
	f := Extract(line)

	if f["level"] != "info" {
		t.Errorf("expected level=info, got %q", f["level"])
	}
	if f["msg"] != "started" {
		t.Errorf("expected msg=started, got %q", f["msg"])
	}
	if f["port"] != "8080" {
		t.Errorf("expected port=8080, got %q", f["port"])
	}
}

func TestExtractKV(t *testing.T) {
	line := `time=2024-01-01T00:00:00Z level=warn msg="disk full" host=web01`
	f := Extract(line)

	if f["level"] != "warn" {
		t.Errorf("expected level=warn, got %q", f["level"])
	}
	if f["host"] != "web01" {
		t.Errorf("expected host=web01, got %q", f["host"])
	}
	if f["msg"] != "disk full" {
		t.Errorf("expected msg='disk full', got %q", f["msg"])
	}
}

func TestExtractEmptyLine(t *testing.T) {
	f := Extract("")
	if len(f) != 0 {
		t.Errorf("expected empty fields for empty line, got %v", f)
	}
}

func TestExtractPlainText(t *testing.T) {
	// Plain text with no key=value pairs should return empty fields.
	f := Extract("this is a plain log message with no structure")
	if len(f) != 0 {
		t.Errorf("expected empty fields for plain text, got %v", f)
	}
}

func TestExtractJSONFallbackOnInvalid(t *testing.T) {
	// Starts with '{' but is invalid JSON — should return empty.
	f := Extract("{not valid json}")
	if len(f) != 0 {
		t.Errorf("expected empty fields for invalid JSON, got %v", f)
	}
}

func TestExtractKVSkipsTokensWithoutEquals(t *testing.T) {
	line := `INFO level=debug some random words host=srv02`
	f := Extract(line)

	if _, ok := f["INFO"]; ok {
		t.Error("INFO should not be extracted as a key")
	}
	if f["level"] != "debug" {
		t.Errorf("expected level=debug, got %q", f["level"])
	}
	if f["host"] != "srv02" {
		t.Errorf("expected host=srv02, got %q", f["host"])
	}
}
