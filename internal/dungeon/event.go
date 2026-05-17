package dungeon

import "time"

type EventID int

const (
	EventRegister      EventID = 1
	EventEnterDungeon  EventID = 2
	EventKillMonster   EventID = 3
	EventNextFloor     EventID = 4
	EventPreviousFloor EventID = 5
	EventEnterBoss     EventID = 6
	EventKillBoss      EventID = 7
	EventLeaveDungeon  EventID = 8
	EventCannotProceed EventID = 9
	EventRestoreHealth EventID = 10
	EventReceiveDamage EventID = 11
)

type Status string

const (
	StatusSuccess Status = "SUCCESS"
	StatusFail    Status = "FAIL"
	StatusDisqual Status = "DISQUAL"
)

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
