package dungeon

import (
	"strings"
	"testing"
)

func TestParseEventsKeepsMultiWordExtra(t *testing.T) {
	events, err := ParseEvents(strings.NewReader(`[10:00:00] 1 9 connection lost near exit`))
	if err != nil {
		t.Fatalf("parse events: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected one event, got %d", len(events))
	}

	if events[0].Extra != "connection lost near exit" {
		t.Fatalf("unexpected extra: %q", events[0].Extra)
	}
}

func TestParseEventsRejectsInvalidNumericExtra(t *testing.T) {
	_, err := ParseEvents(strings.NewReader(`[10:00:00] 1 11 heavy hit`))
	if err == nil {
		t.Fatalf("expected parse error")
	}
}
