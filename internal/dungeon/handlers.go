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
