package metrics

import (
	herokuLog "heroku-logs-exporter/heroku_log"
)

type HerokuRouterMetrics struct {
	Metrics []HerokuMetric
}

func NewHerokuRouterMetrics() *HerokuRouterMetrics {
	labels := []string{"app_name", "dyno", "host", "method", "protocol", "status"}

	return &HerokuRouterMetrics{
		[]HerokuMetric{
			NewHerokuSummaryMetric(
				"service",
				"heroku_router_service_duration_seconds",
				"Request service duration reported by Heroku Router as summary.",
				labels,
				herokuLog.ParseMillis,
			),
			NewHerokuSummaryMetric(
				"connect",
				"heroku_router_connect_duration_seconds",
				"Request connect duration reported by Heroku Router as summary.",
				labels,
				herokuLog.ParseMillis,
			),

			NewHerokuHistogramMetric(
				"service",
				"heroku_router_service_duration_histogram_seconds",
				"Request service duration reported by Heroku Router as histogram.",
				labels,
				nil,
				herokuLog.ParseMillis,
			),
			NewHerokuHistogramMetric(
				"connect",
				"heroku_router_connect_duration_histogram_seconds",
				"Request connect duration reported by Heroku Router as histogram.",
				labels,
				[]float64{.001, .002, .003, .004, .005, .01, 0.025, .05, .1, .25, .5, 1.0, 2.5, 5.0, 10.0, 20.0},
				herokuLog.ParseMillis,
			),
		},
	}
}

func (m *HerokuRouterMetrics) UpdateFromLog(hLog *herokuLog.HerokuLog) {
	if hLog.Source != "heroku" || hLog.Dyno != "router" {
		return
	}

	labels := []string{hLog.AppName, hLog.ValueOrUnknown("dyno"), hLog.ValueOrUnknown("host"), hLog.ValueOrUnknown("method"), hLog.ValueOrUnknown("protocol"), hLog.ValueOrUnknown("status")}
	updateMetricsFromLog(m.Metrics, labels, hLog)
}
