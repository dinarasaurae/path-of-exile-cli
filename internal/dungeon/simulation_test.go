package dungeon

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestReadmeExample(t *testing.T) {
	config := `{
    "Floors": 2,
    "Monsters": 2,
    "OpenAt": "14:05:00",
    "Duration": 2
}`

	events := `[14:00:00] 1 1
[14:00:00] 2 1
[14:10:00] 2 2
[14:10:00] 3 2
[14:11:00] 2 5
[14:12:00] 3 3
[14:14:00] 2 3
[14:27:00] 2 11 60
[14:29:00] 2 11 50
[14:40:00] 1 2
[14:41:00] 1 3
[14:44:00] 1 11 50
[14:45:00] 1 3
[14:48:00] 1 4
[14:48:00] 1 6
[14:49:00] 1 11 25
[14:49:02] 1 10 80
[14:50:00] 1 11 65
[14:59:00] 1 7
[15:04:00] 1 8`

	expected := readGolden(t, "testdata/output.golden")

	got := runScenario(t, config, events)
	if got != expected {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", expected, got)
	}
}

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

func TestFloorClearTimeStopsAtLastMonster(t *testing.T) {
	config := `{
    "Floors": 2,
    "Monsters": 1,
    "OpenAt": "10:00:00",
    "Duration": 1
}`

	events := `[10:00:00] 10 1
[10:00:00] 10 2
[10:00:10] 10 3
[10:00:20] 10 4
[10:00:20] 10 6
[10:00:25] 10 7
[10:00:30] 10 8`

	expected := `[10:00:00] Player [10] registered
[10:00:00] Player [10] entered the dungeon
[10:00:10] Player [10] killed the monster
[10:00:20] Player [10] went to the next floor
[10:00:20] Player [10] entered the boss's floor
[10:00:25] Player [10] killed the boss
[10:00:30] Player [10] left the dungeon
Final report:
[SUCCESS] 10 [00:00:30, 00:00:10, 00:00:05] HP:100
`

	got := runScenario(t, config, events)
	if got != expected {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", expected, got)
	}
}

func TestDamageCanKillPlayer(t *testing.T) {
	config := `{
    "Floors": 2,
    "Monsters": 1,
    "OpenAt": "10:00:00",
    "Duration": 1
}`

	events := `[10:00:00] 6 1
[10:00:00] 6 2
[10:00:01] 6 11 150
[10:00:02] 6 3`

	expected := `[10:00:00] Player [6] registered
[10:00:00] Player [6] entered the dungeon
[10:00:01] Player [6] recieved [150] of damage
[10:00:01] Player [6] is dead
Final report:
[FAIL] 6 [00:00:01, 00:00:00, 00:00:00] HP:0
`

	got := runScenario(t, config, events)
	if got != expected {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", expected, got)
	}
}

func TestCannotProceedKeepsMultiWordReason(t *testing.T) {
	config := `{
    "Floors": 2,
    "Monsters": 1,
    "OpenAt": "10:00:00",
    "Duration": 1
}`

	events := `[10:00:00] 4 1
[10:00:01] 4 2
[10:00:02] 4 9 connection lost near exit`

	expected := `[10:00:00] Player [4] registered
[10:00:01] Player [4] entered the dungeon
[10:00:02] Player [4] cannot continue due to [connection lost near exit]
Final report:
[DISQUAL] 4 [00:00:01, 00:00:00, 00:00:00] HP:100
`

	got := runScenario(t, config, events)
	if got != expected {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", expected, got)
	}
}

func TestDungeonClosesActivePlayer(t *testing.T) {
	config := `{
    "Floors": 2,
    "Monsters": 1,
    "OpenAt": "10:00:00",
    "Duration": 1
}`

	events := `[10:00:00] 5 1
[10:00:01] 5 2
[10:00:02] 5 3`

	expected := `[10:00:00] Player [5] registered
[10:00:01] Player [5] entered the dungeon
[10:00:02] Player [5] killed the monster
Final report:
[FAIL] 5 [00:59:59, 00:00:01, 00:00:00] HP:100
`

	got := runScenario(t, config, events)
	if got != expected {
		t.Fatalf("unexpected output\nwant:\n%s\ngot:\n%s", expected, got)
	}
}

func TestRegisteredPlayerWithoutDungeonEntryFails(t *testing.T) {
	config := `{
    "Floors": 2,
    "Monsters": 1,
    "OpenAt": "10:00:00",
    "Duration": 1
}`

	events := `[10:00:00] 8 1`

	expected := `[10:00:00] Player [8] registered
Final report:
[FAIL] 8 [00:00:00, 00:00:00, 00:00:00] HP:100
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

func readGolden(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}

	return string(content)
}
