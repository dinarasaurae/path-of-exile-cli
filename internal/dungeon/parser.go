package dungeon

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseEvents(reader io.Reader) ([]Event, error) {
	scanner := bufio.NewScanner(reader)
	var events []Event
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		event, err := parseEventLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

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

	if err := validateEventExtra(EventID(eventID), fields[3:], extra); err != nil {
		return Event{}, err
	}

	return Event{
		At:       at,
		PlayerID: playerID,
		ID:       EventID(eventID),
		Extra:    extra,
	}, nil
}

func validateEventExtra(eventID EventID, extraFields []string, extra string) error {
	switch eventID {
	case EventCannotProceed:
		if extra == "" {
			return fmt.Errorf("missing reason")
		}
	case EventRestoreHealth:
		return validateNumericExtra(extraFields, "health")
	case EventReceiveDamage:
		return validateNumericExtra(extraFields, "damage")
	}

	return nil
}

func validateNumericExtra(extraFields []string, name string) error {
	if len(extraFields) != 1 {
		return fmt.Errorf("invalid %s", name)
	}

	value, err := strconv.Atoi(extraFields[0])
	if err != nil || value < 0 {
		return fmt.Errorf("invalid %s", name)
	}

	return nil
}
