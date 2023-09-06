package metrics

import (
	herokuLog "heroku-logs-exporter/heroku_log"
	"strings"
)

type HerokuSystemMetrics struct {
	Metrics []HerokuMetric
}

func NewHerokuSystemMetrics() *HerokuSystemMetrics {
	labels := []string{"app_name", "dyno", "error"}

	return &HerokuSystemMetrics{
		[]HerokuMetric{
			NewHerokuCounterMetric(
				"error",
				"heroku_system_error_count",
				"Heroku errors.",
				labels,
			),
		},
	}
}

func (m *HerokuSystemMetrics) UpdateFromLog(hLog *herokuLog.HerokuLog) {
	if hLog.Source != "heroku" {
		return
	}

	if !strings.HasPrefix(hLog.Line, "Error ") {
		return
	}

	parts := strings.SplitN(hLog.Line, " ", 3)
	errorCode := parts[1]

	labels := []string{hLog.AppName, hLog.Dyno, errorCode}
	updateMetricFromLog(m.Metrics, "error", labels, "")
}
