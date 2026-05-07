package levelfilter

import "testing"

func TestParseKnownLevels(t *testing.T) {
	cases := []struct {
		input string
		want  Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"warning", LevelWarn},
		{"WARN", LevelWarn},
		{"error", LevelError},
		{"err", LevelError},
		{"fatal", LevelFatal},
		{"crit", LevelFatal},
		{"critical", LevelFatal},
		{"unknown_xyz", LevelUnknown},
		{"", LevelUnknown},
	}
	for _, c := range cases {
		got := Parse(c.input)
		if got != c.want {
			t.Errorf("Parse(%q) = %d, want %d", c.input, got, c.want)
		}
	}
}

func TestAllowKVFormat(t *testing.T) {
	f := New(LevelWarn, "")
	if !f.Allow(`time=2024-01-01 level=warn msg="disk low"`) {
		t.Error("expected warn to pass at min=warn")
	}
	if !f.Allow(`time=2024-01-01 level=error msg="disk full"`) {
		t.Error("expected error to pass at min=warn")
	}
	if f.Allow(`time=2024-01-01 level=debug msg="heartbeat"`) {
		t.Error("expected debug to be dropped at min=warn")
	}
	if f.Allow(`time=2024-01-01 level=info msg="starting"`) {
		t.Error("expected info to be dropped at min=warn")
	}
}

func TestAllowJSONFormat(t *testing.T) {
	f := New(LevelError, "level")
	if !f.Allow(`{"time":"2024-01-01","level":"error","msg":"boom"}`) {
		t.Error("expected error to pass at min=error")
	}
	if !f.Allow(`{"time":"2024-01-01","level":"fatal","msg":"crash"}`) {
		t.Error("expected fatal to pass at min=error")
	}
	if f.Allow(`{"time":"2024-01-01","level":"info","msg":"ok"}`) {
		t.Error("expected info to be dropped at min=error")
	}
}

func TestNoLevelFieldPassesThrough(t *testing.T) {
	f := New(LevelError, "level")
	line := "plain text log line without any level field"
	if !f.Allow(line) {
		t.Error("expected line with no level field to pass through")
	}
}

func TestCounters(t *testing.T) {
	f := New(LevelWarn, "level")
	f.Allow(`level=debug msg="a"`)
	f.Allow(`level=info msg="b"`)
	f.Allow(`level=warn msg="c"`)
	f.Allow(`level=error msg="d"`)
	f.Allow("no level here")

	if f.Dropped != 2 {
		t.Errorf("Dropped = %d, want 2", f.Dropped)
	}
	if f.Passed != 3 {
		t.Errorf("Passed = %d, want 3", f.Passed)
	}
}

func TestCustomField(t *testing.T) {
	f := New(LevelInfo, "severity")
	if !f.Allow(`severity=info component=api`) {
		t.Error("expected info to pass with custom field name")
	}
	if f.Allow(`severity=debug component=api`) {
		t.Error("expected debug to be dropped with custom field name")
	}
}

func TestCountersReset(t *testing.T) {
	f := New(LevelWarn, "level")
	f.Allow(`level=debug msg="a"`)
	f.Allow(`level=warn msg="b"`)
	f.Allow(`level=error msg="c"`)

	if f.Dropped != 1 {
		t.Fatalf("before reset: Dropped = %d, want 1", f.Dropped)
	}
	if f.Passed != 2 {
		t.Fatalf("before reset: Passed = %d, want 2", f.Passed)
	}

	f.Reset()

	if f.Dropped != 0 {
		t.Errorf("after reset: Dropped = %d, want 0", f.Dropped)
	}
	if f.Passed != 0 {
		t.Errorf("after reset: Passed = %d, want 0", f.Passed)
	}
}
