package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
  listenAddress = flag.String("web.listen-address", ":9841", "Address to listen on for telemetry")
  metricsPath = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
  logsPath = flag.String("web.logs-path", "/logs", "Path under which to accept Heroku Log Drain")
  logsTokenParamName = flag.String("web.logs-token-param-name", "token", "Parameter name to check against token parameter value in Heroku Log Drain requests")
  logsTokenParamValue = flag.String("web.logs-token-param-value", "", "Token to check against token parameter in Heroku Log Drain requests")
)

var (
  herokuRouterServiceDurationSeconds = promauto.NewSummaryVec(
    prometheus.SummaryOpts{
      Name: "heroku_router_service_duration_seconds",
      Help: "Request service duration reported by Heroku Router.",
      Objectives: map[float64]float64{0.01: 0.001, 0.1: 0.01, 0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
    },
    []string{"app_name", "dyno", "host", "method", "protocol", "status"},
  )

  herokuRouterConnectDurationSeconds = promauto.NewSummaryVec(
    prometheus.SummaryOpts{
      Name: "heroku_router_connect_duration_seconds",
      Help: "Request connect duration reported by Heroku Router.",
      Objectives: map[float64]float64{0.01: 0.001, 0.1: 0.01, 0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
    },
    []string{"app_name", "dyno", "host", "method", "protocol", "status"},
  )

  herokuErrorsTotal = promauto.NewCounterVec(
    prometheus.CounterOpts{
      Name: "heroku_errors_total",
      Help: "Heroku errors.",
    },
    []string{"app_name", "dyno", "error"},
  )

  herokuRuntimeMetricsLoadAvg1M = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_load_avg_1m",
      Help: "The load average for the dyno in the last 1 minute. This reflects the number of CPU tasks that are in the ready queue (i.e. waiting to be processed).",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuRuntimeMetricsLoadAvg5M = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_load_avg_5m",
      Help: "The load average for the dyno in the last 5 minutes. This reflects the number of CPU tasks that are in the ready queue (i.e. waiting to be processed).",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuRuntimeMetricsLoadAvg15M = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_load_avg_15m",
      Help: "The load average for the dyno in the last 15 minutes. This reflects the number of CPU tasks that are in the ready queue (i.e. waiting to be processed).",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuRuntimeMetricsMemoryTotal = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_memory_total_bytes",
      Help: "The total memory being used by the dyno, equal to the sum of resident, cache, and swap memory.",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuRuntimeMetricsMemoryRss = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_memory_rss_bytes",
      Help: "The portion of the dyno’s memory held in RAM.",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuRuntimeMetricsMemoryCache = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_memory_cache_bytes",
      Help: "The portion of the dyno’s memory used as disk cache.",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuRuntimeMetricsMemorySwap = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_memory_swap_bytes",
      Help: "The portion of a dyno’s memory stored on disk.",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuRuntimeMetricsMemoryPgpgin = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_memory_pgpgin_pages",
      Help: "The cumulative total of the pages written to disk. Sudden high variations on this number can indicate short duration spikes in swap usage. The other memory related metrics are point in time snapshots and can miss short spikes.",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuRuntimeMetricsMemoryPgpgout = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_memory_pgpgout_pages",
      Help: "The cumulative total of the pages read from disk. Sudden high variations on this number can indicate short duration spikes in swap usage. The other memory related metrics are point in time snapshots and can miss short spikes.",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuRuntimeMetricsMemoryQuota = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_runtime_metrics_memory_quota_bytes",
      Help: "The resident memory (memory_rss) value at which an R14 is triggered.",
    },
    []string{"app_name", "dyno", "dyno_id"},
  )

  herokuPgMetricsCurrentTransaction = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_transactions_total",
      Help: "The current transaction ID, which can be used to track writes over time.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsDbSize = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_db_size_bytes",
      Help: "The number of bytes contained in the database. This includes all table and index data on disk, including database bloat.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsTables = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_table_count",
      Help: "The number of tables in the database.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsActiveConnections = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_active_connection_count",
      Help: "The number of connections established on the database.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsWaitingConnections = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_waiting_connection_count",
      Help: "Number of connections waiting on a lock to be acquired. If many connections are waiting, this can be a sign of mishandled database concurrency.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsIndexCacheHitRate = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_index_cache_hit_rate",
      Help: "Ratio of index lookups served from shared buffer cache, rounded to five decimal points.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsTableCacheHitRate = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_table_cache_hit_rate",
      Help: "Ratio of table lookups served from shared buffer cache, rounded to five decimal points.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsLoadAvg1M = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_load_avg_1m",
      Help: "The average system load over a period of 1 minute divided by the number of available CPUs.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsLoadAvg5M = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_load_avg_5m",
      Help: "The average system load over a period of 5 minutes divided by the number of available CPUs.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsLoadAvg15M = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_load_avg_15m",
      Help: "The average system load over a period of 15 minutes divided by the number of available CPUs.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsReadIops = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_read_iops",
      Help: "Number of read operations in I/O sizes of 16KB blocks.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsWriteIops = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_write_iops",
      Help: "Number of write operations in I/O sizes of 16KB blocks.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsTmpDiskUsed = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_tmp_disk_used_bytes",
      Help: "Amount of bytes used on tmp mount.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsTmpDiskAvailable = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_tmp_disk_available_bytes",
      Help: "Amount of bytes available on tmp mount.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsMemoryTotal = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_memory_total_bytes",
      Help: "Total amount of server memory available.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsMemoryFree = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_memory_free_bytes",
      Help: "Amount of free memory available.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsMemoryCached = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_memory_cached_bytes",
      Help: "Amount of memory being used by the OS for page cache.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsMemoryPostgres = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_memory_postgres_bytes",
      Help: "Approximate amount of memory used by your database’s Postgres processes. This includes shared buffer cache as well as memory for each connection.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgMetricsWalPercentageUsed = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pg_metrics_wal_percentage_used",
      Help: "Percentage of the WAL disk that has been used, between 0.0 and 1.0.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgbouncerMetricsClientActive = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pgbouncer_metrics_client_active_count",
      Help: "The number of client connections to the pooler that have an active server connection assignment.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgbouncerMetricsClientWaiting = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pgbouncer_metrics_client_waiting_count",
      Help: "The number of client connections to the pooler that are waiting for a server connection assignment.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgbouncerMetricsServerActive = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pgbouncer_metrics_server_active_count",
      Help: "The number of server connections that are currently assigned to a client connection.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgbouncerMetricsServerIdle = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pgbouncer_metrica_server_idle_count",
      Help: "The number of server connections that are not currently assigned to a client connection.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgbouncerMetricsMaxWait = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pgbouncer_metrics_max_wait_seconds",
      Help: "The longest wait time of any client currently waiting for a server connection assignment.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgbouncerMetricsAvgQuery = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pgbouncer_metrics_avg_query_seconds",
      Help: "The average query time of all queries executed through through poolec connections.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgbouncerMetricsAvgRecv = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pgbouncer_metrics_avg_recv_bytes",
      Help: "The average amount of bytes received from clients per second.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuPgbouncerMetricsAvgSent = promauto.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "heroku_pgbouncer_metrics_avg_sent_bytes",
      Help: "The average amount of bytes sent to clients per second.",
    },
    []string{"app_name", "source", "addon"},
  )

  herokuRackTimeoutWaitDurationSeconds = promauto.NewSummaryVec(
    prometheus.SummaryOpts{
      Name: "heroku_rack_timeout_wait_duration_seconds",
      Help: "Request wait duration reported rack-timeout Ruby gem.",
      Objectives: map[float64]float64{0.01: 0.001, 0.1: 0.01, 0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
    },
    []string{"app_name", "dyno"},
  )

  herokuRackTimeoutServiceDurationSeconds = promauto.NewSummaryVec(
    prometheus.SummaryOpts{
      Name: "heroku_rack_timeout_service_duration_seconds",
      Help: "Request service duration reported rack-timeout Ruby gem.",
      Objectives: map[float64]float64{0.01: 0.001, 0.1: 0.01, 0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
    },
    []string{"app_name", "dyno"},
  )
)

type HerokuLogHeader struct {
  AppName string

  Unk1 string
  Unk2 string
  Time string
  Host string
  Source string
  Dyno string
}

func parseHerokuLogHeader(appName string, herokuLogHeaderString string) *HerokuLogHeader {
  if (appName == "") {
    appName = "UNKNOWN"
  }

  herokuLogHeaderTokens := strings.Split(herokuLogHeaderString, " ")
  herokuLogHeader := HerokuLogHeader{appName, herokuLogHeaderTokens[0], herokuLogHeaderTokens[1], herokuLogHeaderTokens[2], herokuLogHeaderTokens[3] ,herokuLogHeaderTokens[4] ,herokuLogHeaderTokens[5]}
  
  return &herokuLogHeader
}

func parseLogToMap(log string) map[string]string {
  tokens := strings.Split(log, " ")
  values := make(map[string]string)
  for _, pair := range tokens {
    if (strings.Contains(pair, "=")) {
      parts := strings.Split(pair, "=")
      values[parts[0]] = parts[1]
    }
  }

  return values
}

func processNumberRuntimeMetric(labels []string, values map[string]string, key string, gauge *prometheus.GaugeVec) {
  if value, ok := values[key]; ok {
    load, _ := strconv.ParseFloat(value, 64)
    gauge.WithLabelValues(labels...).Set(load)
  }
}

func processNumberWithSuffixRuntimeMetric(labels []string, values map[string]string, key string, suffix string, gauge *prometheus.GaugeVec) {
  if v, ok := values[key]; ok {
    pages, _ := strconv.ParseFloat(strings.Replace(v, suffix, "", 1), 64)
    gauge.WithLabelValues(labels...).Set(pages)
  }
}

func processBytesRuntimeMetric(labels []string, values map[string]string, key string, gauge *prometheus.GaugeVec) {
  if value, ok := values[key]; ok {
    multiplier := 1
    raw_value := value

    if (strings.HasSuffix(value, "GB")) {
      multiplier = 1024 * 1024 * 1024;
      raw_value = strings.Replace(value, "MB", "", 1)
    } else if (strings.HasSuffix(value, "MB")) {
      multiplier = 1024 * 1024;
      raw_value = strings.Replace(value, "MB", "", 1)
    } else if (strings.HasSuffix(value, "kB")) {
      multiplier = 1024;
      raw_value = strings.Replace(value, "kB", "", 1)
    } else if (strings.HasSuffix(value, "bytes")) {
      multiplier = 1;
      raw_value = strings.Replace(value, "bytes", "", 1)
    }

    raw_number_value, _ := strconv.ParseFloat(raw_value, 64)
    bytes := raw_number_value * float64(multiplier)
    gauge.WithLabelValues(labels...).Set(bytes)
  }
}

func processHerokuSource(header *HerokuLogHeader, log string) {
  if (header.Dyno == "router") {
    processHerokuRouter(header, log)
  } else if (strings.HasPrefix(header.Dyno, "worker.") || strings.HasPrefix(header.Dyno, "web.")) {
    processHerokuDyno(header, log)
  }
}

func processAppSource(header *HerokuLogHeader, log string) {
  if (strings.HasPrefix(header.Dyno, "web.")) {
    processAppWebDyno(header, log)
  } else if (header.Dyno == "heroku-postgres") {
    processAppHerokuPostgres(header, log)
  } else if (header.Dyno == "heroku-pgbouncer") {
    processAppHerokuPgbouncer(header, log)
  }
}

func processHerokuRouter(header *HerokuLogHeader, log string) {
  values := parseLogToMap(log)

  serviceDuration, _ := strconv.ParseFloat(strings.Replace(values["service"], "ms", "", 1), 64)
  serviceDuration =  serviceDuration / 1000.0
  herokuRouterServiceDurationSeconds.WithLabelValues(header.AppName, values["dyno"], values["host"], values["method"], values["protocol"], values["status"]).Observe(serviceDuration)

  connectDuration, _ := strconv.ParseFloat(strings.Replace(values["connect"], "ms", "", 1), 64)
  connectDuration =  connectDuration / 1000.0
  herokuRouterConnectDurationSeconds.WithLabelValues(header.AppName, values["dyno"], values["host"], values["method"], values["protocol"], values["status"]).Observe(connectDuration)
}

func processHerokuDyno(header *HerokuLogHeader, log string) {
  if (strings.HasPrefix(log, "Error")) {
    parts := strings.SplitN(log, " ", 3)
    errorCode := parts[1]

    herokuErrorsTotal.WithLabelValues(header.AppName, header.Dyno, errorCode).Inc()
  } else {
    values := parseLogToMap(log)
    labels := []string{header.AppName, header.Dyno, values["source"]}

    processNumberRuntimeMetric(labels, values, "sample#load_avg_1m", herokuRuntimeMetricsLoadAvg1M)
    processNumberRuntimeMetric(labels, values, "sample#load_avg_5m", herokuRuntimeMetricsLoadAvg5M)
    processNumberRuntimeMetric(labels, values, "sample#load_avg_15m", herokuRuntimeMetricsLoadAvg15M)

    processBytesRuntimeMetric(labels, values, "sample#memory_total", herokuRuntimeMetricsMemoryTotal)
    processBytesRuntimeMetric(labels, values, "sample#memory_rss", herokuRuntimeMetricsMemoryRss)
    processBytesRuntimeMetric(labels, values, "sample#memory_cache", herokuRuntimeMetricsMemoryCache)
    processBytesRuntimeMetric(labels, values, "sample#memory_swap", herokuRuntimeMetricsMemorySwap)
    processBytesRuntimeMetric(labels, values, "sample#memory_quota", herokuRuntimeMetricsMemoryQuota)

    processNumberWithSuffixRuntimeMetric(labels, values, "sample#memory_pgpgin", "pages", herokuRuntimeMetricsMemoryPgpgin)
    processNumberWithSuffixRuntimeMetric(labels, values, "sample#memory_pgpgout", "pages", herokuRuntimeMetricsMemoryPgpgout)    
  }
}

func processAppWebDyno(header *HerokuLogHeader, log string) {
  if (strings.Contains(log, "source=rack-timeout")) {
    processAppWebDynoRackTimeout(header, log)
  }
}

func processAppWebDynoRackTimeout(header *HerokuLogHeader, log string) {
  if (!strings.Contains(log, "state=completed")) {
    return;
  }

  values := parseLogToMap(log)

  waitDuration, _ := strconv.ParseFloat(strings.Replace(values["wait"], "ms", "", 1), 64)
  waitDuration =  waitDuration / 1000.0
  herokuRackTimeoutWaitDurationSeconds.WithLabelValues(header.AppName, header.Dyno).Observe(waitDuration)

  serviceDuration, _ := strconv.ParseFloat(strings.Replace(values["service"], "ms", "", 1), 64)
  serviceDuration =  serviceDuration / 1000.0
  herokuRackTimeoutServiceDurationSeconds.WithLabelValues(header.AppName, header.Dyno).Observe(serviceDuration)
}

func processAppHerokuPostgres(header *HerokuLogHeader, log string) {
  values := parseLogToMap(log)
  labels := []string{header.AppName, values["source"], values["addon"]}

  processNumberRuntimeMetric(labels, values, "sample#current_transaction", herokuPgMetricsCurrentTransaction)
  processBytesRuntimeMetric(labels, values, "sample#db_size", herokuPgMetricsDbSize)
  processNumberRuntimeMetric(labels, values, "sample#tables", herokuPgMetricsTables)
  processNumberRuntimeMetric(labels, values, "sample#active-connections", herokuPgMetricsActiveConnections)
  processNumberRuntimeMetric(labels, values, "sample#waiting-connections", herokuPgMetricsWaitingConnections)
  processNumberRuntimeMetric(labels, values, "sample#index-cache-hit-rate", herokuPgMetricsIndexCacheHitRate)
  processNumberRuntimeMetric(labels, values, "sample#table-cache-hit-rate", herokuPgMetricsTableCacheHitRate)
  processNumberRuntimeMetric(labels, values, "sample#load-avg-1m", herokuPgMetricsLoadAvg1M)
  processNumberRuntimeMetric(labels, values, "sample#load-avg-5m", herokuPgMetricsLoadAvg5M)
  processNumberRuntimeMetric(labels, values, "sample#load-avg-15m", herokuPgMetricsLoadAvg15M)
  processNumberRuntimeMetric(labels, values, "sample#read-iops", herokuPgMetricsReadIops)
  processNumberRuntimeMetric(labels, values, "sample#write-iops", herokuPgMetricsWriteIops)
  processNumberRuntimeMetric(labels, values, "sample#tmp-disk-used", herokuPgMetricsTmpDiskUsed)
  processNumberRuntimeMetric(labels, values, "sample#tmp-disk-available", herokuPgMetricsTmpDiskAvailable)
  processBytesRuntimeMetric(labels, values, "sample#memory-total", herokuPgMetricsMemoryTotal)
  processBytesRuntimeMetric(labels, values, "sample#memory-free", herokuPgMetricsMemoryFree)
  processBytesRuntimeMetric(labels, values, "sample#memory-cached", herokuPgMetricsMemoryCached)
  processBytesRuntimeMetric(labels, values, "sample#memory-postgres", herokuPgMetricsMemoryPostgres)
  processNumberRuntimeMetric(labels, values, "sample#wal-percentage-used", herokuPgMetricsWalPercentageUsed)
}

func processAppHerokuPgbouncer(header *HerokuLogHeader, log string) {
  values := parseLogToMap(log)
  labels := []string{header.AppName, values["source"], values["addon"]}

  processNumberRuntimeMetric(labels, values, "sample#client_active", herokuPgbouncerMetricsClientActive)
  processNumberRuntimeMetric(labels, values, "sample#client_waiting", herokuPgbouncerMetricsClientWaiting)
  processNumberRuntimeMetric(labels, values, "sample#server_active", herokuPgbouncerMetricsServerActive)
  processNumberRuntimeMetric(labels, values, "sample#server_idle", herokuPgbouncerMetricsServerIdle)
  processNumberRuntimeMetric(labels, values, "sample#max_wait", herokuPgbouncerMetricsMaxWait)
  processNumberRuntimeMetric(labels, values, "sample#avg_query", herokuPgbouncerMetricsAvgQuery)
  processNumberRuntimeMetric(labels, values, "sample#avg_recv", herokuPgbouncerMetricsAvgRecv)
  processNumberRuntimeMetric(labels, values, "sample#avg_sent", herokuPgbouncerMetricsAvgSent)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "heroku-logs-exporter")
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
  if (r.Method != "POST") {
    http.Error(w, "Bad Request", http.StatusBadRequest)
    return
  }

  if (*logsTokenParamName != "" && *logsTokenParamValue != "") {
    if (*logsTokenParamValue != r.URL.Query().Get(*logsTokenParamName)) {
      log.Printf("Token mismatch: %s\n", r.URL)
      http.Error(w, "Bad Request", http.StatusBadRequest)
      return
    }
  }

  scanner := bufio.NewScanner(r.Body)
  for scanner.Scan() {
    line := scanner.Text()
    parts := strings.SplitN(line, " - ", 2)
    
    header := parseHerokuLogHeader(r.URL.Query().Get("app_name"), parts[0])
    log := parts[1]

    if (header.Source == "heroku") {
      processHerokuSource(header, log)
    } else if (header.Source == "app") {
      processAppSource(header, log)
    } else {
      // fmt.Fprintf(os.Stderr, "Unhandled source: %s - %+v\n", header.Source, header)
    }
  }
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
  if (r.Method != "GET") {
    http.Error(w, "Bad Request", http.StatusBadRequest)
    return
  }

  fmt.Fprintf(w, "heroku-logs-exporter: metrics")
}

func main() {
  flag.Parse()

  http.HandleFunc("/", helloHandler)
  
  http.Handle(*metricsPath, promhttp.Handler())
  http.HandleFunc(*logsPath, logsHandler)

  log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
