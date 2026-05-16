package dungeon

import (
	"fmt"
	"time"
)

func (s *Simulator) impossible(player *playerState, event Event) {
	s.log(event.At, fmt.Sprintf("Player [%d] makes imposible move [%d]", player.id, event.ID))
}

func (s *Simulator) expirePlayers(at time.Duration) {}
