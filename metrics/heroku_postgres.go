package metrics

import (
	herokuLog "heroku-logs-exporter/heroku_log"
)

// https://devcenter.heroku.com/articles/heroku-postgres-metrics-logs#database-metrics

type HerokuPostgresMetrics struct {
	Metrics []HerokuMetric
}

func NewHerokuPostgresMetrics() *HerokuPostgresMetrics {
	labels := []string{"app_name", "source", "addon"}

	return &HerokuPostgresMetrics{
		[]HerokuMetric{
			NewHerokuGaugeMetric(
				"sample#db_size",
				"heroku_postgres_metrics_db_size_bytes",
				"The number of bytes contained in the database. This includes all table and index data on disk, including database bloat.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#tables",
				"heroku_postgres_metrics_table_count",
				"The number of tables in the database.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#active-connections",
				"heroku_postgres_metrics_active_connection_count",
				"The number of connections established on the database.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#waiting-connections",
				"heroku_postgres_metrics_waiting_connection_count",
				"Number of connections waiting on a lock to be acquired. If many connections are waiting, this can be a sign of mishandled database concurrency.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#current_transaction",
				"heroku_postgres_metrics_current_transaction",
				"The current transaction ID, which can be used to track writes over time.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#index-cache-hit-rate",
				"heroku_postgres_metrics_index_cache_hit_rate",
				"Ratio of index lookups served from shared buffer cache, rounded to five decimal points.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#table-cache-hit-rate",
				"heroku_postgres_metrics_table_cache_hit_rate",
				"Ratio of table lookups served from shared buffer cache, rounded to five decimal points.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#follower-lag-commits",
				"heroku_postgres_metrics_follower_lag_commit_count",
				"Replication lag, measured as the number of commits that this follower is behind its leader. Replication is asynchronous so a number greater than zero may not indicate an issue, however an increasing value deserves investigation.",
				labels,
				nil,
			),

			NewHerokuGaugeMetric(
				"sample#load-avg-1m",
				"heroku_postgres_metrics_load_avg_1m",
				"The average system load over a period of 1 minute divided by the number of available CPUs. A load-avg of 1.0 indicates that, on average, processes were requesting CPU resources for 100%% of the timespan. This number includes I/O wait.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#load-avg-5m",
				"heroku_postgres_metrics_load_avg_5m",
				"The average system load over a period of 5 minutes divided by the number of available CPUs. A load-avg of 1.0 indicates that, on average, processes were requesting CPU resources for 100%% of the timespan. This number includes I/O wait.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#load-avg-15m",
				"heroku_postgres_metrics_load_avg_15m",
				"The average system load over a period of 15 minutes divided by the number of available CPUs. A load-avg of 1.0 indicates that, on average, processes were requesting CPU resources for 100%% of the timespan. This number includes I/O wait.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#read-iops",
				"heroku_postgres_metrics_read_iops",
				"Number of read operations in I/O sizes of 16KB blocks.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#write-iops",
				"heroku_postgres_metrics_write_iops",
				"Number of write operations in I/O sizes of 16KB blocks.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#memory-total",
				"heroku_postgres_metrics_memory_total_bytes",
				"Total amount of server memory available.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#memory-free",
				"heroku_postgres_metrics_memory_free_bytes",
				"Amount of free memory available.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#memory-cached",
				"heroku_postgres_metrics_memory_cached_bytes",
				"Amount of memory being used by the OS for page cache.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#memory-postgres",
				"heroku_postgres_metrics_memory_postgres_bytes",
				"Approximate amount of memory used by your database’s Postgres processes. This includes shared buffer cache as well as memory for each connection.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#tmp-disk-used",
				"heroku_postgres_metrics_tmp_disk_used_bytes",
				"Amount of bytes used on tmp mount.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#tmp-disk-available",
				"heroku_postgres_metrics_tmp_disk_available_bytes",
				"Amount of bytes available on tmp mount.",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#wal-percentage-used",
				"heroku_postgres_metrics_wal_percentage_used",
				"Percentage of the WAL disk that has been used, between 0.0 and 1.0.",
				labels,
				nil,
			),
		},
	}
}

func (m *HerokuPostgresMetrics) UpdateFromLog(log *herokuLog.HerokuLog) {
	if log.Source != "heroku" || log.Dyno != "heroku-postgres" {
		return
	}

	labels := []string{log.AppName, log.ValueOrUnknown("source"), log.ValueOrUnknown("addon")}
	updateMetricsFromLog(m.Metrics, labels, log)
}
