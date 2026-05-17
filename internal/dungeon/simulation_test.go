package dungeon

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestUnregisteredPlayerIsDisqualified(t *testing.T) {
	config := `{
    "Floors": 2,
    "Monsters": 1,
    "OpenAt": "10:00:00",
    "Duration": 1
}`

	events := `[10:00:00] 3 2`

	expected := `[10:00:00] Player [3] is disqualified
Final report:
[DISQUAL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100
`

	got := runScenario(t, config, events)
	if got != expected {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", expected, got)
	}
}

func TestEventsMustBeSorted(t *testing.T) {
	settings, err := NewSettings(Config{
		Floors:   2,
		Monsters: 1,
		OpenAt:   "10:00:00",
		Duration: 1,
	})
	if err != nil {
		t.Fatalf("new settings: %v", err)
	}

	events := []Event{
		{At: mustClock(t, "10:00:01"), PlayerID: 1, ID: EventRegister},
		{At: mustClock(t, "10:00:00"), PlayerID: 1, ID: EventEnterDungeon},
	}

	_, err = Process(settings, events)
	if !errors.Is(err, ErrEventsNotSorted) {
		t.Fatalf("expected ErrEventsNotSorted, got %v", err)
	}
}

func runScenario(t *testing.T, configText string, eventsText string) string {
	t.Helper()

	settings, err := LoadSettings(strings.NewReader(configText))
	if err != nil {
		t.Fatalf("load settings: %v", err)
	}

	events, err := ParseEvents(strings.NewReader(eventsText))
	if err != nil {
		t.Fatalf("parse events: %v", err)
	}

	result, err := Process(settings, events)
	if err != nil {
		t.Fatalf("process events: %v", err)
	}

	return FormatResult(result)
}

func mustClock(t *testing.T, value string) time.Duration {
	t.Helper()

	parsed, err := ParseClock(value)
	if err != nil {
		t.Fatalf("parse clock: %v", err)
	}

	return parsed
}
