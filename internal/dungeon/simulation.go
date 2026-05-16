package dungeon

import (
	"errors"
	"time"
)

var ErrEventsNotSorted = errors.New("events are not sorted by time")

type Simulator struct {
	settings Settings
	players  map[int]*playerState
	logs     []LogEntry
}

type playerState struct {
	id                  int
	hp                  int
	registered          bool
	entered             bool
	finished            bool
	status              Status
	enteredAt           time.Duration
	endedAt             time.Duration
	currentFloor        int
	floors              []floorState
	onBossFloor         bool
	bossKilled          bool
	bossElapsed         time.Duration
	bossActive          bool
	bossActiveStartedAt time.Duration
	bossKillDuration    time.Duration
}

type floorState struct {
	killed          int
	cleared         bool
	active          bool
	activeStartedAt time.Duration
	elapsed         time.Duration
	clearDuration   time.Duration
	visited         bool
}

func Process(settings Settings, events []Event) (Result, error) {
	simulator := NewSimulator(settings)

	for index, event := range events {
		if index > 0 && event.At < events[index-1].At {
			return Result{}, ErrEventsNotSorted
		}

		simulator.expirePlayers(event.At)

		if err := simulator.handle(event); err != nil {
			return Result{}, err
		}
	}

	simulator.expirePlayers(settings.CloseAt + time.Second)

	return simulator.result(), nil
}

func NewSimulator(settings Settings) *Simulator {
	return &Simulator{
		settings: settings,
		players:  make(map[int]*playerState),
	}
}

func (s *Simulator) handle(event Event) error {
	player := s.player(event.PlayerID)
	if player.finished {
		return nil
	}

	if event.ID == EventRegister {
		return nil
	}

	if !player.registered {
		return nil
	}

	switch event.ID {
	case EventEnterDungeon:
		return nil
	case EventKillMonster:
		return nil
	case EventNextFloor:
		return nil
	case EventPreviousFloor:
		return nil
	case EventEnterBoss:
		return nil
	case EventKillBoss:
		return nil
	case EventLeaveDungeon:
		return nil
	case EventCannotProceed:
		return nil
	case EventRestoreHealth:
		return nil
	case EventReceiveDamage:
		return nil
	default:
		s.impossible(player, event)
	}

	return nil
}

func (s *Simulator) player(id int) *playerState {
	player, ok := s.players[id]
	if ok {
		return player
	}

	player = &playerState{
		id:     id,
		hp:     100,
		status: StatusFail,
		floors: make([]floorState, s.settings.OrdinaryCount),
	}
	s.players[id] = player

	return player
}
