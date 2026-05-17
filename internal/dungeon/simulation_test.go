package dungeon

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestImpossibleMoveDoesNotChangeState(t *testing.T) {
	config := `{
    "Floors": 2,
    "Monsters": 2,
    "OpenAt": "10:00:00",
    "Duration": 1
}`

	events := `[10:00:00] 7 1
[10:00:01] 7 2
[10:00:02] 7 3
[10:00:03] 7 4
[10:00:04] 7 3
[10:00:05] 7 4
[10:00:06] 7 6
[10:00:07] 7 7
[10:00:08] 7 8`

	expected := `[10:00:00] Player [7] registered
[10:00:01] Player [7] entered the dungeon
[10:00:02] Player [7] killed the monster
[10:00:03] Player [7] makes imposible move [4]
[10:00:04] Player [7] killed the monster
[10:00:05] Player [7] went to the next floor
[10:00:06] Player [7] entered the boss's floor
[10:00:07] Player [7] killed the boss
[10:00:08] Player [7] left the dungeon
Final report:
[SUCCESS] 7 [00:00:07, 00:00:03, 00:00:01] HP:100
`

	got := runScenario(t, config, events)
	if got != expected {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", expected, got)
	}
}
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

func TestEnterBeforeOpenIsImpossible(t *testing.T) {
	config := `{
    "Floors": 2,
    "Monsters": 1,
    "OpenAt": "10:00:00",
    "Duration": 1
}`

	events := `[09:59:59] 9 1
[09:59:59] 9 2
[10:00:00] 9 2
[10:00:01] 9 8`

	expected := `[09:59:59] Player [9] registered
[09:59:59] Player [9] makes imposible move [2]
[10:00:00] Player [9] entered the dungeon
[10:00:01] Player [9] left the dungeon
Final report:
[FAIL] 9 [00:00:01, 00:00:00, 00:00:00] HP:100
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
