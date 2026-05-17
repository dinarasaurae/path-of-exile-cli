package dungeon

import (
	"cmp"
	"slices"
	"time"
)

func (p *playerState) finalStatus() Status {
	if p.allCleared() {
		return StatusSuccess
	}

	return StatusFail
}

func (p *playerState) allCleared() bool {
	for index := range p.floors {
		if !p.floors[index].cleared {
			return false
		}
	}

	return p.bossKilled
}

func (s *Simulator) result() Result {
	reports := make([]Report, 0, len(s.players))
	for _, player := range s.players {
		reports = append(reports, player.report())
	}

	slices.SortFunc(reports, func(left Report, right Report) int {
		return cmp.Compare(left.PlayerID, right.PlayerID)
	})

	return Result{
		Logs:    s.logs,
		Reports: reports,
	}
}

func (p *playerState) report() Report {
	var totalTime time.Duration
	if p.entered && p.finished {
		totalTime = p.endedAt - p.enteredAt
	}

	clearedFloors := 0
	var floorTime time.Duration
	for index := range p.floors {
		if p.floors[index].cleared {
			clearedFloors++
			floorTime += p.floors[index].clearDuration
		}
	}

	var averageFloorTime time.Duration
	if clearedFloors > 0 {
		averageFloorTime = floorTime / time.Duration(clearedFloors)
	}

	return Report{
		Status:           p.status,
		PlayerID:         p.id,
		TotalTime:        totalTime,
		AverageFloorTime: averageFloorTime,
		BossTime:         p.bossKillDuration,
		Health:           p.hp,
	}
}

func (s *Simulator) log(at time.Duration, message string) {
	s.logs = append(s.logs, LogEntry{
		At:      at,
		Message: message,
	})
}
