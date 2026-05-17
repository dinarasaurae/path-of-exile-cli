package dungeon

import (
	"fmt"
	"time"
)

func (s *Simulator) register(player *playerState, event Event) {
	if player.registered || player.entered {
		s.impossible(player, event)
		return
	}

	player.registered = true
	s.log(event.At, fmt.Sprintf("Player [%d] registered", player.id))
}

func (s *Simulator) enterDungeon(player *playerState, event Event) {
	if player.entered || event.At < s.settings.OpenAt || event.At > s.settings.CloseAt {
		s.impossible(player, event)
		return
	}

	player.entered = true
	player.enteredAt = event.At
	player.currentFloor = 1
	s.activateCurrentFloor(player, event.At)
	s.log(event.At, fmt.Sprintf("Player [%d] entered the dungeon", player.id))
}

func (s *Simulator) killMonster(player *playerState, event Event) {
	if !player.entered || !s.isOrdinaryFloor(player.currentFloor) {
		s.impossible(player, event)
		return
	}

	floor := &player.floors[player.currentFloor-1]
	if floor.cleared || floor.killed >= s.settings.Monsters {
		s.impossible(player, event)
		return
	}

	floor.killed++
	if floor.killed == s.settings.Monsters {
		s.clearCurrentFloor(player, event.At)
	}

	s.log(event.At, fmt.Sprintf("Player [%d] killed the monster", player.id))
}

func (s *Simulator) disqualify(player *playerState, at time.Duration) {
	s.log(at, fmt.Sprintf("Player [%d] is disqualified", player.id))
	s.finish(player, StatusDisqual, at)
}

func (s *Simulator) impossible(player *playerState, event Event) {
	s.log(event.At, fmt.Sprintf("Player [%d] makes imposible move [%d]", player.id, event.ID))
}

func (s *Simulator) finish(player *playerState, status Status, at time.Duration) {
	if player.finished {
		return
	}

	player.status = status
	player.endedAt = at
	player.finished = true
}

func (s *Simulator) expirePlayers(at time.Duration) {}
