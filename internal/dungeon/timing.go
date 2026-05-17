package dungeon

import "time"

func (s *Simulator) activateCurrentFloor(player *playerState, at time.Duration) {
	if !s.isOrdinaryFloor(player.currentFloor) {
		return
	}

	floor := &player.floors[player.currentFloor-1]
	floor.visited = true

	if floor.cleared {
		return
	}

	if s.settings.Monsters == 0 {
		floor.cleared = true
		floor.clearDuration = 0
		return
	}

	floor.active = true
	floor.activeStartedAt = at
}

func (s *Simulator) clearCurrentFloor(player *playerState, at time.Duration) {
	floor := &player.floors[player.currentFloor-1]
	floor.clearDuration = floor.elapsed + at - floor.activeStartedAt
	floor.cleared = true
	floor.active = false
}

func (s *Simulator) isOrdinaryFloor(floor int) bool {
	return floor >= 1 && floor <= s.settings.OrdinaryCount
}
