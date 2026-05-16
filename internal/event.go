package internal

import "time"

type EventID int
type Status string

type Event struct {
	At       time.Duration
	PlayerID int
	ID       EventID
	Extra    string
}

type Report struct {
	Status           Status
	PlayerID         int
	TotalTime        time.Duration
	AverageFloorTime time.Duration
	BossTime         time.Duration
	Health           int
}

type Result struct {
	Logs    []LogEntry
	Reports []Report
}
type LogEntry struct {
	At      time.Duration
	Message string
}
