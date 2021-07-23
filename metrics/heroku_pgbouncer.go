package metrics

import herokuLog "heroku-logs-exporter/heroku_log"

// https://devcenter.heroku.com/articles/heroku-postgres-metrics-logs#pgbouncer-metrics

type HerokuPgbouncerMetrics struct {
	Metrics []HerokuMetric
}

func NewHerokuPgbouncerMetrics() *HerokuPgbouncerMetrics {
	labels := []string{"app_name", "source", "addon"}

	return &HerokuPgbouncerMetrics{
		[]HerokuMetric{
			NewHerokuGaugeMetric(
				"sample#client_active",
				"heroku_pgbouncer_metrics_client_active_count",
				"The number of client connections to the pooler that have an active server connection assignment.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#client_waiting",
				"heroku_pgbouncer_metrics_client_waiting_count",
				"The number of client connections to the pooler that are waiting for a server connection assignment.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#server_active",
				"heroku_pgbouncer_metrics_server_active_count",
				"The number of server connections that are currently assigned to a client connection.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#server_idle",
				"heroku_pgbouncer_metrics_server_idle_count",
				"The number of server connections that are not currently assigned to a client connection.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#max_wait",
				"heroku_pgbouncer_metrics_max_wait_seconds",
				"The longest wait time of any client currently waiting for a server connection assignment.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#avg_query",
				"heroku_pgbouncer_metrics_avg_query_seconds",
				"The average query time of all queries executed through through poolec connections.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#avg_recv",
				"heroku_pgbouncer_metrics_avg_recv_bytes",
				"The average amount of bytes received from clients per second.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#avg_sent",
				"heroku_pgbouncer_metrics_avg_sent_bytes",
				"The average amount of bytes sent to clients per second.",
				labels,
				nil,
			),
		},
	}
}

func (m *HerokuPgbouncerMetrics) UpdateFromLog(log *herokuLog.HerokuLog) {
	if log.Source != "app" || log.Dyno != "heroku-pgbouncer" {
		return
	}

	labels := []string{log.AppName, log.ValueOrUnknown("source"), log.ValueOrUnknown("addon")}
	updateMetricsFromLog(m.Metrics, labels, log)
}
