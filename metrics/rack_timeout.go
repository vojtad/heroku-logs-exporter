package metrics

import (
	"strings"

	herokuLog "heroku-logs-exporter/heroku_log"
)

type RackTimeoutMetrics struct {
	Metrics []HerokuMetric
}

func NewRackTimeoutMetrics() *RackTimeoutMetrics {
	labels := []string{"app_name", "dyno"}

	return &RackTimeoutMetrics{
		[]HerokuMetric{
			NewHerokuSummaryMetric(
				"wait",
				"heroku_rack_timeout_wait_duration_seconds",
				"Request wait duration reported by rack-timeout as summary.",
				labels,
				herokuLog.ParseMillis,
			),
			NewHerokuSummaryMetric(
				"service",
				"heroku_rack_timeout_service_duration_seconds",
				"Request service duration reported by rack-timeout as summary.",
				labels,
				herokuLog.ParseMillis,
			),

			NewHerokuHistogramMetric(
				"wait",
				"heroku_rack_timeout_wait_duration_histogram_seconds",
				"Request wait duration reported by rack-timeout as histogram.",
				labels,
				nil,
				herokuLog.ParseMillis,
			),
			NewHerokuHistogramMetric(
				"service",
				"heroku_rack_timeout_service_duration_histogram_seconds",
				"Request service duration reported by rack-timeout as histogram.",
				labels,
				nil,
				herokuLog.ParseMillis,
			),
		},
	}
}

func (m *RackTimeoutMetrics) UpdateFromLog(log *herokuLog.HerokuLog) {
	if log.Source != "app" {
		return
	}

	if !strings.HasPrefix(log.Dyno, "web.") {
		return
	}

	if !strings.Contains(log.Line, "source=rack-timeout") || !strings.Contains(log.Line, "state=completed") {
		return
	}

	labels := []string{log.AppName, log.Dyno}
	updateMetricsFromLog(m.Metrics, labels, log)
}
