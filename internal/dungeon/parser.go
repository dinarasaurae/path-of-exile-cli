package dungeon

import (
	"fmt"
	"strconv"
	"strings"
)

func parseEventLine(line string) (Event, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return Event{}, fmt.Errorf("expected at least three fields")
	}

	timeToken := fields[0]
	if len(timeToken) != 10 || timeToken[0] != '[' || timeToken[9] != ']' {
		return Event{}, fmt.Errorf("invalid time token")
	}

	at, err := ParseClock(timeToken[1:9])
	if err != nil {
		return Event{}, err
	}

	playerID, err := strconv.Atoi(fields[1])
	if err != nil {
		return Event{}, fmt.Errorf("invalid player id")
	}

	eventID, err := strconv.Atoi(fields[2])
	if err != nil {
		return Event{}, fmt.Errorf("invalid event id")
	}

	extra := ""
	if len(fields) > 3 {
		extra = strings.Join(fields[3:], " ")
	}

	return Event{
		At:       at,
		PlayerID: playerID,
		ID:       EventID(eventID),
		Extra:    extra,
	}, nil
}
