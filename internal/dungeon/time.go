package dungeon

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseClock(value string) (time.Duration, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("expected HH:MM:SS")
	}

	hours, err := parseClockPart(parts[0], "hours")
	if err != nil {
		return 0, err
	}

	minutes, err := parseClockPart(parts[1], "minutes")
	if err != nil {
		return 0, err
	}

	seconds, err := parseClockPart(parts[2], "seconds")
	if err != nil {
		return 0, err
	}

	if minutes > 59 {
		return 0, fmt.Errorf("minutes out of range")
	}

	if seconds > 59 {
		return 0, fmt.Errorf("seconds out of range")
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second, nil
}

func parseClockPart(value string, name string) (int, error) {
	if len(value) != 2 {
		return 0, fmt.Errorf("%s must have two digits", name)
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be numeric", name)
	}

	if parsed < 0 {
		return 0, fmt.Errorf("%s cannot be negative", name)
	}

	return parsed, nil
}
