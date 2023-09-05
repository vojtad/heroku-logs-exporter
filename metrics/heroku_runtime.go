package metrics

import (
	"strings"

	herokuLog "heroku-logs-exporter/heroku_log"
)

// https://devcenter.heroku.com/articles/log-runtime-metrics

type HerokuRuntimeMetrics struct {
	Metrics []HerokuMetric
}

func NewHerokuRuntimeMetrics() *HerokuRuntimeMetrics {
	labels := []string{"app_name", "dyno", "dyno_id"}

	return &HerokuRuntimeMetrics{
		[]HerokuMetric{
			NewHerokuGaugeMetric(
				"sample#load_avg_1m",
				"heroku_runtime_metrics_load_avg_1m",
				"The load average for the dyno in the last 1 minute. This reflects the number of CPU tasks that are in the ready queue (i.e. waiting to be processed).",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#load_avg_1m",
				"heroku_runtime_metrics_load_avg_5m",
				"The load average for the dyno in the last 1 minute. This reflects the number of CPU tasks that are in the ready queue (i.e. waiting to be processed).",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#load_avg_1m",
				"heroku_runtime_metrics_load_avg_15m",
				"The load average for the dyno in the last 1 minute. This reflects the number of CPU tasks that are in the ready queue (i.e. waiting to be processed).",
				labels,
				nil,
			),
			NewHerokuGaugeMetric(
				"sample#memory_rss",
				"heroku_runtime_metrics_memory_rss_bytes",
				"The portion of the dyno’s memory held in RAM.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#memory_cache",
				"heroku_runtime_metrics_memory_cache_bytes",
				"The portion of the dyno’s memory used as disk cache.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#memory_swap",
				"heroku_runtime_metrics_memory_swap_bytes",
				"The portion of a dyno’s memory stored on disk.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#memory_total",
				"heroku_runtime_metrics_memory_total_bytes",
				"The total memory being used by the dyno, equal to the sum of resident, cache, and swap memory.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#memory_quota",
				"heroku_runtime_metrics_memory_quota_bytes",
				"The resident memory (memory_rss) value at which an R14 is triggered.",
				labels,
				herokuLog.ParseSize,
			),
			NewHerokuGaugeMetric(
				"sample#memory_pgpgin",
				"heroku_runtime_metrics_memory_pgpgin_pages",
				"The cumulative total of the pages written to disk. Sudden high variations on this number can indicate short duration spikes in swap usage. The other memory related metrics are point in time snapshots and can miss short spikes.",
				labels,
				herokuLog.ParseNumberWithPagesSuffix,
			),
			NewHerokuGaugeMetric(
				"sample#memory_pgpgout",
				"heroku_runtime_metrics_memory_pgpgout_pages",
				"The cumulative total of the pages read from disk. Sudden high variations on this number can indicate short duration spikes in swap usage. The other memory related metrics are point in time snapshots and can miss short spikes.",
				labels,
				herokuLog.ParseNumberWithPagesSuffix,
			),
		},
	}
}

func (m *HerokuRuntimeMetrics) UpdateFromLog(log *herokuLog.HerokuLog) {
	if log.Source != "heroku" {
		return
	}

	if !strings.HasPrefix(log.Dyno, "worker.") && !strings.HasPrefix(log.Dyno, "web.") {
		return
	}

	labels := []string{log.AppName, log.Dyno, log.ValueOrUnknown("source")}

	if strings.HasPrefix(log.Line, "State changed") {
		if strings.HasSuffix(log.Line, "to down") {
			deleteMetrics(m.Metrics, labels)
		}

		return
	}

	updateMetricsFromLog(m.Metrics, labels, log)
}
