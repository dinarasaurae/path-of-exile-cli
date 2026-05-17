package dungeon

import (
	"fmt"
	"strconv"
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

func (s *Simulator) moveNext(player *playerState, event Event) {
	if !player.entered || player.currentFloor <= 0 || player.currentFloor >= s.settings.Floors {
		s.impossible(player, event)
		return
	}

	if s.isOrdinaryFloor(player.currentFloor) && !player.floors[player.currentFloor-1].cleared {
		s.impossible(player, event)
		return
	}

	player.currentFloor++
	player.onBossFloor = false
	s.activateCurrentFloor(player, event.At)
	s.log(event.At, fmt.Sprintf("Player [%d] went to the next floor", player.id))
}

func (s *Simulator) movePrevious(player *playerState, event Event) {
	if !player.entered || player.currentFloor <= 1 {
		s.impossible(player, event)
		return
	}

	s.pauseCurrentFloor(player, event.At)
	player.currentFloor--
	player.onBossFloor = false
	s.activateCurrentFloor(player, event.At)
	s.log(event.At, fmt.Sprintf("Player [%d] went to the previous floor", player.id))
}

func (s *Simulator) enterBoss(player *playerState, event Event) {
	if !player.entered || player.currentFloor != s.settings.Floors || player.onBossFloor || player.bossKilled {
		s.impossible(player, event)
		return
	}

	if !s.allOrdinaryFloorsCleared(player) {
		s.impossible(player, event)
		return
	}

	player.onBossFloor = true
	player.bossActive = true
	player.bossActiveStartedAt = event.At
	s.log(event.At, fmt.Sprintf("Player [%d] entered the boss's floor", player.id))
}

func (s *Simulator) killBoss(player *playerState, event Event) {
	if !player.entered || !player.onBossFloor || player.bossKilled {
		s.impossible(player, event)
		return
	}

	player.bossKillDuration = player.bossElapsed + event.At - player.bossActiveStartedAt
	player.bossActive = false
	player.bossKilled = true
	s.log(event.At, fmt.Sprintf("Player [%d] killed the boss", player.id))
}

func (s *Simulator) leaveDungeon(player *playerState, event Event) {
	if !player.entered {
		s.impossible(player, event)
		return
	}

	s.pauseCurrentFloor(player, event.At)
	s.log(event.At, fmt.Sprintf("Player [%d] left the dungeon", player.id))
	s.finish(player, player.finalStatus(), event.At)
}

func (s *Simulator) cannotProceed(player *playerState, event Event) {
	if event.Extra == "" {
		s.impossible(player, event)
		return
	}

	s.pauseCurrentFloor(player, event.At)
	s.log(event.At, fmt.Sprintf("Player [%d] cannot continue due to [%s]", player.id, event.Extra))
	s.finish(player, StatusDisqual, event.At)
}

func (s *Simulator) restoreHealth(player *playerState, event Event) error {
	if !player.entered {
		s.impossible(player, event)
		return nil
	}

	health, err := strconv.Atoi(event.Extra)
	if err != nil {
		return fmt.Errorf("invalid health for player %d", player.id)
	}

	player.hp = min(player.hp+health, 100)

	s.log(event.At, fmt.Sprintf("Player [%d] has restored [%d] of health", player.id, health))
	return nil
}

func (s *Simulator) receiveDamage(player *playerState, event Event) error {
	if !player.entered {
		s.impossible(player, event)
		return nil
	}

	damage, err := strconv.Atoi(event.Extra)
	if err != nil {
		return fmt.Errorf("invalid damage for player %d", player.id)
	}

	player.hp = max(player.hp-damage, 0)

	s.log(event.At, fmt.Sprintf("Player [%d] recieved [%d] of damage", player.id, damage))

	if player.hp == 0 {
		s.pauseCurrentFloor(player, event.At)
		s.log(event.At, fmt.Sprintf("Player [%d] is dead", player.id))
		s.finish(player, StatusFail, event.At)
	}

	return nil
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

func (s *Simulator) expirePlayers(at time.Duration) {
	if at <= s.settings.CloseAt {
		return
	}

	for _, player := range s.players {
		if !player.entered || player.finished {
			continue
		}

		s.pauseCurrentFloor(player, s.settings.CloseAt)
		s.finish(player, player.finalStatus(), s.settings.CloseAt)
	}
}
