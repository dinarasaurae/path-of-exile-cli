package dungeon

import "time"

func (s *Simulator) log(at time.Duration, message string) {
	s.logs = append(s.logs, LogEntry{
		At:      at,
		Message: message,
	})
}
