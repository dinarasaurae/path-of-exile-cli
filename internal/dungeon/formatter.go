package dungeon

import (
	"fmt"
	"strings"
)

func FormatResult(result Result) string {
	var builder strings.Builder

	for _, log := range result.Logs {
		_, _ = fmt.Fprintf(&builder, "[%s] %s\n", FormatClock(log.At), log.Message)
	}

	_, _ = builder.WriteString("Final report:\n")
	for _, report := range result.Reports {
		_, _ = fmt.Fprintf(
			&builder,
			"[%s] %d [%s, %s, %s] HP:%d\n",
			report.Status,
			report.PlayerID,
			FormatClock(report.TotalTime),
			FormatClock(report.AverageFloorTime),
			FormatClock(report.BossTime),
			report.Health,
		)
	}

	return builder.String()
}
